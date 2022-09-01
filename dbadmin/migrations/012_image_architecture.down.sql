DROP INDEX image_digest_idx;
ALTER TABLE image ADD CONSTRAINT image_digest_key UNIQUE(digest);
ALTER TABLE image DROP COLUMN arch_id;
DROP TABLE arch;
