create table if not exists users
(
    id       serial primary key,
    username VARCHAR(50) unique not null,
    is_active boolean unique     not null default true,
    team_name VARCHAR(50) NOT NULL REFERENCES teams(name)
);

create table if not exists teams
(
    id       serial primary key,
    name VARCHAR(50) unique not null
);

create table if not exists pull_requests
(
    id              serial primary key,
    pull_request_name VARCHAR(50) unique not null,
    author_id serial NOT NULL REFERENCES users(id)
);

create table if not exists assigned_reviewers
(
    pull_request_id serial NOT NULL REFERENCES pull_requests(id),
    author_id serial NOT NULL REFERENCES users(id)
);

create table if not exists statuses
(
    id     serial primary key,
    Status VARCHAR(50) unique not null
);

insert into statuses(Status)
values ('OPENED');
insert into statuses(Status)
values ('MERGED')