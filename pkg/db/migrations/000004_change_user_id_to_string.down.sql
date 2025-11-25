

ALTER TABLE users ADD COLUMN id SERIAL;
ALTER TABLE pull_requests ADD COLUMN author_id_old SERIAL;
ALTER TABLE assigned_reviewers ADD COLUMN reviewer_id_old SERIAL;


UPDATE users SET id = CAST(SUBSTRING(user_id FROM 2) AS INTEGER);


UPDATE pull_requests pr
SET author_id_old = u.id
FROM users u
WHERE pr.author_id = u.user_id;


UPDATE assigned_reviewers ar
SET reviewer_id_old = u.id
FROM users u
WHERE ar.reviewer_id = u.user_id;


ALTER TABLE pull_requests DROP CONSTRAINT pull_requests_author_user_id_fkey;
ALTER TABLE assigned_reviewers DROP CONSTRAINT assigned_reviewers_reviewer_user_id_fkey;


ALTER TABLE users DROP CONSTRAINT users_pkey;


ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);


ALTER TABLE pull_requests RENAME COLUMN author_id TO author_id_new;
ALTER TABLE pull_requests RENAME COLUMN author_id_old TO author_id;


ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_id TO reviewer_id_new;
ALTER TABLE assigned_reviewers RENAME COLUMN reviewer_id_old TO reviewer_id;


ALTER TABLE pull_requests ADD CONSTRAINT pull_requests_author_id_fkey 
    FOREIGN KEY (author_id) REFERENCES users(id);
ALTER TABLE assigned_reviewers ADD CONSTRAINT assigned_reviewers_reviewer_id_fkey 
    FOREIGN KEY (reviewer_id) REFERENCES users(id);


ALTER TABLE users DROP COLUMN user_id;
ALTER TABLE pull_requests DROP COLUMN author_id_new;
ALTER TABLE assigned_reviewers DROP COLUMN reviewer_id_new;