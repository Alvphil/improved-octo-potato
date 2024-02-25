package database

import "strconv"

type Chirp struct {
	Author_id int    `json:"author_id"`
	Body      string `json:"body"`
	ID        int    `json:"id"`
}

func (db *DB) CreateChirp(body, author string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	authorId, err := strconv.Atoi(author)
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:        id,
		Body:      body,
		Author_id: authorId,
	}
	dbStructure.Chirps[id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return ErrNotExist
	}

	deleted := Chirp{
		ID:        chirp.ID,
		Body:      "",
		Author_id: -1,
	}
	dbStructure.Chirps[id] = deleted
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
