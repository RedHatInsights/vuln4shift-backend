ALTER TABLE image DROP COLUMN docker_image_digest;

ALTER TABLE image DROP COLUMN manifest_schema2_digest;

ALTER TABLE image ALTER COLUMN manifest_list_digest SET NOT NULL;
ALTER TABLE image RENAME CONSTRAINT image_manifest_list_digest_check TO image_digest_check;
ALTER INDEX image_manifest_list_digest_idx RENAME TO image_digest_idx;
ALTER TABLE image RENAME COLUMN manifest_list_digest TO digest;
