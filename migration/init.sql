CREATE TABLE IF NOT EXISTS user (
    userID INT AUTO_INCREMENT PRIMARY KEY,
    userName VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    createdAt DATETIME NOT NULL
    google_id VARCHAR(255) UNIQUE,
    github_id VARCHAR(255) UNIQUE,
    avatar TEXT,
    auth_provider VARCHAR(50) DEFAULT 'local',
    is_verified BOOLEAN DEFAULT FALSE,
    createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS session (
    id INT PRIMARY KEY AUTO_INCREMENT,
    userid INT NOT NULL,
    sessiontoken VARCHAR(255) NOT NULL,
    expiresat VARCHAR(255) NOT NULL,
    FOREIGN KEY (userid) REFERENCES user(userid)
    );

CREATE INDEX IF NOT EXISTS idx_user_google_id ON user(google_id);
CREATE INDEX IF NOT EXISTS idx_user_auth_provider ON user(auth_provider);
CREATE INDEX IF NOT EXISTS idx_user_email ON user(email);
CREATE INDEX idx_user_github_id ON user(github_id);

INSERT IGNORE INTO user (userName, email, password, auth_provider, is_verified)
VALUES ('admin', 'admin@example.com', '$2a$10$hash...', 'local', TRUE);


CREATE TABLE IF NOT EXISTS category (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
    );

INSERT INTO category (name) VALUES ('nutrition') ON DUPLICATE KEY UPDATE name=name;
INSERT INTO category (name) VALUES ('entrainement') ON DUPLICATE KEY UPDATE name=name;
INSERT INTO category (name) VALUES ('equipement') ON DUPLICATE KEY UPDATE name=name;
INSERT INTO category (name) VALUES ('motivation') ON DUPLICATE KEY UPDATE name=name;
INSERT INTO category (name) VALUES ('football') ON DUPLICATE KEY UPDATE name=name;

CREATE TABLE IF NOT EXISTS post (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    categoryid INT NOT NULL,
    userid INT NOT NULL,
    createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(100),
    FOREIGN KEY (categoryid) REFERENCES category(id),
    FOREIGN KEY (userid) REFERENCES user(userid)
    );

CREATE TABLE IF NOT EXISTS comment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    postid INT NOT NULL,
    userid INT NOT NULL,
    createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (postid) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES user(userid)
    );

CREATE TABLE IF NOT EXISTS `like` (
    id INT AUTO_INCREMENT PRIMARY KEY,
    postid INT NOT NULL,
    userid INT NOT NULL,
    createdat TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (postid) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (userid) REFERENCES user(userid),
    UNIQUE KEY unique_like (userid, postid)
    );