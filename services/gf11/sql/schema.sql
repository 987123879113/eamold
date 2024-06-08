CREATE TABLE gf11dm10_card_profile (
    game_type INTEGER NOT NULL DEFAULT 0,
    gdid INTEGER NOT NULL,

    cardid TEXT NOT NULL UNIQUE,
    irid TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    pass VARCHAR(4) NOT NULL,
    type INTEGER NOT NULL DEFAULT 0,
    update_flag INTEGER NOT NULL DEFAULT 0, -- TODO: What is this?
    puzzle_no INTEGER NOT NULL DEFAULT 0,
    recovery INTEGER NOT NULL DEFAULT 0,
    skill INTEGER NOT NULL DEFAULT 0,
    expired INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(game_type, gdid)
);

CREATE TABLE gf11dm10_puzzle (
    game_type INTEGER NOT NULL DEFAULT 0,
    gdid INTEGER NOT NULL,

    puzzle_no INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    hidden INTEGER NOT NULL,

    PRIMARY KEY(game_type, gdid, puzzle_no)
);

CREATE TABLE gf11dm10_lands (
    game_type INTEGER NOT NULL DEFAULT 0,
    team INTEGER NOT NULL PRIMARY KEY,
    round INTEGER NOT NULL,
    area INTEGER NOT NULL,
    hidden INTEGER NOT NULL,
    point INTEGER NOT NULL,
    spoint INTEGER NOT NULL
);

CREATE TABLE gf11dm10_scores (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    gdid INTEGER NOT NULL,
    game_type INTEGER NOT NULL DEFAULT 0,

    netid INTEGER NOT NULL DEFAULT -1,
    courseid INTEGER NOT NULL DEFAULT -1,
    seq_mode INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    score INTEGER NOT NULL,
    clear INTEGER NOT NULL,
    combo INTEGER NOT NULL,
    skill INTEGER NOT NULL,
    perc INTEGER NOT NULL,
    irall INTEGER NOT NULL DEFAULT 0,
    ircom INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE gf11dm10_shops (
    game_type INTEGER NOT NULL DEFAULT 0,
    sid VARCHAR(12) NOT NULL,

    name VARCHAR(16) NOT NULL,
    pref INTEGER NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (game_type, sid)
);
