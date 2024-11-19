ALTER TABLE account_cve_cache ADD UNIQUE (account_id, cve_id),
                              DROP CONSTRAINT account_cve_cache_pkey;
ALTER TABLE cluster_cve_cache ADD UNIQUE (cluster_id, cve_id),
                              DROP CONSTRAINT cluster_cve_cache_pkey;
ALTER TABLE cluster_image ADD UNIQUE (cluster_id, image_id),
                          DROP CONSTRAINT cluster_image_pkey;
ALTER TABLE image_cve ADD UNIQUE (image_id, cve_id),
                      DROP CONSTRAINT image_cve_pkey;
ALTER TABLE repository_image ADD UNIQUE (repository_id, image_id),
                             DROP CONSTRAINT repository_image_pkey;
