CREATE TABLE gf10dm9_card_profile (
    game_type INTEGER NOT NULL,
    cardid TEXT NOT NULL UNIQUE,

    name TEXT NOT NULL,
    pass VARCHAR(4) NOT NULL,
    type INTEGER NOT NULL DEFAULT 0,
    update_flag INTEGER NOT NULL DEFAULT 0, -- TODO: What is this?
    recovery INTEGER NOT NULL DEFAULT 0,
    -- game_type text,
    skill INTEGER NOT NULL DEFAULT 0,
    expired INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY(game_type, cardid)
);

CREATE TABLE gf10dm9_puzzle (
    game_type INTEGER NOT NULL,
    cardid TEXT NOT NULL,

    puzzle_no INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    hidden INTEGER NOT NULL,

    PRIMARY KEY(game_type, cardid)
);

CREATE TABLE gf10dm9_scores (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    cardid TEXT NOT NULL,
    game_type INTEGER NOT NULL,

    netid INTEGER NOT NULL DEFAULT -1,
    courseid INTEGER NOT NULL DEFAULT -1,
    seq_mode INTEGER NOT NULL,
    flags INTEGER NOT NULL,
    score INTEGER NOT NULL,
    clear INTEGER NOT NULL,
    combo INTEGER NOT NULL,
    skill INTEGER NOT NULL,
    perc INTEGER NOT NULL,
    irall INTEGER NOT NULL DEFAULT 0, -- TODO: Remove these defaults before release
    ircom INTEGER NOT NULL DEFAULT 0
);

CREATE TABLE gf10dm9_shops (
    game_type INTEGER NOT NULL,
    sid VARCHAR(12) NOT NULL,

    name VARCHAR(16) NOT NULL,
    pref INTEGER NOT NULL,
    points INTEGER NOT NULL DEFAULT 0,

    PRIMARY KEY (game_type, sid)
);

