ALTER TABLE account DROP COLUMN org_id;
ALTER TABLE account DROP COLUMN account_number;
ALTER TABLE account ADD COLUMN name TEXT NOT NULL UNIQUE CHECK (NOT empty(name));
