CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    category_id INT REFERENCES categories (id) ON DELETE CASCADE,
    image TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS news_translate (
    id SERIAL PRIMARY KEY,
    news_id INT REFERENCES news (id) ON DELETE CASCADE,
    lang_id INT REFERENCES languages (id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL
);


