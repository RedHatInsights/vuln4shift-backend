REVOKE INSERT, UPDATE, DELETE ON cve IN SCHEMA public FROM vmaas_gatherer;
REVOKE DELETE ON account_cve_cache IN SCHEMA public FROM vmaas_gatherer;
REVOKE DELETE ON image_cve IN SCHEMA public FROMvmaas_gatherer;
REVOKE DELETE ON cluster_cve_cache IN SCHEMA public FROM vmaas_gatherer;
