package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

type PriceTick struct {
	ID           uuid.UUID `db:"id"`
	Price        float64   `db:"price"`
	TokenIn      string    `db:"token_in"`
	AmountIn     float64   `db:"amount_in"`
	TokenOut     string    `db:"token_out"`
	AmountOut    float64   `db:"amount_out"`
	Chain        string    `db:"chain"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (e PriceTick) String() string {
	je, _ := json.Marshal(e)
	return string(je)
}

type PriceTicks []PriceTick

func (c PriceTicks) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

func (e *PriceTick) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (e *PriceTick) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

func (e *PriceTick) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}