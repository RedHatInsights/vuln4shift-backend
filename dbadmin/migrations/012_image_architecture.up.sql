-- arch
CREATE TABLE IF NOT EXISTS arch
(
    id   BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE CHECK (NOT empty(name))
) TABLESPACE pg_default;

GRANT SELECT, INSERT ON arch TO pyxis_gatherer;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO archive_db_writer;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO archive_db_writer;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO pyxis_gatherer;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO pyxis_gatherer;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO vmaas_gatherer;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO vmaas_gatherer;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO cve_aggregator;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO cve_aggregator;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO manager;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO manager;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO cleaner;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO cleaner;

ALTER TABLE image ADD COLUMN arch_id BIGINT;
ALTER TABLE image ADD CONSTRAINT image_arch_id_fkey FOREIGN KEY (arch_id) REFERENCES arch (id);

ALTER TABLE image DROP CONSTRAINT image_digest_key;

CREATE INDEX ON image(digest);
