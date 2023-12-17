CREATE TABLE riddles (
    id SERIAL PRIMARY KEY,
    riddle TEXT NOT NULL,
    solution VARCHAR(255) NOT NULL,
    synonyms TEXT,
    published BOOLEAN DEFAULT false,
    username VARCHAR(255) DEFAULT NULL,
    user_email VARCHAR(255) DEFAULT NULL,
    date_created TIMESTAMP DEFAULT NOW(),
    last_modified TIMESTAMP DEFAULT NOW()
)

UPDATE riddles
SET published = true;

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.last_modified = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_last_modified
BEFORE UPDATE ON riddles
FOR EACH ROW
EXECUTE PROCEDURE update_modified_column();
