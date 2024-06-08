CREATE TABLE gf8dm7puv_card_profile (
    game_type INTEGER NOT NULL,
    cardid TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL,
    color INTEGER NOT NULL,
    recovery INTEGER NOT NULL,
    styles INTEGER NOT NULL,
    hidden INTEGER NOT NULL,
    expired INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(game_type, cardid)
);

CREATE TABLE gf8dm7puv_puzzles (
    game_type INTEGER NOT NULL,
    cardid TEXT NOT NULL UNIQUE,

    number INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    out INTEGER NOT NULL,

    PRIMARY KEY(game_type, cardid, number)
);

CREATE TABLE gf8dm7puv_favorites (
    game_type INTEGER NOT NULL,
    musicid INTEGER NOT NULL,

    count INTEGER NOT NULL,

    PRIMARY KEY(game_type, musicid)
);

CREATE TABLE gf8dm7puv_shops (
    game_type INTEGER NOT NULL,
    pref INTEGER NOT NULL,
    name TEXT NOT NULL,

    points INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(game_type, pref, name)
);

CREATE TABLE gf8dm7puv_shop_machines (
    game_type INTEGER NOT NULL,
    pcbid TEXT NOT NULL,

    pref INTEGER NOT NULL,
    name TEXT NOT NULL,

    PRIMARY KEY(game_type, pcbid)
);

CREATE TABLE gf8dm7puv_scores (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    cardid TEXT NOT NULL,
    game_type INTEGER NOT NULL,
    musicid INTEGER NOT NULL,
    musicnum INTEGER NOT NULL,
    seq INTEGER NOT NULL,
    score INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    clear INTEGER NOT NULL,
    skill INTEGER NOT NULL,
    combo INTEGER NOT NULL,
    encore INTEGER NOT NULL DEFAULT 0,
    extra INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE gf8dm7puv_messages (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    game_type INTEGER NOT NULL,

    enabled INTEGER NOT NULL,
    message TEXT NOT NULL
);

CREATE TABLE gf8dm7puv_demomusic (
    game_type INTEGER NOT NULL,
    musicid INTEGER NOT NULL,

    PRIMARY KEY(game_type, musicid)
);
