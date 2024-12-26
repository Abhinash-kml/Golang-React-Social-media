CREATE TABLE IF NOT EXISTS messages(
    id SERIAL NOT NULL,
    senderid UUID NOT NULL,
    recieverid UUID NOT NULL,
    body VARCHAR(64) NOT NULL,
    status INT NOT NULL DEFAULT(0),
    time TIMESTAMP NOT NULL DEFAULT(now()::TIMESTAMP),

    CONSTRAINT messages_id_pkey
    PRIMARY KEY(id),
    CONSTRAINT messages_sender_fkey
    FOREIGN KEY(senderid) REFERENCES users(userid)
    ON DELETE CASCADE,
    CONSTRAINT messages_reciever_fkey
    FOREIGN KEY(recieverid) REFERENCES users(userid)
    ON DELETE CASCADE
);

