CREATE TABLE IF NOT EXISTS segment
(
    name VARCHAR PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_segment
(
    user_id INT NOT NULL,
    segment_name VARCHAR REFERENCES segment (name) ON UPDATE CASCADE,
    ttl TIMESTAMP NOT NULL DEFAULT '9999-01-01 00:00:00',
    CONSTRAINT user_segment_pkey PRIMARY KEY (user_id, segment_name)
);

CREATE TABLE IF NOT EXISTS user_segment_history
(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    segment_name VARCHAR REFERENCES segment (name) ON UPDATE CASCADE,
    operation VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);