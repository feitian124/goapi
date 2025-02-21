CREATE DATABASE IF NOT EXISTS `testdb` DEFAULT CHARACTER SET utf8mb4;

USE testdb;

SET NAMES utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;

DROP TRIGGER IF EXISTS update_posts_updated;
DROP VIEW IF EXISTS post_comments;
DROP TABLE IF EXISTS `hyphen-table`;
DROP TABLE IF EXISTS CamelizeTable;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS comment_stars;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS user_options;
DROP TABLE IF EXISTS users;

CREATE TABLE users
(
    id       int PRIMARY KEY AUTO_INCREMENT,
    username varchar(50) UNIQUE  NOT NULL,
    password varchar(50)         NOT NULL,
    email    varchar(150) UNIQUE NOT NULL COMMENT 'ex. user@example.com',
    created  timestamp           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated  timestamp DEFAULT '1970-01-01 08:00:00'
) COMMENT = 'Users table' AUTO_INCREMENT = 100;

CREATE TABLE user_options
(
    user_id    int PRIMARY KEY,
    show_email boolean   NOT NULL DEFAULT false,
    created    timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated    timestamp DEFAULT '1970-01-01 08:00:00',
    UNIQUE (user_id),
    CONSTRAINT user_options_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE CASCADE
) COMMENT = 'User options table';

CREATE TABLE posts
(
    id        bigint AUTO_INCREMENT,
    user_id   int          NOT NULL,
    title     varchar(180) NOT NULL DEFAULT 'Untitled',
    body      text         NOT NULL,
    post_type enum('public', 'private', 'draft') NOT NULL COMMENT 'public/private/draft',
    created   datetime     NOT NULL,
    updated   datetime,
    CONSTRAINT posts_id_pk PRIMARY KEY (id),
    CONSTRAINT posts_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE NO ACTION ON DELETE CASCADE,
    UNIQUE (user_id, title)
) COMMENT = 'Posts table';
CREATE INDEX posts_user_id_idx ON posts (id) USING BTREE;

CREATE TABLE comments
(
    id           bigint AUTO_INCREMENT,
    post_id      bigint   NOT NULL,
    user_id      int      NOT NULL,
    comment      text     NOT NULL COMMENT 'Comment\nMulti-line\r\ncolumn\rcomment',
    post_id_desc bigint GENERATED ALWAYS AS (post_id * -1) VIRTUAL,
    created      datetime NOT NULL,
    updated      datetime,
    CONSTRAINT comments_id_pk PRIMARY KEY (id),
    CONSTRAINT comments_post_id_fk FOREIGN KEY (post_id) REFERENCES posts (id),
    CONSTRAINT comments_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE (post_id, user_id)
) COMMENT = 'Comments\nMulti-line\r\ntable\rcomment';
CREATE INDEX comments_post_id_user_id_idx ON comments (post_id, user_id) USING HASH;

CREATE TABLE comment_stars
(
    id              bigint AUTO_INCREMENT,
    user_id         int       NOT NULL,
    comment_post_id bigint    NOT NULL,
    comment_user_id int       NOT NULL,
    created         timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated         timestamp DEFAULT '1970-01-01 08:00:00',
    CONSTRAINT comment_stars_id_pk PRIMARY KEY (id),
    CONSTRAINT comment_stars_user_id_post_id_fk FOREIGN KEY (comment_post_id, comment_user_id) REFERENCES comments (post_id, user_id),
    CONSTRAINT comment_stars_user_id_fk FOREIGN KEY (comment_user_id) REFERENCES users (id),
    UNIQUE (user_id, comment_post_id, comment_user_id)
);

CREATE TABLE logs
(
    id              bigint PRIMARY KEY AUTO_INCREMENT,
    user_id         int      NOT NULL,
    post_id         bigint,
    comment_id      bigint,
    comment_star_id bigint,
    payload         text,
    created         datetime NOT NULL
) COMMENT = 'Auditログ';

CREATE VIEW post_comments AS
(
SELECT c.id, p.title, u2.username AS post_user, c.comment, u2.username AS comment_user, c.created, c.updated
FROM posts AS p
         LEFT JOIN comments AS c on p.id = c.post_id
         LEFT JOIN users AS u on u.id = p.user_id
         LEFT JOIN users AS u2 on u2.id = c.user_id
    );

CREATE TABLE CamelizeTable
(
    id      bigint PRIMARY KEY AUTO_INCREMENT,
    created datetime NOT NULL
);

CREATE TABLE `hyphen-table`
(
    id              bigint PRIMARY KEY AUTO_INCREMENT,
    `hyphen-column` text     NOT NULL,
    created         datetime NOT NULL
);

CREATE TRIGGER update_posts_updated
    BEFORE UPDATE
    ON posts
    FOR EACH ROW SET NEW.updated = CURRENT_TIMESTAMP();
