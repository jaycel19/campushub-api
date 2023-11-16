CREATE TABLE IF NOT EXISTS profiles (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "username" varchar NOT NULL,
  "age" varchar NOT NULL,
  "program" varchar NOT NULL,
  "year" varchar NOT NULL,
  "profile_background" varchar NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
