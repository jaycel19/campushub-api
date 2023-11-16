CREATE TABLE IF NOT EXISTS profiles (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "name" varchar NOT NULL,
  "username" varchar NOT NULL,
  "age" varchar NOT NULL,
  "program" varchar NOT NULL,
  "year" varchar NOT NULL,
  "profile_background" varchar NOT NULL DEFAULT 'FFFFFF7F',
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
