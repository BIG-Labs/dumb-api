CREATE TABLE IF NOT EXISTS ticks (
    id UUID PRIMARY KEY,
    pool_address VARCHAR(42) NOT NULL,
    tick_index INTEGER NOT NULL,
    liquidity_gross TEXT NOT NULL,
    liquidity_net TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_ticks_pool_address ON ticks(pool_address);
CREATE UNIQUE INDEX idx_ticks_pool_tick ON ticks(pool_address, tick_index); 