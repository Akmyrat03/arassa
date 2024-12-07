CREATE TABLE IF NOT EXISTS cat_translate (
    id SERIAL PRIMARY KEY,
    cat_id INT REFERENCES categories (id) ON DELETE CASCADE,
    lang_id INT REFERENCES languages (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL
);