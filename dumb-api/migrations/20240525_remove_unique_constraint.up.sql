-- Remove the unique constraint
ALTER TABLE tokens DROP CONSTRAINT IF EXISTS tokens_address_chain_id_key;

-- Add created_at column if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                  WHERE table_name = 'tokens' AND column_name = 'created_at') THEN
        ALTER TABLE tokens ADD COLUMN created_at timestamp NOT NULL DEFAULT NOW();
    END IF;
END $$; 