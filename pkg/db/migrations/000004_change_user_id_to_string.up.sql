-- Up migration to change user ID from integer to string
-- This is a complex migration that requires multiple steps

-- First, add a new column for string user IDs
ALTER TABLE users ADD COLUMN user_id VARCHAR(50) UNIQUE;

-- Populate the new user_id column with string versions of the existing id values
UPDATE users SET user_id = 'u' || id::TEXT;

-- Add user_id column to pull_requests table
ALTER TABLE pull_requests ADD COLUMN author_user_id VARCHAR(50);

-- Update pull_requests to populate author_user_id from users table
UPDATE pull_requests pr
SET author_user_id = u.user_id
FROM users u
WHERE pr.author_id = u.id;

-- Add user_id column to assigned_reviewers table
ALTER TABLE assigned_reviewers ADD COLUMN reviewer_user_id VARCHAR(50);

-- Update assigned_reviewers to populate reviewer_user_id from users table
UPDATE assigned_reviewers ar
SET reviewer_user_id = u.user_id
FROM users u
WHERE ar.reviewer_id = u.id;

-- Update foreign key references to use user_id instead of id
-- Drop old foreign key constraints
ALTER TABLE pull_requests DROP CONSTRAINT pull_requests_author_id_fkey;
ALTER TABLE assigned_reviewers DROP CONSTRAINT assigned_reviewers_reviewer_id_fkey;

-- Drop the old primary key constraint
ALTER TABLE users DROP CONSTRAINT users_pkey;

-- Make user_id column NOT NULL
ALTER TABLE users ALTER COLUMN user_id SET NOT NULL;

-- Add primary key constraint on user_id
ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);

-- Add new foreign key constraints
ALTER TABLE pull_requests ADD CONSTRAINT pull_requests_author_user_id_fkey 
    FOREIGN KEY (author_user_id) REFERENCES users(user_id);
ALTER TABLE assigned_reviewers ADD CONSTRAINT assigned_reviewers_reviewer_user_id_fkey 
    FOREIGN KEY (reviewer_user_id) REFERENCES users(user_id);

-- Update pull_requests to use author_user_id as the main column
ALTER TABLE pull_requests RENAME COLUMN author_id TO author_id_old;
ALTER TABLE pull_requests RENAME COLUMN author_user_id TO author_id;

-- Update assigned_reviewers to use reviewer_user_id as the main column
ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_id TO reviewer_id_old;
ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_user_id TO reviewer_id;

-- Drop old columns
ALTER TABLE pull_requests DROP COLUMN author_id_old;
ALTER TABLE assigned_reviewers DROP COLUMN reviewer_id_old;

-- Drop old id column from users
ALTER TABLE users DROP COLUMN id;