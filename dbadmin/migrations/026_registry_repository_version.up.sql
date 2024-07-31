CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- store CONCAT(repository.registry, '/', repository.repository, ':', repository_image.version)
ALTER TABLE repository_image ADD COLUMN registry_repository_version TEXT;

SELECT ri.version
FROM repository_image ri
ORDER BY ri.repository_id, ri.image_id
FOR NO KEY UPDATE OF ri;

WITH cte AS (
    SELECT repository_id, image_id, CONCAT(r.registry, '/', r.repository, ':', ri.version) as registry_repository_version
    FROM repository r
    JOIN repository_image ri ON r.id = ri.repository_id
    ORDER BY ri.repository_id, ri.image_id
)
UPDATE repository_image
SET registry_repository_version = cte.registry_repository_version
FROM cte
WHERE repository_image.repository_id = cte.repository_id
    AND repository_image.image_id = cte.image_id;

ALTER TABLE repository_image ALTER COLUMN registry_repository_version SET NOT NULL;

-- GIN tri gram index for ILIKE text search
CREATE INDEX ON repository_image USING gin ( registry_repository_version gin_trgm_ops);
