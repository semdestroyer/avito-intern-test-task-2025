-- Migration to add test data for the PR Reviewer Assignment Service

-- Insert test teams
INSERT INTO teams (name) VALUES 
  ('backend'),
  ('frontend'),
  ('mobile')
ON CONFLICT (name) DO NOTHING;

-- Insert test users with integer IDs (migration 000004 will convert them to strings)
INSERT INTO users (username, is_active, team_name) VALUES 
  ('alice', true, 'backend'),
  ('bob', true, 'backend'),
  ('charlie', true, 'backend'),
  ('david', true, 'backend'),
  ('eve', true, 'frontend'),
  ('frank', true, 'frontend'),
  ('grace', true, 'mobile'),
  ('henry', true, 'mobile')
ON CONFLICT (username) DO NOTHING;

-- Insert a pull request with string ID
INSERT INTO pull_requests (id, pull_request_name, author_id, status_id) 
SELECT 'pr-1001', 'Add user authentication', u.id, 1
FROM users u
WHERE u.username = 'alice'
ON CONFLICT DO NOTHING;

-- Assign reviewers to the pull request
INSERT INTO assigned_reviewers (pull_request_id, reviewer_id)
SELECT 'pr-1001', u.id
FROM users u
WHERE u.username IN ('bob', 'charlie')
ON CONFLICT DO NOTHING;