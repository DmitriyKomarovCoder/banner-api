CREATE TABLE IF NOT EXISTS features (
    feature_id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS banners (
    banner_id SERIAL PRIMARY KEY,
    content JSONB NOT NULL,
    active BOOLEAN NOT NULL,
    feature_id INT REFERENCES features(feature_id),
    created_at timestamp DEFAULT now(),
    update_at timestamp DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tags (
    tag_id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS banner_tags (
    banner_id INT REFERENCES banners(banner_id),
    tag_id INT REFERENCES tags(tag_id),
    PRIMARY KEY (banner_id, tag_id)
);

-- CREATE TABLE deletion_requests (
--     request_id SERIAL PRIMARY KEY,
--     feature_id INT,
--     tag_id INT
-- );

-- =====================================

INSERT INTO tags (name) VALUES ('Tag1'), ('Tag2'), ('Tag3'), ('Tag4');

INSERT INTO features (name) VALUES ('Feature1'), ('Feature2'), ('Feature3'), ('Feature4');
