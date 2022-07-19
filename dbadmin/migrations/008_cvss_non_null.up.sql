UPDATE cve SET cvss2_score = 0.0 WHERE cvss2_score IS NULL;
ALTER TABLE cve ALTER COLUMN cvss2_score SET NOT NULL;

UPDATE cve SET cvss3_score = 0.0 WHERE cvss3_score IS NULL;
ALTER TABLE cve ALTER COLUMN cvss3_score SET NOT NULL;
