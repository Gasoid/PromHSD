package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageError_Error(t *testing.T) {
	type fields struct {
		Text string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test error",
			fields: fields{Text: "test", Err: errors.New("Some error")},
			want:   "Storage error: test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &StorageError{
				Text: tt.fields.Text,
				Err:  tt.fields.Err,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
			assert.Error(t, e.Unwrap())
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Text string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test error",
			fields: fields{Text: "test", Err: errors.New("Some error")},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Text: tt.fields.Text,
				Err:  tt.fields.Err,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
			assert.Error(t, e.Unwrap())
		})
	}
}

func TestNotFoundError_Error(t *testing.T) {
	type fields struct {
		Text string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test error",
			fields: fields{Text: "test", Err: errors.New("Some error")},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &NotFoundError{
				Text: tt.fields.Text,
				Err:  tt.fields.Err,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
			assert.Error(t, e.Unwrap())
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	type fields struct {
		Text string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test error",
			fields: fields{Text: "test", Err: errors.New("Some error")},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ValidationError{
				Text: tt.fields.Text,
				Err:  tt.fields.Err,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
			assert.Error(t, e.Unwrap())
		})
	}
}

func TestConflictError_Error(t *testing.T) {
	type fields struct {
		Text string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "test error",
			fields: fields{Text: "test", Err: errors.New("Some error")},
			want:   "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ConflictError{
				Text: tt.fields.Text,
				Err:  tt.fields.Err,
			}
			got := e.Error()
			assert.Equal(t, tt.want, got)
			assert.Error(t, e.Unwrap())
		})
	}
}
