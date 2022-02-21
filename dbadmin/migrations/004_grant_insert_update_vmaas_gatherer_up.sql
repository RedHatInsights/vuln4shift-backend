GRANT SELECT, INSERT, UPDATE, DELETE ON cve TO vmaas_gatherer;
GRANT SELECT, DELETE ON account_cve_cache TO vmaas_gatherer;
GRANT SELECT, DELETE ON image_cve TO vmaas_gatherer;
GRANT SELECT, DELETE ON cluster_cve_cache TO vmaas_gatherer;