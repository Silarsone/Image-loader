CREATE TABLE tg_auth (
id serial PRIMARY KEY,
user_id int4 REFERENCES users(id),
telegram_id int8
);
