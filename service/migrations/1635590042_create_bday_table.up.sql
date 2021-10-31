CREATE TABLE IF NOT EXISTS bday (
  id         SERIAL PRIMARY KEY,
  username   varchar(20) NOT NULL UNIQUE,
  bday       timestamp NOT NULL
);
