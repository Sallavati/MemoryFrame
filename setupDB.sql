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
    person_id INTEGER,
    action_id INTEGER,

    FOREIGN KEY (person_id) REFERENCES persons(id),
    FOREIGN KEY (action_id) REFERENCES  actions(id)
);

CREATE TABLE tags_in_action (
    tag_id INTEGER,
    action_id INTEGER,

    FOREIGN KEY (tag_id) REFERENCES tags(id),
    FOREIGN KEY (action_id) REFERENCES  actions(id)
);

CREATE TABLE diashows (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    folder TEXT,
    infotext TEXT
);

CREATE TABLE action_settings(
    id INTEGER PRIMARY KEY  AUTOINCREMENT,
    title TEXT,
    editing_allowed bool,
    upload_from_gallery_allowed bool,
    show_actions bool,
    filter_with_tags bool,
    filter_with_persons bool,
    filter_with_time bool,
    time_start datetime,
    time_end datetime
);

CREATE TABLE tags_in_action_filter (
    tag_id INTEGER,
    action_settings_id INTEGER,

    FOREIGN KEY (tag_id) REFERENCES tags(id),
    FOREIGN KEY (action_settings_id) REFERENCES  action_settings(id)
);

CREATE TABLE persons_in_action_filter (
    person_id INTEGER,
    action_settings_id INTEGER,

    FOREIGN KEY (person_id) REFERENCES persons(id),
    FOREIGN KEY (action_settings_id) REFERENCES  action_settings(id)
);

CREATE TABLE diashow_settings(
    id INTEGER PRIMARY KEY  AUTOINCREMENT,
    title TEXT,
    editing_allowed bool
);

CREATE TABLE diashow_in_diashow_filter (
    diashow_id INTEGER,
    diashow_settings_id INTEGER,

    FOREIGN KEY (diashow_id) REFERENCES diashows(id),
    FOREIGN KEY (diashow_settings_id) REFERENCES  diahow_settings(id)
);