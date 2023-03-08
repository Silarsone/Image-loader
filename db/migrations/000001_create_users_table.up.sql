CREATE TABLE users (
    id serial PRIMARY KEY,
    name text NOT NULL,
    description text,
    login text not null,
    password text not null
);
