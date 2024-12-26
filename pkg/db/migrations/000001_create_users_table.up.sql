
CREATE TYPE BAN_LEVEL
AS
ENUM('NONE', 'TEXT', 'VOICE', 'TEXTVOICE', 'IP');

CREATE TABLE IF NOT EXISTS users(
    id SERIAL,
    userid UUID,
    name VARCHAR(32),
    email VARCHAR(64),
    password VARCHAR(150),
    dob DATE,
    created_at TIMESTAMP,
    modified_at TIMESTAMP,
    last_login TIMESTAMP,
    country VARCHAR(32),
    city VARCHAR(32),
    ban_level BAN_LEVEL,
    ban_duration TIME(6),

    CONSTRAINT users_userid_pkey
    PRIMARY KEY(userid)
);