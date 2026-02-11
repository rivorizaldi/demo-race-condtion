CREATE SCHEMA IF NOT EXISTS demo;

CREATE TABLE demo.history (
    id VARCHAR PRIMARY KEY,
    result VARCHAR
);

CREATE TABLE demo.child (
    id VARCHAR PRIMARY KEY,
    history_id VARCHAR NOT NULL,
    result VARCHAR,
    CONSTRAINT fk_history
        FOREIGN KEY(history_id)
        REFERENCES demo.history(id)
);

CREATE INDEX idx_history_id
ON demo.child(history_id);


INSERT INTO demo.history (id, result)
VALUES ('H1', NULL);

INSERT INTO demo.child (id, history_id, result)
VALUES 
('S1', 'H1', NULL),
('S2', 'H1', NULL);