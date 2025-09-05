CREATE TABLE
    `users` (
        id INT NOT NULL AUTO_INCREMENT,
        -- 
        email VARCHAR(255) UNIQUE NOT NULL,
        passwordHash VARCHAR(255) NOT NULL,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        -- 
        PRIMARY KEY (id)
    );

CREATE TABLE
    `events` (
        id INT NOT NULL AUTO_INCREMENT,
        -- 
        ownerUserId INT NOT NULL,
        -- 
        label VARCHAR(800),
        -- 
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        -- 
        PRIMARY KEY (id),
        FOREIGN KEY (ownerUserId) REFERENCES users (id)
    );