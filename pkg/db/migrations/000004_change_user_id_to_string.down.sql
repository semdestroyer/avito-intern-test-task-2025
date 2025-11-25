-- Down migration to change user ID back to integer

-- This is a complex operation and might not be fully reversible
-- In a real production environment, this would require careful planning

-- Add back the old integer columns
ALTER TABLE users ADD COLUMN id SERIAL;
ALTER TABLE pull_requests ADD COLUMN author_id_old SERIAL;
ALTER TABLE assigned_reviewers ADD COLUMN reviewer_id_old SERIAL;

-- Create mapping from user_id to id
UPDATE users SET id = CAST(SUBSTRING(user_id FROM 2) AS INTEGER);

-- Update pull_requests to populate author_id_old from users table
UPDATE pull_requests pr
SET author_id_old = u.id
FROM users u
WHERE pr.author_id = u.user_id;

-- Update assigned_reviewers to populate reviewer_id_old from users table
UPDATE assigned_reviewers ar
SET reviewer_id_old = u.id
FROM users u
WHERE ar.reviewer_id = u.user_id;

-- Drop foreign key constraints
ALTER TABLE pull_requests DROP CONSTRAINT pull_requests_author_user_id_fkey;
ALTER TABLE assigned_reviewers DROP CONSTRAINT assigned_reviewers_reviewer_user_id_fkey;

-- Drop the user_id primary key constraint
ALTER TABLE users DROP CONSTRAINT users_pkey;

-- Make id the primary key again
ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

-- Update pull_requests to use author_id_old as the main column
ALTER TABLE pull_requests RENAME COLUMN author_id TO author_id_new;
ALTER TABLE pull_requests RENAME COLUMN author_id_old TO author_id;

-- Update assigned_reviewers to use reviewer_id_old as the main column
ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_id TO reviewer_id_new;
ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_id_old TO reviewer_id;

-- Add foreign key constraints back
ALTER TABLE pull_requests ADD CONSTRAINT pull_requests_author_id_fkey 
    FOREIGN KEY (author_id) REFERENCES users(id);
ALTER TABLE assigned_reviewers ADD CONSTRAINT assigned_reviewers_reviewer_id_fkey 
    FOREIGN KEY (reviewer_id) REFERENCES users(id);

-- Drop the new columns
ALTER TABLE users DROP COLUMN user_id;
ALTER TABLE pull_requests DROP COLUMN author_id_new;
ALTER TABLE assigned_reviewers DROP COLUMN reviewer_id_new;