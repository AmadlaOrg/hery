-- All the tables used in the caching

-- This table stores information about entities
CREATE TABLE IF NOT EXISTS entities (
    entity_type TEXT NOT NULL,  -- _type: entity type URI with version
    entity_extends TEXT,        -- _extends: URI to extended entity instance
    uri TEXT UNIQUE,            -- Full entity URI (type@version)
    name TEXT,                  -- Simple name of the entity
    repo_url TEXT,              -- Full URL to the repository
    origin TEXT,                -- Partial path of the entity
    version TEXT,               -- Entity version
    is_latest_version BOOLEAN,
    is_pseudo_version BOOLEAN,
    abs_path TEXT UNIQUE,       -- Full system path to the entity files
    have BOOLEAN,               -- Whether the entity is on the local machine
    hash TEXT UNIQUE,           -- Hash of the entity content for validation
    exist BOOLEAN,              -- Whether the repository was found
    schema_json TEXT,           -- JSON schema of the entity
    meta_json TEXT,             -- _meta as JSON string
    body_json TEXT,             -- _body as JSON string
    requires_json TEXT,         -- _requires as JSON string
    merged_json TEXT,           -- Full merged entity as JSON
    source_file TEXT,           -- Origin .hery file path
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_entities_type ON entities(entity_type);
CREATE INDEX IF NOT EXISTS idx_entities_name ON entities(name);
CREATE INDEX IF NOT EXISTS idx_entities_uri ON entities(uri);
CREATE INDEX IF NOT EXISTS idx_entities_repo_url ON entities(repo_url);
CREATE INDEX IF NOT EXISTS idx_entities_version ON entities(version);
CREATE INDEX IF NOT EXISTS idx_entities_is_latest_version ON entities(is_latest_version);
CREATE INDEX IF NOT EXISTS idx_entities_have ON entities(have);

-- Manifest of the source files a project cache (.hery.cache) was built from.
-- Directory rows record the sorted .hery filenames so added/removed files are
-- caught without depending on directory mtimes (which the cache file itself
-- would perturb).
CREATE TABLE IF NOT EXISTS cache_manifest (
    path TEXT PRIMARY KEY,      -- Absolute path of a source .hery file or directory
    is_dir BOOLEAN,             -- TRUE for directory rows
    mtime_ns INTEGER,           -- File modification time (ns); 0 for directories
    size INTEGER,               -- File size in bytes; 0 for directories
    hery_names TEXT             -- Directory rows: JSON array of .hery filenames
);
