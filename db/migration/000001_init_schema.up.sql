CREATE TABLE "url" (
  "id" uuid UNIQUE NOT NULL,
  "url" varchar(255) UNIQUE NOT NULL,
  "created_at" timestamp(3) NOT NULL DEFAULT (now() at time zone 'UTC'),
  PRIMARY KEY ("id")
);
