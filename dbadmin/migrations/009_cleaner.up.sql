-- user for cleaner
CREATE USER cleaner;

-- grant select/delete for cluster deletion cleaner job
GRANT SELECT, DELETE ON cluster TO cleaner;
GRANT SELECT, DELETE ON cluster_cve_cache TO cleaner;
GRANT SELECT, DELETE ON cluster_image TO cleaner;
