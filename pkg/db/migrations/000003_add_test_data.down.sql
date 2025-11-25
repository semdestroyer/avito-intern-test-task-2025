-- Down migration to remove test data

-- Remove assigned reviewers
DELETE FROM assigned_reviewers WHERE pull_request_id = 'pr-1001';

-- Remove pull request
DELETE FROM pull_requests WHERE id = 'pr-1001';

-- Remove test users
DELETE FROM users WHERE username IN ('alice', 'bob', 'charlie', 'david', 'eve', 'frank', 'grace', 'henry');

-- Remove test teams
DELETE FROM teams WHERE name IN ('backend', 'frontend', 'mobile');