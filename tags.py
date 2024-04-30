import database

def db_request_to_map(tag):
    return {"id": tag[0], "name": tag[1]}

def db_request_to_list(tags):
    return [db_request_to_map(t) for t in tags]

def get_by_id(id):
    res = database.query_db_with_atributes("SELECT * FROM tags where id = ?;", (id, ))
    return db_request_to_map(res.fetchall()[0])

def get_all():
    res = database.query_db("SELECT * FROM tags;")
    return db_request_to_list(res.fetchall())

def get_all():
    res = database.query_db("SELECT * FROM tags;")
    return db_request_to_list(res.fetchall())

def get_by_action_id(id):
    res = database.query_db_with_atributes("SELECT * FROM tags where id IN (SELECT tag_id FROM tags_in_action where action_id = ?);", (id, ))
    return db_request_to_list(res.fetchall())

def get_by_action_setting_id(id):
    res = database.query_db_with_atributes("SELECT * FROM tags where id IN (SELECT tag_id FROM tags_in_action_filter where action_setting_id = ?);", (id, ))
    return db_request_to_list(res.fetchall())