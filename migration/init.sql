-- Table des utilisateurs
CREATE TABLE IF NOT EXISTS user (
    userID INT AUTO_INCREMENT PRIMARY KEY,
    userName VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    createdAt DATETIME NOT NULL
);

-- Table des sessions
CREATE TABLE IF NOT EXISTS sessions (
    id INT PRIMARY KEY AUTO_INCREMENT,
    userID INT NOT NULL,
    sessionToken VARCHAR(255) NOT NULL,
    expiresAt DATETIME NOT NULL,
    FOREIGN KEY (userID) REFERENCES users(userID)
);

-- Table des catégories
CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Insertion des catégories initiales
INSERT INTO categories (name) VALUES
    ('nutrition'),
    ('entrainement'),
    ('equipement'),
    ('motivation'),
    ('football')
ON DUPLICATE KEY UPDATE name=name;

-- Table des posts
CREATE TABLE IF NOT EXISTS post (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    userID INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userID) REFERENCES users(userID)
);

-- Table de relation many-to-many entre post et catégories
CREATE TABLE IF NOT EXISTS post_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    postID INT NOT NULL,
    categoryID INT NOT NULL,
    FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (categoryID) REFERENCES categories(id),
    UNIQUE KEY unique_post_category (postID, categoryID)
);

-- Table des commentaires
CREATE TABLE IF NOT EXISTS comment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    postID INT NOT NULL,
    userID INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (userID) REFERENCES users(userID)
);

-- Table des likes
CREATE TABLE IF NOT EXISTS like_post (
    id INT AUTO_INCREMENT PRIMARY KEY,
    postID INT NOT NULL,
    userID INT NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (postID) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (userID) REFERENCES users(userID),
    UNIQUE KEY unique_like (userID, postID)
)