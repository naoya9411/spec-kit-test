package todo

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// TodoID represents a unique identifier for a todo
type TodoID struct {
	value uuid.UUID
}

// NewTodoID creates a new TodoID
func NewTodoID() TodoID {
	return TodoID{value: uuid.New()}
}

// NewTodoIDFromString creates a TodoID from string
func NewTodoIDFromString(id string) (TodoID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return TodoID{}, fmt.Errorf("invalid todo ID: %w", err)
	}
	return TodoID{value: parsedID}, nil
}

// String returns the string representation of TodoID
func (id TodoID) String() string {
	return id.value.String()
}

// Value returns the underlying UUID value
func (id TodoID) Value() uuid.UUID {
	return id.value
}

// Equals checks if two TodoIDs are equal
func (id TodoID) Equals(other TodoID) bool {
	return id.value == other.value
}

// Scan implements the sql.Scanner interface for GORM
func (id *TodoID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch s := value.(type) {
	case string:
		parsed, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		id.value = parsed
	case []byte:
		parsed, err := uuid.Parse(string(s))
		if err != nil {
			return err
		}
		id.value = parsed
	default:
		return errors.New("cannot scan into TodoID")
	}
	
	return nil
}

// Value implements the driver.Valuer interface for GORM
func (id TodoID) DBValue() (driver.Value, error) {
	return id.value.String(), nil
}

// Title represents a todo title with validation
type Title struct {
	value string
}

// NewTitle creates a new Title with validation
func NewTitle(title string) (Title, error) {
	t := Title{value: strings.TrimSpace(title)}
	if err := t.Validate(); err != nil {
		return Title{}, err
	}
	return t, nil
}

// String returns the string representation of Title
func (t Title) String() string {
	return t.value
}

// Value returns the underlying string value
func (t Title) Value() string {
	return t.value
}

// Validate validates the title
func (t Title) Validate() error {
	if len(t.value) == 0 {
		return errors.New("title cannot be empty")
	}
	if len(t.value) > 255 {
		return errors.New("title cannot exceed 255 characters")
	}
	return nil
}

// Equals checks if two Titles are equal
func (t Title) Equals(other Title) bool {
	return t.value == other.value
}

// Scan implements the sql.Scanner interface for GORM
func (t *Title) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	
	switch s := value.(type) {
	case string:
		t.value = s
	case []byte:
		t.value = string(s)
	default:
		return errors.New("cannot scan into Title")
	}
	
	return nil
}

// DBValue implements the driver.Valuer interface for GORM
func (t Title) DBValue() (driver.Value, error) {
	return t.value, nil
}
