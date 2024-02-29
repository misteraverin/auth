package postgres

import (
	"auth/models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "testdb"
)

func NewPostgreSQLRepository() (*PostgreSQLRepository, error) {
	db, err := OpenConnectWithDB()
	if err != nil {
		return nil, err
	}

	err = InitializeDBTable(db)
	if err != nil {
		return nil, err
	}
	return &PostgreSQLRepository{db: db}, nil
}

func InitializeDBTable(db *sql.DB) error {
	// Создаем таблицу пользователей, если она еще не существует
	_, err := db.Exec(`
		DROP TABLE IF EXISTS Tokens;
		DROP TABLE IF EXISTS Users;
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS Users (
			Id SERIAL PRIMARY KEY,
			Username VARCHAR(50) UNIQUE,
			Password VARCHAR(50)
		);
		
		CREATE TABLE IF NOT EXISTS Tokens (
			Token VARCHAR(250) PRIMARY KEY,
			Username VARCHAR(50),
			FOREIGN KEY (username) REFERENCES users(username)
		);  
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec(`		       
		INSERT INTO Users (Id, Username, Password) VALUES
			('1', 'user1', 'password1'),
			('2', 'user2', 'password2'),
			('3', 'user3', 'password3')
    `)
	return nil
}

func OpenConnectWithDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db, err
}

type PostgreSQLRepository struct {
	db *sql.DB
}

func (r *PostgreSQLRepository) SaveUser(user *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}

func (r *PostgreSQLRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, username, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgreSQLRepository) GetUserByUserName(username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, username, password FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *PostgreSQLRepository) SaveToken(token, username string) error {
	_, err := r.db.Exec("INSERT INTO tokens (token, username) VALUES ($1, $2)", token, username)
	return err
}

func (r *PostgreSQLRepository) GetUserNameByToken(token string) (string, error) {
	var username string
	err := r.db.QueryRow("SELECT username FROM tokens WHERE token = $1", token).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (r *PostgreSQLRepository) UpdateToken(token, username string) error {
	_, err := r.db.Exec("UPDATE tokens SET token = $1 WHERE username = $2", token, username)
	return err
}

func (r *PostgreSQLRepository) CheckUserHasToken(username string) (bool, error) {
	var token models.Token
	err := r.db.QueryRow("SELECT username FROM tokens WHERE username = $1", username).Scan(&token.Username)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Store представляет интерфейс для хранилища пользователей и JWT токенов
type Repository interface {
	SaveUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	GetUserByUserName(username string) (*models.User, error)

	SaveToken(token, username string) error
	GetUserNameByToken(token string) (string, error)
	UpdateToken(token, username string) error
	CheckUserHasToken(username string) (bool, error)
}
