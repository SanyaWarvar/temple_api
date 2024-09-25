CREATE TABLE users (
    id varchar(62) NOT NULL UNIQUE,
    email varchar(62) NOT NULL UNIQUE,
    username varchar(62) NOT NULL UNIQUE,
    password_hash varchar(62) NOT NULL,
    confirmed_email boolean NOT NULL DEFAULT false
);
