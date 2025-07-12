CREATE TABLE email_providers (
                                      id            serial PRIMARY KEY,
                                      workspace_id       BIGINT NOT NULL,  -- Reference to user
                                      provider      VARCHAR(50) NOT NULL,
                                      smtp_host     VARCHAR(255) NOT NULL,
                                      smtp_port     INT NOT NULL,
                                      email         VARCHAR(255) NOT NULL UNIQUE,
                                      password      TEXT, -- Encrypted if using password-based SMTP
                                      oauth_token   TEXT, -- Stores OAuth token for Gmail/Outlook
                                      oauth_refresh_token TEXT, -- Refresh token for long-term use
                                      oauth_expires_at TIMESTAMP, -- Token expiry time
                                      use_tls       BOOLEAN DEFAULT TRUE,
                                      created_at    TIMESTAMP DEFAULT NOW(),
                                      updated_at    TIMESTAMP DEFAULT NOW(),
                                      FOREIGN KEY (workspace_id) REFERENCES workspaces(id)
);
