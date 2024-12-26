CREATE TABLE IF NOT EXISTS posts(
    id UUID NOT NULL,
    userid UUID NOT NULL,
    title VARCHAR(64) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    body VARCHAR(248) NOT NULL,
    likes INT DEFAULT(0),
    comments INT DEFAULT(0),
    media_url VARCHAR(64),
    hashtag VARCHAR(10),
    created_at TIMESTAMP DEFAULT(now()::TIMESTAMP),
    modified_at TIMESTAMP DEFAULT(NULL),

    CONSTRAINT posts_id_pkey
    PRIMARY KEY(id),
    CONSTRAINT posts_userid_fkey
    FOREIGN KEY(userid) REFERENCES users(userid)
    ON DELETE CASCADE
);

