package repository

import (
	"BasicAuth/internal/domain/entity"
	"BasicAuth/internal/errdomain"
	"sync"
)

type MapDB struct {
	users map[string]entity.User
	mux   sync.RWMutex
}

func NewMapDB() (*MapDB, error) {
	db := MapDB{
		users: make(map[string]entity.User),
		mux:   sync.RWMutex{},
	}

	return &db, nil
}

func (db *MapDB) Exist(login string) bool {
	db.mux.RLock()
	defer db.mux.RUnlock()

	_, ok := db.users[login]
	return ok
}

func (db *MapDB) IsCorrectPassword(login string, password string) bool {
	db.mux.RLock()
	defer db.mux.RUnlock()

	user, ok := db.users[login]

	if !ok {
		return false
	}

	return user.Password == password
}

func (db *MapDB) GetUser(login string) (entity.User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	user, ok := db.users[login]
	var err error

	if !ok {
		err = errdomain.ErrUserIsNotExist
	}

	return user, err
}

func (db *MapDB) Save(login string, password string) {
	db.mux.Lock()
	defer db.mux.Unlock()

	db.users[login] = entity.User{Login: login, Password: password}
}
