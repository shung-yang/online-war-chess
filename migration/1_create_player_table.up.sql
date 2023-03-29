CREATE TABLE IF NOT EXISTS player (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(100) UNIQUE NOT NULL,
  email varchar(100) UNIQUE DEFAULT NULL,
  password varbinary(100) DEFAULT NULL,
  level tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (id)
);