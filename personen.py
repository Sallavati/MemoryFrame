import database

def db_request_to_map(person):
    return {"id": person[0], "name": person[1]}

def db_request_to_list(persons):
    return [db_request_to_map(p) for p in persons]

def get_by_id(id):
    res = database.query_db_with_atributes("SELECT * FROM persons where id = ?;", (id, ))
    return db_request_to_map(res.fetchall()[0])

def get_all():
    res = database.query_db("SELECT * FROM persons;")
    return db_request_to_list(res.fetchall())

def get_all():
    res = database.query_db("SELECT * FROM persons;")
    return db_request_to_list(res.fetchall())

def get_by_action_id(id):
    res = database.query_db_with_atributes("SELECT * FROM persons where id IN (SELECT person_id FROM persons_in_action where action_id = ?);", (id, ))
    return db_request_to_list(res.fetchall())

def get_by_action_setting_id(id):
    res = database.query_db_with_atributes("SELECT * FROM persons where id IN (SELECT person_id FROM persons_in_action_filter where action_setting_id = ?);", (id, ))
    return db_request_to_list(res.fetchall())