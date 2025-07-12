ALTER TABLE email_requests
    ADD COLUMN correlation_id VARCHAR(255) DEFAULT NULL;