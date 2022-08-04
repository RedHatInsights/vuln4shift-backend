DROP USER cleaner;

REVOKE DELETE ON cluster IN SCHEMA public FROM cleaner;
REVOKE DELETE ON cluster_cve_cache IN SCHEMA public FROM cleaner;
REVOKE DELETE ON cluster_image IN SCHEMA public FROM cleaner;
