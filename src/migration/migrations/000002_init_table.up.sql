CREATE TABLE "email_templates" (
                                   "id" SERIAL PRIMARY KEY,
                                   "name" varchar(255),
                                   "subject" varchar(255),
                                   "body" text,
                                   "variables" json,
                                   "workspace_id" int NOT NULL UNIQUE,
                                   "created_by" int,
                                   "last_updated_by" int,
                                   created_at TIMESTAMPTZ DEFAULT now(),
                                   updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE "email_requests" (
                                  "id" SERIAL PRIMARY KEY,
                                  "template_id" int,
                                  "recipient" varchar(255),
                                  "data" json,
                                  "status" varchar(50),
                                  created_at TIMESTAMPTZ DEFAULT now(),
                                  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE "users" (
                         "id" SERIAL PRIMARY KEY,
                         "full_name" varchar(255),
                         "email" varchar(255) UNIQUE,
                         "password" text,
                         created_at TIMESTAMPTZ DEFAULT now(),
                         updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE "workspaces" (
                              "id" SERIAL PRIMARY KEY,
                              "name" varchar(255),
                              "description" text,
                            "code" varchar(255) UNIQUE,
                              created_at TIMESTAMPTZ DEFAULT now(),
                              updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE "workspace_users" (
                                   "id" SERIAL PRIMARY KEY,
                                   "user_id" bigint,
                                   "workspace_id" bigint,
                                   "role" varchar(20),
                                   created_at TIMESTAMPTZ DEFAULT now(),
                                   updated_at TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE "email_templates" ADD FOREIGN KEY ("workspace_id") REFERENCES "workspaces" ("id");

ALTER TABLE "email_templates" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "email_templates" ADD FOREIGN KEY ("last_updated_by") REFERENCES "users" ("id");

ALTER TABLE "email_requests" ADD FOREIGN KEY ("template_id") REFERENCES "email_templates" ("id");

ALTER TABLE "workspace_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "workspace_users" ADD FOREIGN KEY ("workspace_id") REFERENCES "workspaces" ("id");