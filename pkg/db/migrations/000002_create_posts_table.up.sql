CREATE TABLE IF NOT EXISTS posts(
    id UUID,
    userid UUID,
    title VARCHAR(64),
    created_at TIMESTAMP,
    body VARCHAR(248),
    likes INT,
    comments INT,
    media_url VARCHAR(64),
    hashtag VARCHAR(10),

    CONSTRAINT posts_id_pkey
    PRIMARY KEY(id),
    CONSTRAINT posts_userid_fkey
    FOREIGN KEY(userid) REFERENCES users(userid)
    ON DELETE CASCADE
);

