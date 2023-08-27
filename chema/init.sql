CREATE TABLE segments
(
    slug varchar(255) PRIMARY KEY
);


CREATE TABLE user_segments
(
    user_id integer,
    segment varchar(255) REFERENCES segments(slug) ON DELETE CASCADE,
    PRIMARY KEY (user_id, segment)
);
