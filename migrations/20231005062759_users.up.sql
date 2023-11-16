CREATE TABLE IF NOT EXISTS users (
  "username" varchar PRIMARY KEY NOT NULL UNIQUE,
  "email" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
