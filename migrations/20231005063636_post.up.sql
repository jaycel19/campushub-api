CREATE TABLE IF NOT EXISTS posts (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "author" varchar NOT NULL,
  "image" varchar NOT NULL,
  "post_content" varchar NOT NULL,
  "likes" integer DEFAULT 0,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE posts ADD FOREIGN KEY ("author") REFERENCES users ("username");