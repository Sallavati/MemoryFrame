import sqlite3
import json
import os

def setupDB():
    dataBaseNew = False
    try:
        open("data.db")
    except FileNotFoundError:
        dataBaseNew = True

    dbConnection = sqlite3.connect("data.db")

    if dataBaseNew:
        print("newDB")
        files = os.listdir("./")

        #import old db when existent
        if "aktionen.json" in files and "settings.json" in files:
            migrate_old_db()

        setupDB = open("setupDB.sql").read()
        dbConnection.cursor().executescript(setupDB)

def hardReset():
    dbConnection = sqlite3.connect("data.db")

    dropDB = open("dropDB.sql").read()
    dbConnection.cursor().executescript(dropDB)

    setupDB = open("setupDB.sql").read()
    dbConnection.cursor().executescript(setupDB)

def query_db(query):
    dbConnection = sqlite3.connect("data.db")
    res = dbConnection.cursor().execute(query)
    return res


def query_db_with_atributes(query, atributes):
    dbConnection = sqlite3.connect("data.db")
    res = dbConnection.cursor().execute(query, atributes)
    dbConnection.commit()
    return res

## Migrate old DB

def add_zero(n):
    if n < 10:
        return "0" + str(n)
    return str(n)

def time_obj_to_datetime(t):
    y = t["Year"]
    m = add_zero(t["Month"])
    d = add_zero(t["Day"])
    h = add_zero(t["Hour"])
    mi = add_zero(t["Minute"])
    return f"{y}-{m}-{d} {h}:{mi}"

def insert_old_actions():
    with open("./aktionen.json") as aktionen_file:
        aktionen = json.load(aktionen_file)

    for aktion in aktionen:
        personenIds = []
        for person in aktion["Visitors"]: #personen einfügen, falls noch nicht vorhanden
            res = query_db_with_atributes("SELECT * from persons WHERE name = ?;", (person["Name"],))
            if len(res.fetchall()) < 1:
                res = query_db_with_atributes("INSERT into persons (name) VALUES (?);", (person["Name"],))
                if res.rowcount != 1:
                    print("something went wrong inserting the person")
            
            res = query_db_with_atributes("SELECT * from persons WHERE name = ?;", (person["Name"],))  #id der person merken
            personenIds.append(res.fetchall()[0][0])

        tagIds = []
        for tag in aktion["Tags"]: #personen einfügen, falls noch nicht vorhanden
            res = query_db_with_atributes("SELECT * from tags WHERE name = ?;", (tag["Name"],))
            if len(res.fetchall()) < 1:
                res = query_db_with_atributes("INSERT into tags (name) VALUES (?);", (tag["Name"],))
                if res.rowcount != 1:
                    print("something went wrong inserting the tag")
            
            res = query_db_with_atributes("SELECT * from tags WHERE name = ?;", (tag["Name"],))  #id der person merken
            tagIds.append(res.fetchall()[0][0])

        res = query_db_with_atributes("INSERT into actions (title, folder, infotext, start) VALUES (?, ?, ?, ?);", (aktion["Title"],aktion["Folder"], aktion["InfoText"], time_obj_to_datetime(aktion["Start"])))
        if res.rowcount != 1:
            print("something went wrong inserting the aktion")
        else:
            res = query_db_with_atributes("SELECT * from actions WHERE title = ?;", (aktion["Title"], ))  #id der Aktion merken
            actionsID = res.fetchall()[0][0]

        for id in personenIds:
            res = query_db_with_atributes("INSERT into persons_in_action (person_id, action_id) VALUES (?, ?);", (id, actionsID))
            if res.rowcount != 1:
                print("something went wrong inserting the person")

        for id in tagIds:
            res = query_db_with_atributes("INSERT into tags_in_action (tag_id, action_id) VALUES (?, ?);", (id, actionsID))
            if res.rowcount != 1:
                print("something went wrong inserting the tag")

def insert_old_settings():
    with open("./settings.json") as settings_file:
        settings = json.load(settings_file)

    for d in settings["AdditionalPics"]:
        res = query_db_with_atributes("INSERT into diashows (title, folder, infotext) VALUES (?, ?, ?);", (d["Name"], d["Path"], d["InfoText"]))
        if res.rowcount != 1:
            print("something went wrong inserting the diashow")
    


def migrate_old_db():
    insert_old_actions()
    insert_old_settings()

        

    