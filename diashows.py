import database

def db_request_to_map(a):
    return {"id": a[0], 
            "title": a[1], 
            "folder": a[2], 
            "infotext": a[3]
            }

def db_request_to_list(diashows):
    return [db_request_to_map(d) for d in diashows]

def get_by_id(id):
    res = database.query_db_with_atributes("SELECT * FROM diashows where id = ?;", (id, ))
    return db_request_to_map(res.fetchall()[0])

def get_all():
    res = database.query_db("SELECT * FROM diashows;")
    return db_request_to_list(res.fetchall())


def get_by_diashow_filter_id(id):
    res = database.query_db_with_atributes("SELECT * FROM diashows where id IN (SELECT diashow_id FROM diashow_in_diashow_filter where diashow_settings_id = ?) ;", (id, ))
    return db_request_to_list(res.fetchall())