package storage

import "errors"

// Общие ошибки для работы с хранилищем
var (
	ErrUnknown = errors.New("Unknown DB error")
)
