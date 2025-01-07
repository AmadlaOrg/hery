-- This table stores information about the tables different columns
CREATE TABLE IF NOT EXISTS schema_metadata (
   table_name TEXT,
   column_name TEXT,
   description TEXT
);

CREATE INDEX IF NOT EXISTS idx_schema_metadata_table_name_and_column_name ON schema_metadata(table_name, column_name);

-- This table stores information about entities
CREATE TABLE IF NOT EXISTS entities (
    _entity TEXT UNIQUE,  -- Contains the entities URI (not unique)
    Name TEXT,     -- Just the name of the entity (not unique)
    RepoUrl TEXT,  -- The full URL to the repository containing the entity
    Origin TEXT,   -- Contains the partial path of the entity
    Version TEXT,  -- The entity version
    IsLatestVersion BOOLEAN,  -- Whether the entity version is the latest
    AbsPath TEXT UNIQUE, -- The full system path to the entity files (unique)
    Have BOOLEAN,  -- Indicates if the entity is on the local machine
    Hash TEXT UNIQUE,     -- The hash of the entity content for validation
    Exist BOOLEAN, -- Indicates if the repository was found
    Schema TEXT,   -- JSON schema of the entity
    insert_date_time DATETIME,  -- The date when the row was added
    update_date_time DATETIME   -- The latest date when the row was updated
);

-- Indexes for frequently queried columns
CREATE INDEX IF NOT EXISTS idx_entities_name ON entities(Name);
CREATE INDEX IF NOT EXISTS idx_entities_repo_url ON entities(RepoUrl);
CREATE INDEX IF NOT EXISTS idx_entities_version ON entities(Version);
CREATE INDEX IF NOT EXISTS idx_entities_exist ON entities(Exist);
CREATE INDEX IF NOT EXISTS idx_entities_is_latest_version ON entities(IsLatestVersion);
CREATE INDEX IF NOT EXISTS idx_entities_have ON entities(Have);

-- This table stores data for entities
CREATE TABLE IF NOT EXISTS body (
    entities_rowid INTEGER,  -- The rowid from the entities table
    _id TEXT DEFAULT (
        lower(hex(randomblob(4))
            || '-'
            || hex(randomblob(2))
            || '-4'
            || substr(hex(randomblob(2)), 2)
            || '-'
            || substr('89ab', abs(random() % 4) + 1, 1)
            || substr(hex(randomblob(2)), 2)
            || '-'
            || hex(randomblob(6)))
        ),
    body TEXT,  -- Contains the content from _body (JSON)
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
    PRIMARY KEY (entities_rowid, _id),
    FOREIGN KEY (entities_rowid) REFERENCES entities(rowid)  -- Explicit foreign key constraint
);

-- This table stores information for meta entities
CREATE TABLE IF NOT EXISTS meta (
    entities_rowid INTEGER, -- The rowid from the entities table
    for_body_rowid INTEGER, -- The rowid from the body indicates what the meta row is for which body
    body_rowid INTEGER, -- The rowid from the body
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
    PRIMARY KEY (entities_rowid, for_body_rowid, body_rowid),
    FOREIGN KEY (entities_rowid) REFERENCES entities(rowid),  -- Explicit foreign key constraint
    FOREIGN KEY (for_body_rowid) REFERENCES body(rowid),
    FOREIGN KEY (body_rowid) REFERENCES body(rowid)
);

CREATE TABLE IF NOT EXISTS body_data_TEXT (
    body_rowid INTEGER,
    property_name TEXT,
    property_value TEXT,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
    PRIMARY KEY (body_rowid, property_name),
    FOREIGN KEY (body_rowid) REFERENCES body(rowid)
);

CREATE TABLE IF NOT EXISTS body_data_INTEGER (
    body_rowid INTEGER,
    property_name TEXT,
    property_value INTEGER,
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
    PRIMARY KEY (body_rowid, property_name),
    FOREIGN KEY (body_rowid) REFERENCES body(rowid)
);

CREATE INDEX IF NOT EXISTS idx_body_data_INTEGER_property_value ON body_data_INTEGER(property_value);

CREATE TABLE IF NOT EXISTS body_data_REAL (
      body_rowid INTEGER,
      property_name TEXT,
      property_value REAL,
      insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
      update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
      PRIMARY KEY (body_rowid, property_name),
      FOREIGN KEY (body_rowid) REFERENCES body(rowid)
);

CREATE INDEX IF NOT EXISTS idx_body_data_REAL_property_value ON body_data_REAL(property_value);

CREATE TABLE IF NOT EXISTS body_data_BOOLEAN (
     body_rowid INTEGER,
     property_name TEXT,
     property_value BOOLEAN,
     insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
     update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
     PRIMARY KEY (body_rowid, property_name),
     FOREIGN KEY (body_rowid) REFERENCES body(rowid)
);

CREATE INDEX IF NOT EXISTS idx_body_data_BOOLEAN_property_value ON body_data_BOOLEAN(property_value);

