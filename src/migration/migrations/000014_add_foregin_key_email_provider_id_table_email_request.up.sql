ALTER TABLE email_requests
ADD CONSTRAINT fk_email_provider_id
FOREIGN KEY (email_provider_id) REFERENCES email_providers(id);