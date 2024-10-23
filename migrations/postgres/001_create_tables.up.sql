CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(100),
    user_name VARCHAR(40),
    password VARCHAR(100),
    bio VARCHAR(255),
    profile_picture TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at INT DEFAULT 0
);

CREATE TABLE tweets (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    content TEXT,
    media TEXT,
    views_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE views (
    user_id UUID REFERENCES users(id),
    tweet_id UUID REFERENCES tweets(id),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, tweet_id)
);

CREATE TABLE likes (
    user_id UUID REFERENCES users(id),
    tweet_id UUID REFERENCES tweets(id),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, tweet_id)
);

CREATE TABLE retweets (
    user_id UUID REFERENCES users(id),
    tweet_id UUID REFERENCES tweets(id),
    PRIMARY KEY (user_id, tweet_id)
);

CREATE TABLE follows (
    follower_id UUID REFERENCES users(id),
    following_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id)
);

