CREATE TABLE IF NOT EXISTS title (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS title_translate (
    id SERIAL PRIMARY KEY,
    title_id INT REFERENCES title (id) ON DELETE CASCADE,
    lang_id INT REFERENCES languages (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS images (
    id SERIAL PRIMARY KEY,
    title_id INT REFERENCES title (id) ON DELETE CASCADE,
    image_path VARCHAR(255) NOT NULL
);