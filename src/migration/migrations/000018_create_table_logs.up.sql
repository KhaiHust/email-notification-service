CREATE TABLE email_logs (
                            id BIGSERIAL PRIMARY KEY,
                            email_request_id BIGINT NOT NULL,                -- Ties back to an email request
                            request_id varchar(255) NOT NULL,                       -- Ties back to an email request
                            workspace_id bigint NOT NULL,                     -- For multi-tenant systems
                            template_id BIGINT NOT NULL,                    -- Which email template was used
                            recipient VARCHAR(255) NOT NULL,          -- Email address of the recipient
                            status VARCHAR(32) NOT NULL,                    -- queued, sent, failed, opened, clicked
                            error_message TEXT,                             -- Error if sending failed
                            email_provider_id BIGINT,                           -- gmail, outlook, etc.
                            retry_count INT DEFAULT 0,                      -- How many times retried
                            logged_at TIMESTAMPTZ DEFAULT now(),
                            created_at TIMESTAMPTZ DEFAULT now(),
                            updated_at TIMESTAMPTZ DEFAULT now()
);
