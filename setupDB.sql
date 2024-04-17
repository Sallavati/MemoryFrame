CREATE TABLE persons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT
);

CREATE TABLE actions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    folder TEXT,
    infotext TEXT,
    start datetime
);

CREATE TABLE persons_in_action (
    person_id INTEGER FOREIGN KEY referencespersons(id),
    action_id INTEGER FOREIGN KEY references actions(id)
);

CREATE TABLE tags_in_action (
    tag_id INTEGER FOREIGN KEY references tags(id),
    action_id INTEGER FOREIGN KEY references actions(id)
);

CREATE TABLE diashows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    folder TEXT,
    infotext TEXT
);

CREATE TABLE tags_in_diahow (
    tag_id INTEGER FOREIGN KEY references tags(id),
    diashow_id INTEGER FOREIGN KEY references diashows(id)
);

CREATE TABLE action_settings(
    id INTEGER PRIMARY KEY,
    title TEXT,
    editing_allowed bool,
    show_actions bool,
    filter_with_tags bool,
    filter_with_persons bool,
    filter_with_time bool,
    time_start datetime,
    time_end datetime
);

CREATE TABLE tags_in_action_filter (
    tag_id INTEGER FOREIGN KEY references tags(id),
    action_settings_id INTEGER FOREIGN KEY references action_settings(id)
);

CREATE TABLE persons_in_action_filter (
    person_id INTEGER FOREIGN KEY references persons(id),
    action_settings_id INTEGER FOREIGN KEY references action_settings(id)
);

CREATE TABLE diahow_settings(
    id INTEGER PRIMARY KEY,
    title TEXT,
    editing_allowed bool,
    show_diahows bool,
    filter_with_tags bool
);

CREATE TABLE tags_in_diahow_filter (
    tag_id INTEGER FOREIGN KEY references tags(id),
    diashow_settings_id INTEGER FOREIGN KEY references diahow_settings(id)
);