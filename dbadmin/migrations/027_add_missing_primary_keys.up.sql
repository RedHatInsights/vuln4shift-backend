ALTER TABLE account_cve_cache ADD PRIMARY KEY (account_id, cve_id),
                              DROP CONSTRAINT account_cve_cache_account_id_cve_id_key;
ALTER TABLE cluster_cve_cache ADD PRIMARY KEY (cluster_id, cve_id),
                              DROP CONSTRAINT cluster_cve_cache_cluster_id_cve_id_key;
ALTER TABLE cluster_image ADD PRIMARY KEY (cluster_id, image_id),
                          DROP CONSTRAINT cluster_image_cluster_id_image_id_key;
ALTER TABLE image_cve ADD PRIMARY KEY (image_id, cve_id),
                      DROP CONSTRAINT image_cve_image_id_cve_id_key;
ALTER TABLE repository_image ADD PRIMARY KEY (repository_id, image_id),
                             DROP CONSTRAINT repository_image_repository_id_image_id_key;
