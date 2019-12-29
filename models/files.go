package models

import "database/sql"

type File struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
}

func (f *File) parseFile(row *sql.Row) error {
	err := row.Scan(&f.ID, &f.UserID, &f.Name, &f.Filename, &f.ContentType)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) parseFiles(rows *sql.Rows) error {
	err := rows.Scan(&f.ID, &f.UserID, &f.Name, &f.Filename, &f.ContentType)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) NewFile(db *sql.DB) error {
	q := `INSERT INTO files (id, user_id, name, filename, content_type) VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(q, f.ID, f.UserID, f.Name, f.Filename, f.ContentType)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) GetFile(db *sql.DB) error {
	q := `SELECT * FROM files WHERE id = $1 LIMIT 1`

	row := db.QueryRow(q, f.ID)

	err := f.parseFile(row)
	if err != nil {
		return err
	}
	return nil
}

func (f *File) FindFiles(db *sql.DB, name string) ([]*File, error) {
	q := `SELECT * FROM files WHERE name LIKE $1 LIMIT 1`
	var err error
	var files = make([]*File, 0)
	prepped := name + "%"

	rows, err := db.Query(q, prepped)
	if err != nil {
		return files, err
	}

	for rows.Next() {
		f1 := new(File)
		err = f1.parseFiles(rows)
		if err != nil {
			return files, err
		}
	}

	if err != nil {
		return files, err
	}
	return files, nil
}

func (f *File) GetFilesByUserId(db *sql.DB, userId string) ([]*File, error) {
	q := `SELECT * FROM files WHERE user_id = $1`
	var err error
	var files = make([]*File, 0)

	rows, err := db.Query(q, userId)
	if err != nil {
		return files, err
	}

	for rows.Next() {
		f1 := new(File)
		err = f1.parseFiles(rows)
		if err != nil {
			return files, err
		}
	}

	if err != nil {
		return files, err
	}
	return files, nil
}
