-- Restore the unique constraint
ALTER TABLE tokens ADD CONSTRAINT tokens_address_chain_id_key UNIQUE (address, chain_id);

-- Remove created_at column if it exists
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
              WHERE table_name = 'tokens' AND column_name = 'created_at') THEN
        ALTER TABLE tokens DROP COLUMN created_at;
    END IF;
END $$; 