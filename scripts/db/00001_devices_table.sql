-- +goose Up
-- SQL in this section is executed when the migration is applied. 
CREATE TABLE devices (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(50) NOT NULL,
  "model" varchar(10) NOT NULL,
  "user" integer NOT NULL,
  "created_at" date,
  "updated_at" date,
  "deleted_at" date,
  "exchanged" boolean DEFAULT FALSE
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE devices;
