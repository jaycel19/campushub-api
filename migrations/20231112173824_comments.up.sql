CREATE TABLE IF NOT EXISTS comments (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (uuid_generate_v4()),
  "author" varchar NOT NULL,
  "post_id" uuid NOT NULL,
  "comment_body" varchar NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE comments ADD FOREIGN KEY ("author") REFERENCES users ("username");
ALTER TABLE comments ADD FOREIGN KEY ("post_id") REFERENCES posts ("id");
