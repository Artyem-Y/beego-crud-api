
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied


CREATE TYPE roles AS ENUM (
'user',
'admin',
'super-admin');

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  Name VARCHAR(64) DEFAULT NULL,
  Email VARCHAR(64) NOT NULL,
  Phone VARCHAR(64) NOT NULL,
  Password VARCHAR(128) NOT NULL,
  Role roles DEFAULT NULL,
  Access_token VARCHAR DEFAULT NULL,
  Email_validation_code VARCHAR DEFAULT NULL,
  Validation_code_sent TIMESTAMPTZ,
  Email_confirmed bool DEFAULT false,
  Created_at TIMESTAMPTZ,
  Updated_at TIMESTAMPTZ,
  Recent_login TIMESTAMPTZ
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE users;