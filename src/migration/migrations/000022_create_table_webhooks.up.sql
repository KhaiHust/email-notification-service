CREATE TABLE webhooks(
                         id SERIAL PRIMARY KEY,
                         workspace_id BIGINT NOT NULL,
                         name VARCHAR(100) NOT NULL,
                         type VARCHAR(20),
                         url TEXT NOT NULL,
                         enabled BOOLEAN DEFAULT TRUE,
                         created_at TIMESTAMPTZ DEFAULT now(),
                         updated_at TIMESTAMPTZ DEFAULT now()
);
AlTER TABLE webhooks
    ADD CONSTRAINT fk_workspace_id FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE;