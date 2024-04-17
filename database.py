import sqlite3

def setupDB():
    dataBaseNew = False
    try:
        open("data.db")
    except FileNotFoundError:
        dataBaseNew = True

    dbConnection = sqlite3.connect("data.db")

    if dataBaseNew:
        print("newDB")
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
    return res

def migrate_old_db():
    