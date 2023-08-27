CREATE TABLE segments
(
    slug varchar(255) PRIMARY KEY
);

CREATE TABLE users
(
    id integer PRIMARY KEY
);

CREATE TABLE user_segments
(
    user_id integer REFERENCES users(id),
    segment varchar(255) REFERENCES segments(slug),
    PRIMARY KEY (user_id, segment)
);
