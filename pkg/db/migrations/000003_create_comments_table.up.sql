CREATE TABLE IF NOT EXISTS comments(
    id UUID NOT NULL,
    postid UUID NOT NULL,
    body VARCHAR(64) NOT NULL,
    created_at TIMESTAMP DEFAULT(now()::TIMESTAMP),
    modified_at TIMESTAMP,

    CONSTRAINT comments_id_pkey
    PRIMARY KEY(id),
    CONSTRAINT comments_postid_fkey
    FOREIGN KEY(postid) REFERENCES posts(id)
    ON DELETE CASCADE
);