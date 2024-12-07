CREATE TABLE IF NOT EXISTS motto (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    language_id INT REFERENCES languages (id) ON DELETE CASCADE,
    image_url VARCHAR(255) NOT NULL
);


