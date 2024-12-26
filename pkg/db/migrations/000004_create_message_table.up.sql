CREATE TABLE IF NOT EXISTS messages(
    id SERIAL,
    senderid UUID,
    recieverid UUID,
    body VARCHAR(64) NOT NULL,
    status INT NOT NULL,
    time TIMESTAMP NOT NULL,

    CONSTRAINT messages_id_pkey
    PRIMARY KEY(id),
    CONSTRAINT messages_sender_fkey
    FOREIGN KEY(senderid) REFERENCES users(userid)
    ON DELETE CASCADE,
    CONSTRAINT messages_reciever_fkey
    FOREIGN KEY(recieverid) REFERENCES users(userid)
    ON DELETE CASCADE
);

