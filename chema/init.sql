CREATE TABLE segments
(
    slug varchar(255) PRIMARY KEY,
    percent integer 
);


CREATE TABLE user_segments
(
    user_id integer NOT NULL,
    segment varchar(255) NOT NULL REFERENCES segments(slug) ON DELETE CASCADE,
    expiration_time timestamp,
    PRIMARY KEY (user_id, segment)
);

CREATE TABLE user_segments_history
(
    user_id integer NOT NULL,
    segment varchar(255) NOT NULL,
    operation boolean NOT NULL,
    operation_datetime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);
