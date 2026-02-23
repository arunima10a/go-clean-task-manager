CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    user_id INTEGER REFERENCES users(id)
);

ALTER TABLE tasks ADD COLUMN category_id INTEGER REFERENCES categories(id);