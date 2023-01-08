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
	Text string
	Err  error
}

func (e *Error) Error() string {
	return e.Text
}

func (e *Error) Unwrap() error {
	return e.Err
}

type NotFoundError struct {
	Text string
	Err  error
}

func (e *NotFoundError) Error() string {
	return e.Text
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

type ValidationError struct {
	Text string
	Err  error
}

func (e *ValidationError) Error() string {
	return e.Text
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

type ConflictError struct {
	Text string
	Err  error
}

func (e *ConflictError) Error() string {
	return e.Text
}

func (e *ConflictError) Unwrap() error {
	return e.Err
}

var (
	ErrNotFound   *NotFoundError   = &NotFoundError{Text: "Target was not found"}
	ErrValidation *ValidationError = &ValidationError{Text: "Provided data is not valid"}
	ErrConflict   *ConflictError   = &ConflictError{Text: "ID exists"}
)
