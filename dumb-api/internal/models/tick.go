package models

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/gofrs/uuid"
)

type Tick struct {
	ID             uuid.UUID `json:"id" db:"id"`
	PoolAddress    string    `json:"pool_address" db:"pool_address"`
	Index          int       `json:"index" db:"tick_index"`
	LiquidityGross string    `json:"liquidity_gross" db:"liquidity_gross"` // stored as string due to big.Int
	LiquidityNet   string    `json:"liquidity_net" db:"liquidity_net"`     // stored as string due to big.Int
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// TableName returns the table name for this model
func (t *Tick) TableName() string {
	return "ticks"
}

// GetLiquidityGross returns the LiquidityGross as *big.Int
func (t *Tick) GetLiquidityGross() *big.Int {
	if t.LiquidityGross == "" {
		return big.NewInt(0)
	}
	result := new(big.Int)
	result.SetString(t.LiquidityGross, 10)
	return result
}

// GetLiquidityNet returns the LiquidityNet as *big.Int
func (t *Tick) GetLiquidityNet() *big.Int {
	if t.LiquidityNet == "" {
		return big.NewInt(0)
	}
	result := new(big.Int)
	result.SetString(t.LiquidityNet, 10)
	return result
}

// SetLiquidityGross sets the LiquidityGross from a *big.Int
func (t *Tick) SetLiquidityGross(value *big.Int) {
	if value == nil {
		t.LiquidityGross = "0"
		return
	}
	t.LiquidityGross = value.String()
}

// SetLiquidityNet sets the LiquidityNet from a *big.Int
func (t *Tick) SetLiquidityNet(value *big.Int) {
	if value == nil {
		t.LiquidityNet = "0"
		return
	}
	t.LiquidityNet = value.String()
}

// String is not required by pop and may be deleted
func (t Tick) String() string {
	jt, _ := json.Marshal(t)
	return string(jt)
}
