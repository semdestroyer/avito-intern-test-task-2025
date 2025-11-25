DELETE FROM assigned_reviewers WHERE pull_request_id = 'pr-1001';
DELETE FROM pull_requests WHERE id = 'pr-1001';
DELETE FROM users WHERE username IN ('alice', 'bob', 'charlie', 'david', 'eve', 'frank', 'grace', 'henry');
DELETE FROM teams WHERE name IN ('backend', 'frontend', 'mobile');