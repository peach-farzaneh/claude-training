package patterns

import (
	"errors"
	"net/http"
	"testing"
)

func TestHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "not found error",
			err:  &NotFoundError{Resource: "user", ID: "123"},
			want: http.StatusNotFound,
		},
		{
			name: "validation error",
			err:  &ValidationError{Field: "email", Message: "invalid format"},
			want: http.StatusBadRequest,
		},
		{
			name: "conflict error",
			err:  &ConflictError{Resource: "user", Message: "email already exists"},
			want: http.StatusConflict,
		},
		{
			name: "unauthorized sentinel",
			err:  ErrUnauthorized,
			want: http.StatusUnauthorized,
		},
		{
			name: "forbidden sentinel",
			err:  ErrForbidden,
			want: http.StatusForbidden,
		},
		{
			name: "wrapped not found",
			err:  fmt.Errorf("getting user: %w", &NotFoundError{Resource: "user", ID: "456"}),
			want: http.StatusNotFound,
		},
		{
			name: "wrapped unauthorized",
			err:  fmt.Errorf("checking auth: %w", ErrUnauthorized),
			want: http.StatusUnauthorized,
		},
		{
			name: "unknown error",
			err:  errors.New("something unexpected"),
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HTTPStatusCode(tt.err)
			if got != tt.want {
				t.Errorf("HTTPStatusCode(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}
