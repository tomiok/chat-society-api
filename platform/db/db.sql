CREATE TABLE rooms
(
    ID          VARCHAR(255) NOT NULL,
    Title       VARCHAR(255) UNIQUE,
    Description TEXT,
    Owner       VARCHAR(255),
    IsModerated BOOLEAN,
    Moderator   VARCHAR(255),
    IsOnlyAudio BOOLEAN,
    IsOnlyText  BOOLEAN,
    IsBoth      BOOLEAN,
    Max         INT,
    URL         VARCHAR(255),
    CreatedAt   TIMESTAMP,
    PRIMARY KEY (ID)
);

CREATE TABLE participants (
                              id int auto_increment NOT NULL,
                              nick varchar(250) NOT NULL UNIQUE,
                              password varchar(100) NOT NULL,
                              created_at DATETIME DEFAULT NOW() NULL,
                              CONSTRAINT participants_PK PRIMARY KEY (id)
);
-- docker run --name=tl-mysql -p3306:3306 -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=root -d mysql/mysql-server:8.0.20