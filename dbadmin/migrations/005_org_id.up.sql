ALTER TABLE account DROP COLUMN name;
ALTER TABLE account ADD COLUMN account_number TEXT NOT NULL UNIQUE CHECK (NOT empty(account_number));
ALTER TABLE account ADD COLUMN org_id TEXT NOT NULL UNIQUE CHECK (NOT empty(org_id));
