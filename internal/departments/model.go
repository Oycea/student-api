package faculties

import (
	"database/sql"
	"errors"
	"time"

	"student-api/internal/storage/postgres"
)

type Faculty struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// -------------------- CREATE --------------------

func CreateFaculty(storage *postgres.Storage, f *Faculty) error {
	query := `
		INSERT INTO faculties (name)
		VALUES ($1)
		RETURNING id
	`

	return storage.DB().QueryRow(query, f.Name).Scan(&f.ID)
}

// -------------------- READ ONE --------------------

func GetFacultyByID(storage *postgres.Storage, id int64) (*Faculty, error) {
	var f Faculty

	query := `
		SELECT id, name
		FROM faculties
		WHERE id=$1
	`

	err := storage.DB().QueryRow(query, id).Scan(
		&f.ID,
		&f.Name,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("faculty not found")
		}
		return nil, err
	}

	return &f, nil
}

// -------------------- READ ALL --------------------

func GetAllFaculties(storage *postgres.Storage) ([]Faculty, error) {
	rows, err := storage.DB().Query(`
		SELECT id, name
		FROM faculties
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var faculties []Faculty

	for rows.Next() {
		var f Faculty
		if err := rows.Scan(&f.ID, &f.Name); err != nil {
			return nil, err
		}
		faculties = append(faculties, f)
	}

	return faculties, nil
}

// -------------------- UPDATE --------------------

func UpdateFaculty(storage *postgres.Storage, f *Faculty) error {
	query := `
		UPDATE faculties
		SET name=$1
		WHERE id=$2
	`

	res, err := storage.DB().Exec(query, f.Name, f.ID)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("faculty not found")
	}

	return nil
}

// -------------------- DELETE --------------------

func DeleteFaculty(storage *postgres.Storage, id int64) error {
	res, err := storage.DB().Exec("DELETE FROM faculties WHERE id=$1", id)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("faculty not found")
	}

	return nil
}
