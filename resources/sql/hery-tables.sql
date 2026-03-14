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

-- Body and meta data are stored as JSON in the entities table (meta_json, body_json, merged_json).
-- The separate body_data tables below are kept for potential future use with indexed property queries.

CREATE TABLE IF NOT EXISTS body_data_TEXT (
    body_rowid INTEGER,
    property_name TEXT,
    property_value TEXT,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (body_rowid, property_name)
);

CREATE INDEX IF NOT EXISTS idx_body_data_property_name ON body_data_TEXT(property_name);

CREATE TABLE IF NOT EXISTS body_data_NUMERIC (
    body_rowid INTEGER,
    property_name TEXT,
    property_value INTEGER,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (body_rowid, property_name)
);

CREATE INDEX IF NOT EXISTS idx_body_data_NUMERIC_property_value ON body_data_NUMERIC(property_value);

CREATE TABLE IF NOT EXISTS body_data_REAL (
    body_rowid INTEGER,
    property_name TEXT,
    property_value REAL,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (body_rowid, property_name)
);

CREATE INDEX IF NOT EXISTS idx_body_data_REAL_property_value ON body_data_REAL(property_value);

CREATE TABLE IF NOT EXISTS body_data_BOOLEAN (
    body_rowid INTEGER,
    property_name TEXT,
    property_value BOOLEAN,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (body_rowid, property_name)
);

CREATE INDEX IF NOT EXISTS idx_body_data_BOOLEAN_property_value ON body_data_BOOLEAN(property_value);

CREATE TABLE IF NOT EXISTS body_data_DATE (
    body_rowid INTEGER,
    property_name TEXT,
    property_value DATE,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (body_rowid, property_name)
);

CREATE INDEX IF NOT EXISTS idx_body_data_DATE_property_value ON body_data_DATE(property_value);

CREATE TABLE IF NOT EXISTS body_data_DATETIME (
    body_rowid INTEGER,
    property_name TEXT,
    property_value DATETIME,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (body_rowid, property_name)
);

CREATE INDEX IF NOT EXISTS idx_body_data_DATETIME_property_value ON body_data_DATETIME(property_value);

CREATE TABLE IF NOT EXISTS body_data_connection (
    parent_body_data_TEXT_rowid INTEGER,
    body_data_TEXT_rowid INTEGER,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (parent_body_data_TEXT_rowid, body_data_TEXT_rowid),
    FOREIGN KEY (parent_body_data_TEXT_rowid) REFERENCES body_data_TEXT(rowid),
    FOREIGN KEY (body_data_TEXT_rowid) REFERENCES body_data_TEXT(rowid)
);
