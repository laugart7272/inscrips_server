CREATE TABLE "events" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "init_date" date NOT NULL,
  "end_date" date NOT NULL,
  "active" bool NOT NULL DEFAULT true,
  "created_at" timestamptz DEFAULT 'now'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "hashed_password" varchar NOT NULL,
  "phone" varchar NOT NULL,
  "user_type" varchar NOT NULL DEFAULT 'Participante',
  "avatar_path" varchar NOT NULL DEFAULT 'assets/images/user/user.png',
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz DEFAULT 'now'
);

CREATE TABLE "cientifics_works" (
  "id" bigserial PRIMARY KEY,
  "title" varchar UNIQUE NOT NULL,
  "author_id" integer NOT NULL,
  "inscription_id" integer NOT NULL,
  "resume" text NOT NULL DEFAULT 'Resumo do trabalho',
  "file" varchar NOT NULL,
  "presentation_type" varchar NOT NULL DEFAULT 'Presencial',
  "exposition_type" varchar NOT NULL DEFAULT 'Apresantaçâo',
  "created_at" timestamptz DEFAULT 'now'
);

CREATE TABLE "inscriptions" (
  "id" bigserial PRIMARY KEY,
  "user_id" integer NOT NULL,
  "event_id" integer NOT NULL,
  "created_at" timestamptz DEFAULT 'now'
);

CREATE TABLE "payment_rate" (
  "id" bigserial PRIMARY KEY,
  "event_id" integer NOT NULL,
  "user_type" varchar NOT NULL DEFAULT 'Participante',
  "presentation_type" varchar NOT NULL DEFAULT 'Presencial',
  "amount" integer NOT NULL DEFAULT 1000,
  "created_at" timestamptz DEFAULT 'now'
);

CREATE TABLE "payment" (
  "id" bigserial PRIMARY KEY,
  "user_id" integer NOT NULL,
  "event_id" integer NOT NULL,
  "inscription_id" integer UNIQUE NOT NULL,
  "user_type" varchar NOT NULL DEFAULT 'Participante',
  "presentation_type" varchar NOT NULL DEFAULT 'Presencial',
  "amount" integer NOT NULL DEFAULT 1000,
  "is_paid" bool NOT NULL DEFAULT false,
  "proof_payment_path" varchar NOT NULL DEFAULT 'assets/proof_payment',
  "created_at" timestamptz DEFAULT 'now'
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "cientifics_works" ("title");

CREATE INDEX ON "payment" ("user_id");

ALTER TABLE "cientifics_works" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id");

ALTER TABLE "cientifics_works" ADD FOREIGN KEY ("inscription_id") REFERENCES "inscriptions" ("id");

ALTER TABLE "inscriptions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "inscriptions" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");

ALTER TABLE "payment_rate" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");

ALTER TABLE "payment" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "payment" ADD FOREIGN KEY ("event_id") REFERENCES "events" ("id");

ALTER TABLE "payment" ADD FOREIGN KEY ("inscription_id") REFERENCES "inscriptions" ("id");

------------------------------------------- Feeds

INSERT INTO events (name, init_date, end_date) VALUES ('II das Jornadas Científicas do Instituto Politécnico de Saurimo', '2024-05-09', '2024-05-10');

----- User Type: ('Docente'), ('Estudante'), ('Empresa'), ('Instituições Estaduais'), ('Participante')

----- Presentation Type: ('Presencial'), ('Remota')

----- Exposition Type: ('Apresantaçâo'), ('Poster')


INSERT INTO payment_rate (event_id, user_type, presentation_type, amount) VALUES (1, 'Docente', 'Presencial', 10000), (1, 'Docente', 'Remota', 5000), (1, 'Estudante', 'Presencial', 5000), (1, 'Estudante', 'Remota', 1000), (1, 'Empresa', 'Presencial', 0), (1, 'Empresa', 'Remota', 0), (1, 'Instituições Estaduais', 'Presencial', 0), (1, 'Instituições Estaduais', 'Remota', 0), (1, 'Participante', 'Presencial', 1000), (1, 'Participante', 'Remota', 1000);

