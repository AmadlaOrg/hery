-- This table stores information about the tables different columns
CREATE TABLE IF NOT EXISTS schema_metadata (
   table_name TEXT,
   column_name TEXT,
   description TEXT
);

-- This table stores information about entities
CREATE TABLE IF NOT EXISTS entities (
    _entity TEXT,  -- Contains the entities URI (not unique)
    Name TEXT,     -- Just the name of the entity (not unique)
    RepoUrl TEXT,  -- The full URL to the repository containing the entity
    Origin TEXT,   -- Contains the partial path of the entity
    Version TEXT,  -- The entity version
    IsLatestVersion BOOLEAN,  -- Whether the entity version is the latest
    AbsPath TEXT UNIQUE, -- The full system path to the entity files (unique)
    Have BOOLEAN,  -- Indicates if the entity is on the local machine
    Hash TEXT,     -- The hash of the entity content for validation
    Exist BOOLEAN, -- Indicates if the repository was found
    Schema TEXT,   -- JSON schema of the entity
    insert_date_time DATETIME,  -- The date when the row was added
    update_date_time DATETIME   -- The latest date when the row was updated
);

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
    body TEXT,  -- Contains the content from _body
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
    PRIMARY KEY (entities_rowid, _id),
    FOREIGN KEY (entities_rowid) REFERENCES entities(rowid)  -- Explicit foreign key constraint
);

CREATE TABLE IF NOT EXISTS meta (
    entities_rowid INTEGER, -- The rowid from the entities table
    body_rowid INTEGER, --
    insert_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The date when this row was added
    update_date_time DATETIME DEFAULT CURRENT_TIMESTAMP,  -- The latest date when this row was updated
    PRIMARY KEY (entities_rowid, body_rowid),
    FOREIGN KEY (entities_rowid) REFERENCES entities(rowid)  -- Explicit foreign key constraint
);
