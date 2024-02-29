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
			HashPassword VARCHAR(250)
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

	//password1, err := bcrypt.GenerateFromPassword([]byte("password1"+salt), bcrypt.DefaultCost)
	//хэш от "password1"+salt
	//хэш от "password2"+salt
	//хэш от "password3"+salt
	_, err = db.Exec(`
		INSERT INTO Users (Id, Username, HashPassword) VALUES
			('1', 'user1', '$2a$10$JmYRexxhbYvJBR1rfS9VA.4AeBbl7r82r3ZQ9duFMlZrwJSHzRmQW'),			
			('2', 'user2', '$2a$10$Nu4Fa6IIROeXNGQ6iHunqOP71ojTCZWeXz9t4VaQ.ehAL6tRU5YFK'),
			('3', 'user3', '$2a$10$diGAXGuA7k7v8tcUSbJUdec0QVUS2N9hUM1KaD.j4z670Yv8VrVl6')
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

func (r *PostgreSQLRepository) GetUserByUserName(username string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow("SELECT id, username, HashPassword FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.HashPassword)
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
	//SaveUser(user *models.User) error
	//GetUserByID(id int) (*models.User, error)
	GetUserByUserName(username string) (*models.User, error)

	SaveToken(token, username string) error
	GetUserNameByToken(token string) (string, error)
	UpdateToken(token, username string) error
	CheckUserHasToken(username string) (bool, error)
}

//func (r *PostgreSQLRepository) SaveUser(user *models.User) error {
//	_, err := r.db.Exec("INSERT INTO users (username, HashPassword) VALUES ($1, $2)", user.Username, user.HashPassword)
//	return err
//}

//func (r *PostgreSQLRepository) GetUserByID(id int) (*models.User, error) {
//	var user models.User
//	err := r.db.QueryRow("SELECT id, username, HashPassword FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.HashPassword)
//	if err != nil {
//		return nil, err
//	}
//	return &user, nil
//}
