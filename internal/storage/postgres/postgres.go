package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Storage struct {
	// Подключение к БД
	db *sql.DB
}

func New(connectStr string) (*Storage, error) {
	const op = "storage.postgres.new"

	// Пытаемся подключиться к БД по переданным аргументам
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем пинг
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Println("Connected to DB succesfully")

	return &Storage{db: db}, nil
}
