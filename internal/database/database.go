package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

var ErrNotExist = errors.New("resource does not exist")

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps               map[int]Chirp                   `json:"chirps"`
	Users                map[int]User                    `json:"users"`
	RevokedRefreshTokens map[string]RevokedRefreshTokens `json:"revokedRefreshTokens"`
}

type RevokedRefreshTokens struct {
	Refresh_token string    `json:"refresh_token"`
	TimeStamp     time.Time `json:"timestamp"`
}

func NewDB(path string, debug bool) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	if debug {
		os.Remove(db.path)
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:               map[int]Chirp{},
		Users:                map[int]User{},
		RevokedRefreshTokens: map[string]RevokedRefreshTokens{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	timestamp := time.Now().UTC()
	revoked := RevokedRefreshTokens{
		Refresh_token: refreshToken,
		TimeStamp:     timestamp,
	}
	dbStructure.RevokedRefreshTokens[refreshToken] = revoked
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) GetRevokedToken(refreshToken string) (string, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return "", err
	}

	token, ok := dbStructure.RevokedRefreshTokens[refreshToken]
	if !ok {
		return "", ErrNotExist
	}
	//fmt.Println(token.Refresh_token)
	return token.Refresh_token, nil

}
