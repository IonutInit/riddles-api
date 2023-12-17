CREATE TABLE images (
    id SERIAL PRIMARY KEY,
    riddleId INT,
    image TEXT,
    date_created TIMESTAMP DEFAULT NOW()
)