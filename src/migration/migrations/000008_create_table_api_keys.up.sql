CREATE TABLE api_keys (
                          id SERIAL PRIMARY KEY,
                          workspace_id INT NOT NULL REFERENCES workspaces(id),
                          name VARCHAR(255),
                          key_hash TEXT NOT NULL,          -- hashed version of full key
                          raw_prefix VARCHAR(16) NOT NULL, -- first part of key for fast lookup
                          environment VARCHAR(32) NOT NULL,
                          created_at TIMESTAMPTZ DEFAULT now(),
                          expires_at TIMESTAMPTZ DEFAULT NULL,
                          updated_at TIMESTAMPTZ DEFAULT now(),
                          revoked BOOLEAN DEFAULT FALSE
);
CREATE UNIQUE INDEX api_keys_workspace_id_environment_idx ON api_keys (workspace_id, environment);
