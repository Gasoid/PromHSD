package db

import "fmt"

type StorageError struct {
	Text string
	Err  error
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("Storage error: %s", e.Text)
}

func (e *StorageError) Unwrap() error {
	return e.Err
}

type Error struct {
	text string
	err  error
}

func (e *Error) Error() string {
	return e.text
}

func (e *Error) Unwrap() error {
	return e.err
}

var (
	ErrNotFound   *Error = &Error{text: "Target was not found"}
	ErrValidation *Error = &Error{text: "Provided data is not valid"}
)
