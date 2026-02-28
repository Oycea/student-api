package users

import (
	"database/sql"
	"errors"
	"time"

	"student-api/internal/storage/postgres"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"` // только при создании
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// -------------------- CREATE --------------------

func CreateUser(storage *postgres.Storage, u *User) error {
	query := `
		INSERT INTO users (username, password_hash, full_name, role_id, email, created_at)
		VALUES ($1,$2,$3,$4,$5,NOW())
		RETURNING id, created_at
	`
	return storage.DB().QueryRow(
		query,
		u.Username,
		u.Password, // пока без хэширования
		u.FullName,
		u.RoleID,
		u.Email,
	).Scan(&u.ID, &u.CreatedAt)
}

// -------------------- READ ONE --------------------

func GetUserByID(storage *postgres.Storage, id int64) (*User, error) {
	var u User

	query := `
		SELECT id, username, full_name, role_id, email, created_at
		FROM users
		WHERE id=$1
	`

	err := storage.DB().QueryRow(query, id).Scan(
		&u.ID,
		&u.Username,
		&u.FullName,
		&u.RoleID,
		&u.Email,
		&u.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

// -------------------- READ ALL --------------------

func GetAllUsers(storage *postgres.Storage) ([]User, error) {
	rows, err := storage.DB().Query(`
		SELECT id, username, full_name, role_id, email, created_at
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.FullName,
			&u.RoleID,
			&u.Email,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// -------------------- UPDATE --------------------

func UpdateUser(storage *postgres.Storage, u *User) error {
	query := `
		UPDATE users
		SET username=$1,
		    full_name=$2,
		    role_id=$3,
		    email=$4
		WHERE id=$5
	`
	res, err := storage.DB().Exec(
		query,
		u.Username,
		u.FullName,
		u.RoleID,
		u.Email,
		u.ID,
	)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// -------------------- DELETE --------------------

func DeleteUser(storage *postgres.Storage, id int64) error {
	res, err := storage.DB().Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return errors.New("user not found")
	}

	return nil
}
