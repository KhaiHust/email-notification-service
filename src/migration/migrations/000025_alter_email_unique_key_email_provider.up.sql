ALTER TABLE email_providers
DROP CONSTRAINT email_providers_email_key;

ALTER TABLE email_providers
ADD CONSTRAINT email_providers_email_key UNIQUE (email, workspace_id);