CREATE TABLE images (
    id serial PRIMARY KEY,
    user_id int4 REFERENCES users(id),
    name text NOT NULL,
    extension text not null
)