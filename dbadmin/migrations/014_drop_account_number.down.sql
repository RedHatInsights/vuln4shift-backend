ALTER TABLE account ADD COLUMN account_number TEXT NOT NULL UNIQUE CHECK (NOT empty(account_number));
