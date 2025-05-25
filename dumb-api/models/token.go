package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Token represents a row in the tokens table
type Token struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Address   string    `json:"address" db:"address"`
	ChainID   string    `json:"chainId" db:"chain_id"`
	Price     float64   `json:"price" db:"price"`
	Icon      string    `json:"icon" db:"icon"`
	Name      string    `json:"name" db:"name"`
	Symbol    string    `json:"symbol" db:"symbol"`
	Decimals  int       `json:"decimals" db:"decimals"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// String returns a JSON representation of Token
func (t Token) String() string {
	js, _ := json.Marshal(t)
	return string(js)
}

// Tokens is a slice of Token
type Tokens []Token

// String returns a JSON representation of Tokens
func (t Tokens) String() string {
	js, _ := json.Marshal(t)
	return string(js)
}

// Validate gets run every time you call a "pop.Validate*" method.
func (t *Token) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: t.ID, Name: "ID"},
		&validators.StringIsPresent{Field: t.Address, Name: "Address"},
		&validators.StringIsPresent{Field: t.ChainID, Name: "ChainID"},
		&validators.StringIsPresent{Field: t.Name, Name: "Name"},
		&validators.StringIsPresent{Field: t.Symbol, Name: "Symbol"},
		&validators.IntIsGreaterThan{Field: t.Decimals, Name: "Decimals", Compared: -1},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
func (t *Token) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
func (t *Token) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
