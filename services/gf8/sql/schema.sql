CREATE TABLE gf8dm7_scores (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    game_type INTEGER NOT NULL,
    musicid INTEGER NOT NULL,
    seq INTEGER NOT NULL,
    score INTEGER NOT NULL
);

CREATE TABLE gf8dm7_ranked_scores (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    game_type INTEGER NOT NULL,
    musicid INTEGER NOT NULL,
    seq INTEGER NOT NULL,
    score INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE gf8dm7_favorites (
    game_type INTEGER NOT NULL,
    musicid INTEGER NOT NULL,

    count INTEGER NOT NULL,

    PRIMARY KEY(game_type, musicid)
);

CREATE TABLE gf8dm7_shops (
    game_type INTEGER NOT NULL,
    pref INTEGER NOT NULL,
    name TEXT NOT NULL,

    points INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(game_type, pref, name)
);

CREATE TABLE gf8dm7_messages (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    game_type INTEGER NOT NULL,

    enabled INTEGER NOT NULL,
    message TEXT NOT NULL
);

CREATE TABLE gf8dm7_demomusic (
    game_type INTEGER NOT NULL,
    musicid INTEGER NOT NULL,

    PRIMARY KEY(game_type, musicid)
);
