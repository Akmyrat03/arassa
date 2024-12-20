CREATE TABLE IF NOT EXISTS video_titles (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS video_title_translation (
    id SERIAL PRIMARY KEY,
    video_title_id INT REFERENCES video_titles (id) ON DELETE CASCADE,
    lang_id INT REFERENCES languages (id) ON DELETE CASCADE,
    video_title VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS videos (
    id SERIAL PRIMARY KEY,
    video_title_id INT REFERENCES video_titles (id) ON DELETE CASCADE,
    video_path VARCHAR(255) NOT NULL
);