CREATE SCHEMA IF NOT EXISTS ml_analysis;

CREATE TABLE IF NOT EXISTS ml_analysis.analyzes
(
    id SERIAL PRIMARY KEY,
    query text NOT NULL,
    answer text NOT NULL,
    is_user_satisfied boolean NOT NULL
);