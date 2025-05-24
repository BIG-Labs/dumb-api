package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)
type PoolState struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Token0    string    `json:"token0" db:"token0"`
	Token1    string    `json:"token1" db:"token1"`
	Pair      string    `json:"pair" db:"pair"`
	Factory   string    `json:"factory" db:"factory"`
	ChainID   string    `json:"chain_id" db:"chain_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Status    string    `json:"status" db:"status"`
}

// String returns a JSON representation of PoolState
func (p PoolState) String() string {
	js, _ := json.Marshal(p)
	return string(js)
}

// PoolStates is a slice of PoolState
type PoolStates []PoolState

// String returns a JSON representation of PoolStates
func (p PoolStates) String() string {
	js, _ := json.Marshal(p)
	return string(js)
}

// Validate runs validations every time "pop.Validate*" methods are called.
func (p *PoolState) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: p.ID, Name: "ID"},
		&validators.StringIsPresent{Field: p.Token0, Name: "Token0"},
		&validators.StringIsPresent{Field: p.Token1, Name: "Token1"},
		&validators.StringIsPresent{Field: p.Pair, Name: "Pair"},
		&validators.StringIsPresent{Field: p.Factory, Name: "Factory"},
		&validators.StringIsPresent{Field: p.ChainID, Name: "ChainID"},
		&validators.StringIsPresent{Field: p.Status, Name: "Status"},
	), nil
}

// ValidateCreate gets run every time "pop.ValidateAndCreate" is called.
func (p *PoolState) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time "pop.ValidateAndUpdate" is called.
func (p *PoolState) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
