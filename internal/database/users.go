package database

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	PasswordHashed []byte `json:"password"`
	Token          string `json:"token"`
}

type PublicUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type UserJWT struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Token         string `json:"token"`
	Refresh_token string `json:"refresh_token"`
}

func (db *DB) CreateUser(email string, passwordHashed []byte) (PublicUser, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return PublicUser{}, err
	}
	_, err = db.CheckUserExists(email)
	if err == nil {
		return PublicUser{}, errors.New("User already exists")
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:             id,
		Email:          email,
		PasswordHashed: passwordHashed,
		Token:          "",
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return PublicUser{}, err
	}
	publicUser := PublicUser{
		ID:    user.ID,
		Email: user.Email,
	}

	return publicUser, nil
}

func (db *DB) GetUser(id int) (PublicUser, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return PublicUser{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return PublicUser{}, ErrNotExist
	}
	publicUser := PublicUser{
		ID:    user.ID,
		Email: user.Email,
	}

	return publicUser, nil
}

func (db *DB) GetLoggedInUser(id int, token, refresh string) (UserJWT, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserJWT{}, err
	}

	user, ok := dbStructure.Users[id]
	if !ok {
		return UserJWT{}, ErrNotExist
	}

	loginResponse := UserJWT{
		ID:            user.ID,
		Email:         user.Email,
		Token:         token,
		Refresh_token: refresh,
	}

	return loginResponse, nil
}

func (db *DB) CheckUserExists(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, errors.New("user does not exist")
}

func (db *DB) UpdateUserJWTToken(jwtToken string, id int) (UserJWT, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserJWT{}, err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return UserJWT{}, ErrNotExist
	}
	user.Token = jwtToken

	UserJWT := UserJWT{
		ID:    id,
		Email: user.Email,
		Token: user.Token,
	}

	return UserJWT, nil

}

func (db *DB) UpdateUser(id int, email, password string) (PublicUser, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return PublicUser{}, err
	}
	user, ok := dbStructure.Users[id]
	if !ok {
		return PublicUser{}, err
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(password), 0)
		user.PasswordHashed = passwordHashed
	}

	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return PublicUser{}, err
	}

	publicUser := PublicUser{
		ID:    user.ID,
		Email: user.Email,
	}
	return publicUser, nil
}
