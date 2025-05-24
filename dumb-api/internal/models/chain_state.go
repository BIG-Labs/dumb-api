package models

import (
	"encoding/json"
	"time"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type ChainState struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ChainID   string    `json:"chain_id" db:"chain_id"`
	LastBlock int64     `json:"last_block" db:"last_block"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (cs ChainState) TableName() string {
	return "chain_states"
}

func (cs ChainState) String() string {
	je, _ := json.Marshal(cs)
	return string(je)
}

type ChainStates []ChainState

func (cs *ChainState) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (cs *ChainState) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (cs *ChainState) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
