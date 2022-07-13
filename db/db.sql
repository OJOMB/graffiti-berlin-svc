CREATE TABLE districts (
    id int AUTO_INCREMENT,
    name varchar(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE users (
    id varchar(36),
    user_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(40)  NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT uc_username UNIQUE (username),
    CONSTRAINT uc_email UNIQUE (email)
);

CREATE INDEX idx_user_email ON users (email);
CREATE INDEX idx_user_username ON users (username);

CREATE TABLE artists (
    id varchar(36),
    name varchar(100) NOT NULL,
    instagram varchar(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE aliases (
    artist varchar(36) NOT NULL,
    alias varchar(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (artist, alias),
    CONSTRAINT fk_aliases_artist FOREIGN KEY (artist) REFERENCES artists(id),
    CONSTRAINT fk_aliases_alias FOREIGN KEY (alias) REFERENCES artists(id)
);

CREATE TABLE crews (
    id varchar(36),
    name varchar(100) NOT NULL,
    acronym varchar(20),
    instagram varchar(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE affiliations (
    artist varchar(36) NOT NULL,
    crew varchar(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (artist, crew),
    CONSTRAINT fk_affiliations_artist FOREIGN KEY (artist) REFERENCES artists(id),
    CONSTRAINT fk_affiliations_crew FOREIGN KEY (crew) REFERENCES countries(id)
);

CREATE INDEX idx_affiliations_crew ON affiliations (crew);

CREATE TABLE pieces (
    id varchar(36),
    img varchar(255) NOT NULL,
    type int NOT NULL,
    uploaded_by varchar(36) NOT NULL,
    district int,
    geo_location point,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_pieces_district FOREIGN KEY (district) REFERENCES districts(id),
    CONSTRAINT fk_pieces_uploadedBy FOREIGN KEY (uploadedBy) REFERENCES users(id),
    CONSTRAINT fk_pieces_type FOREIGN KEY (type) REFERENCES piece_types(id)
);

CREATE TABLE duplicates (
    original varchar(36),
    duplicate varchar(36),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (original, duplicate)
);

CREATE TABLE piece_artists (
    piece int AUTO_INCREMENT,
    artist varchar(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (piece, artist)
);

CREATE TABLE piece_crews (
    piece int AUTO_INCREMENT,
    crew varchar(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (piece, crew)
);

CREATE TABLE piece_types (
    id int AUTO_INCREMENT,
    name varchar(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE piece_tags (
    piece varchar(36),
    tag varchar(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (piece, tag)
);
