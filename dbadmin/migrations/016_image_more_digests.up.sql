ALTER TABLE image RENAME COLUMN digest TO manifest_list_digest;
ALTER INDEX image_digest_idx RENAME TO image_manifest_list_digest_idx;
ALTER TABLE image RENAME CONSTRAINT image_digest_check TO image_manifest_list_digest_check;
ALTER TABLE image ALTER COLUMN manifest_list_digest DROP NOT NULL;

ALTER TABLE image ADD COLUMN manifest_schema2_digest TEXT CHECK (NOT empty(manifest_schema2_digest));
CREATE INDEX ON image(manifest_schema2_digest);

ALTER TABLE image ADD COLUMN docker_image_digest TEXT CHECK (NOT empty(docker_image_digest));
CREATE INDEX ON image(docker_image_digest);
