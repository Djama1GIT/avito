CREATE TABLE segments
(
    id serial PRIMARY KEY,
    slug varchar(255) UNIQUE
);

CREATE TABLE users
(
    id serial PRIMARY KEY
);

CREATE TABLE user_segments
(
    user_id integer REFERENCES users(id),
    segment_id integer REFERENCES segments(id),
    PRIMARY KEY (user_id, segment_id)
);
