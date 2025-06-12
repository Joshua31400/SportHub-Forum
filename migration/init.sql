CREATE TABLE IF NOT EXISTS user (
                                    userID INT AUTO_INCREMENT PRIMARY KEY,
                                    userName VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NULL, -- Optionnel pour OAuth
    google_id VARCHAR(255) UNIQUE, -- Ajouter dès la création
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
    FOREIGN KEY (userid) REFERENCES user(userID) -- Correction: userID au lieu de userid
    );

CREATE INDEX IF NOT EXISTS idx_user_google_id ON user(google_id);
CREATE INDEX IF NOT EXISTS idx_user_auth_provider ON user(auth_provider);
CREATE INDEX IF NOT EXISTS idx_user_email ON user(email);

INSERT IGNORE INTO user (userName, email, password, auth_provider, is_verified)
VALUES ('admin', 'admin@example.com', '$2a$10$hash...', 'local', TRUE);

ALTER TABLE user ADD COLUMN github_id VARCHAR(255) UNIQUE;
CREATE INDEX idx_user_github_id ON user(github_id);