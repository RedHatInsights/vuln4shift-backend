CREATE COLLATION IF NOT EXISTS numeric (provider = icu, locale = 'en-u-kn-true');
ALTER TABLE cluster ALTER COLUMN version SET DATA TYPE TEXT COLLATE "numeric";
