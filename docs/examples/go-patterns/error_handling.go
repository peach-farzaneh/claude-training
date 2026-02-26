// Package patterns demonstrates Go error handling patterns
// used in this project, following the conventions in CLAUDE.md.
//
// Key rules:
// - All errors must be wrapped with context using fmt.Errorf with %w
// - Prefer explicit error handling over panics
// - Use custom error types for domain errors
// - Use errors.Is / errors.As for error inspection
package patterns

import (
	"errors"
	"fmt"
	"net/http"
)

// --- Custom Error Types ---

// NotFoundError indicates a requested resource does not exist.
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
}

// ValidationError indicates invalid input from the client.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Message)
}

// ConflictError indicates a resource state conflict (e.g. duplicate).
type ConflictError struct {
	Resource string
	Message  string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict on %s: %s", e.Resource, e.Message)
}

// --- Sentinel Errors ---

// ErrUnauthorized is returned when authentication is missing or invalid.
var ErrUnauthorized = errors.New("unauthorized")

// ErrForbidden is returned when the caller lacks permission.
var ErrForbidden = errors.New("forbidden")

// --- Error to HTTP Status Mapping ---

// HTTPStatusCode maps a domain error to the appropriate HTTP status code.
func HTTPStatusCode(err error) int {
	var notFound *NotFoundError
	var validation *ValidationError
	var conflict *ConflictError

	switch {
	case errors.As(err, &notFound):
		return http.StatusNotFound
	case errors.As(err, &validation):
		return http.StatusBadRequest
	case errors.As(err, &conflict):
		return http.StatusConflict
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// --- Error Wrapping Examples ---

// Example_wrapping shows the correct way to wrap errors at each layer.
//
// Store layer:
//
//	func (s *Store) GetUser(ctx context.Context, id string) (*User, error) {
//	    row := s.db.QueryRowContext(ctx, "SELECT ... WHERE id = $1", id)
//	    var u User
//	    if err := row.Scan(&u.ID, &u.Name); err != nil {
//	        if errors.Is(err, sql.ErrNoRows) {
//	            return nil, &NotFoundError{Resource: "user", ID: id}
//	        }
//	        return nil, fmt.Errorf("querying user %s: %w", id, err)
//	    }
//	    return &u, nil
//	}
//
// Service layer:
//
//	func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
//	    user, err := s.store.GetUser(ctx, id)
//	    if err != nil {
//	        return nil, fmt.Errorf("getting user: %w", err)
//	    }
//	    return user, nil
//	}
//
// Handler layer:
//
//	func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
//	    user, err := h.service.GetUser(r.Context(), chi.URLParam(r, "id"))
//	    if err != nil {
//	        status := HTTPStatusCode(err)
//	        http.Error(w, err.Error(), status)
//	        return
//	    }
//	    json.NewEncoder(w).Encode(user)
//	}
func Example_wrapping() {}
