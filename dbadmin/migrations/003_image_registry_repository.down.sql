REVOKE INSERT, UPDATE, DELETE ON image_cve IN SCHEMA public FROM pyxis_gatherer;
REVOKE INSERT ON cve IN SCHEMA public FROM pyxis_gatherer;
REVOKE INSERT, UPDATE, DELETE ON image IN SCHEMA public FROM pyxis_gatherer;

ALTER TABLE image DROP CONSTRAINT image_digest_check;
ALTER TABLE image DROP COLUMN modified_date;
ALTER TABLE image DROP COLUMN pyxis_id;

-- repository_image
DROP TABLE IF EXISTS repository_image;

-- repository
DROP TABLE IF EXISTS repository;
