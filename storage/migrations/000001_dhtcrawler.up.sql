CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS shares (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   name VARCHAR NOT NULL,
   shareSize NUMERIC NOT NULL,
   fileTree VARCHAR NOT NULL,
   magnetLink VARCHAR NOT NULL
);
