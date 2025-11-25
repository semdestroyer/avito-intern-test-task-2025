ALTER TABLE users ADD COLUMN user_id VARCHAR(50) UNIQUE;

UPDATE users SET user_id = 'u' || id::TEXT;


ALTER TABLE pull_requests ADD COLUMN author_user_id VARCHAR(50);


UPDATE pull_requests pr
SET author_user_id = u.user_id
FROM users u
WHERE pr.author_id = u.id;

ALTER TABLE assigned_reviewers ADD COLUMN reviewer_user_id VARCHAR(50);


UPDATE assigned_reviewers ar
SET reviewer_user_id = u.user_id
FROM users u
WHERE ar.reviewer_id = u.id;


ALTER TABLE pull_requests DROP CONSTRAINT pull_requests_author_id_fkey;
ALTER TABLE assigned_reviewers DROP CONSTRAINT assigned_reviewers_reviewer_id_fkey;

ALTER TABLE users DROP CONSTRAINT users_pkey;

ALTER TABLE users ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);

ALTER TABLE pull_requests ADD CONSTRAINT pull_requests_author_user_id_fkey 
    FOREIGN KEY (author_user_id) REFERENCES users(user_id);
ALTER TABLE assigned_reviewers ADD CONSTRAINT assigned_reviewers_reviewer_user_id_fkey 
    FOREIGN KEY (reviewer_user_id) REFERENCES users(user_id);

ALTER TABLE pull_requests RENAME COLUMN author_id TO author_id_old;
ALTER TABLE pull_requests RENAME COLUMN author_user_id TO author_id;

ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_id TO reviewer_id_old;
ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_user_id TO reviewer_id;

ALTER TABLE pull_requests DROP COLUMN author_id_old;
ALTER TABLE assigned_reviewers DROP COLUMN reviewer_id_old;

ALTER TABLE users DROP COLUMN id;