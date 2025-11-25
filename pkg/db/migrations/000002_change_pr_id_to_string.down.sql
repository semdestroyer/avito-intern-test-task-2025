-- Down migration to change pull request ID back to integer
ALTER TABLE assigned_reviewers DROP CONSTRAINT assigned_reviewers_pull_request_id_fkey;
ALTER TABLE pull_requests DROP CONSTRAINT pull_requests_pkey;
ALTER TABLE pull_requests ALTER COLUMN id TYPE INTEGER USING id::INTEGER;
ALTER TABLE pull_requests ADD PRIMARY KEY (id);
ALTER TABLE assigned_reviewers ALTER COLUMN pull_request_id TYPE INTEGER USING pull_request_id::INTEGER;
ALTER TABLE assigned_reviewers ADD CONSTRAINT assigned_reviewers_pull_request_id_fkey FOREIGN KEY (pull_request_id) REFERENCES pull_requests(id);