CREATE TABLE IF NOT EXISTS comments(
    id UUID,
    postid UUID,
    body VARCHAR(64),

    CONSTRAINT comments_id_pkey
    PRIMARY KEY(id),
    CONSTRAINT comments_postid_fkey
    FOREIGN KEY(postid) REFERENCES posts(id)
    ON DELETE CASCADE
);