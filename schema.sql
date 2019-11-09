CREATE TABLE "items" (
  "id" SERIAL PRIMARY KEY,
  "item_id" int,
  "guild_id" int,
  "quantity" int,
  "date_added" timestamp
);

CREATE TABLE "guild" (
  "id" SERIAL PRIMARY KEY,
  "key" varchar(50),
  "bank_char_name" varchar(60),
  "guild_name" varchar(60),
  "realm_id" int,
  "total_slots" int,
  "date_added" timestamp
);

CREATE TABLE "realms" (
  "id" SERIAL PRIMARY KEY,
  "realm_name" varchar(100)
);

ALTER TABLE "items" ADD FOREIGN KEY ("guild_id") REFERENCES "guild" ("id");

ALTER TABLE "guild" ADD FOREIGN KEY ("realm_id") REFERENCES "realms" ("id");
