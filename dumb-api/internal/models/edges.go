package models

import (
	"encoding/json"
	"time"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// EdgeState represents the stored information about a given edge in the routing graph.
type EdgeState struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChainID   string    `json:"chain_id" db:"chain_id"`
	Token0    string    `json:"token0" db:"token0"`
	Token1    string    `json:"token1" db:"token1"`
	PoolID    string    `json:"pool_id" db:"pool_id"`
	EdgeData  []byte    `json:"edge_data" db:"edge_data"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String returns a JSON representation of EdgeState
func (e EdgeState) String() string {
	js, _ := json.Marshal(e)
	return string(js)
}

// EdgeStates is a slice of EdgeState
type EdgeStates []EdgeState

// String returns a JSON representation of EdgeStates
func (e EdgeStates) String() string {
	js, _ := json.Marshal(e)
	return string(js)
}

// Validate runs validations every time "pop.Validate*" methods are called.
func (e *EdgeState) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: e.ID, Name: "ID"},
		&validators.StringIsPresent{Field: e.ChainID, Name: "ChainID"},
		&validators.StringIsPresent{Field: e.Token0, Name: "Token0"},
		&validators.StringIsPresent{Field: e.Token1, Name: "Token1"},
		&validators.StringIsPresent{Field: e.PoolID, Name: "PoolID"},
	), nil
}

// ValidateCreate gets run every time "pop.ValidateAndCreate" is called.
func (e *EdgeState) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time "pop.ValidateAndUpdate" is called.
func (e *EdgeState) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
