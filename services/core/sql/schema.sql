CREATE TABLE core_pcbid (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    pcbid TEXT NOT NULL UNIQUE,
    status INTEGER NOT NULL -- 0 = not enabled, 1 = enabled, 2 = blacklisted
);

CREATE TABLE core_assigned_card_numbers (
    number TEXT NOT NULL,
    label TEXT NOT NULL,
    PRIMARY KEY(number, label)
);
