CREATE TABLE videos (
    id SERIAL PRIMARY KEY,
    video_path VARCHAR(255) NOT NULL,
    uploaded_at TIMESTAMP
);

