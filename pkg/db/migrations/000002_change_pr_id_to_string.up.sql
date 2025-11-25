ALTER TABLE assigned_reviewers DROP CONSTRAINT assigned_reviewers_pull_request_id_fkey;
ALTER TABLE pull_requests DROP CONSTRAINT pull_requests_pkey;
ALTER TABLE pull_requests ALTER COLUMN id TYPE VARCHAR(50) USING id::VARCHAR(50);
ALTER TABLE pull_requests ADD PRIMARY KEY (id);
ALTER TABLE assigned_reviewers ALTER COLUMN pull_request_id TYPE VARCHAR(50) USING pull_request_id::VARCHAR(50);
ALTER TABLE assigned_reviewers ADD CONSTRAINT assigned_reviewers_pull_request_id_fkey FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id);