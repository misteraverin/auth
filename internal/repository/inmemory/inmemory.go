package inmemory

import (
	"auth/internal/domain/user"
	"auth/internal/errdomain"
	"context"
	"sync"
)

type DB struct {
	users map[string]user.User
	mux   *sync.RWMutex
}

func NewMapDB() (*DB, error) {
	m := make(map[string]user.User)
	db := DB{
		users: m,
		mux:   &sync.RWMutex{},
	}

	return &db, nil
}

func (db *DB) Exist(ctx context.Context, login string) (bool, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	_, ok := db.users[login]
	return ok, nil
}

func (db *DB) IsCorrectPassword(ctx context.Context, login string, password string) (bool, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	u, ok := db.users[login]

	if !ok {
		return false, nil
	}

	return u.Password == password, nil
}

func (db *DB) GetUser(ctx context.Context, login string) (*user.User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	u, ok := db.users[login]

	if !ok {
		return nil, errdomain.ErrUserIsNotExist
	}

	return &u, nil
}

func (db *DB) Save(ctx context.Context, login string, password string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	db.users[login] = user.User{Login: login, Password: password}

	return nil
}
