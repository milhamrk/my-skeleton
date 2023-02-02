BEGIN;

CREATE TABLE users (
       id TEXT NOT NULL PRIMARY KEY,
       email TEXT NOT NULL UNIQUE,
       name TEXT NOT NULL,
       surname TEXT NOT NULL,
       patronymic TEXT NOT NULL,
       role TEXT NOT NULL,
       encrypted_password TEXT NOT NULL
);

COMMIT;