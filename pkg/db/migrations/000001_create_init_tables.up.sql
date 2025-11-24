create table if not exists teams
(
    id   serial primary key,
    name VARCHAR(50) unique not null
);
create table if not exists users
(
    id        serial primary key,
    username  VARCHAR(50) unique not null,
    is_active boolean            not null default true,
    team_name VARCHAR(50)        not null references teams (name)
);
create table if not exists statuses
(
    id     serial primary key,
    Status VARCHAR(50) unique not null
);

INSERT INTO statuses (status)
SELECT 'OPENED'
WHERE NOT EXISTS (SELECT 1 FROM statuses WHERE status = 'OPENED');

INSERT INTO statuses (status)
SELECT 'MERGED'
WHERE NOT EXISTS (SELECT 1 FROM statuses WHERE status = 'MERGED');

create table if not exists pull_requests
(
    id                serial primary key,
    pull_request_name VARCHAR(50) unique not null,
    author_id         serial             not null references users (id),
    status_id         INTEGER NOT NULL REFERENCES statuses(id)
        DEFAULT 1
);
create table if not exists assigned_reviewers
(
    pull_request_id serial not null references pull_requests (id),
    reviewer_id       serial not null references users (id)
);

