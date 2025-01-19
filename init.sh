#!/bin/bash
set -e

mysql -u root -proot <<-EOSQL

CREATE DATABASE IF NOT EXISTS user_tracker;
USE user_tracker;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

EOSQL
