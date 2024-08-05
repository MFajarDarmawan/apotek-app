CREATE DATABASE apotek;

USE apotek;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role ENUM('admin', 'customer') NOT NULL
);

CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    quantity INT,
    description TEXT,
    price DECIMAL(10,2),
    image_url VARCHAR(255)
);


CREATE TABLE purchases (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255),
    product_id INT,
    quantity INT,
    purchased_at DATETIME,
    price DECIMAL(10,2),
    total_price DECIMAL(10,2),
    paid_amount DECIMAL(10,2),
    change_amount DECIMAL(10,2),
    FOREIGN KEY (product_id) REFERENCES products(id)
);


