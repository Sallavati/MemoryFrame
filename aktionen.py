import database
import personen
import tags

def db_request_to_map(a):
    return {"id": a[0], 
            "title": a[1], 
            "folder": a[2], 
            "infotext": a[3], 
            "start": a[4], 
            "persons": personen.get_by_action_id(a[0]),
            "tags": tags.get_by_action_id(a[0])
            }

def db_request_to_list(actions):
    return [db_request_to_map(a) for a in actions]

def get_by_id(id):
    res = database.query_db_with_atributes("SELECT * FROM actions where id = ?;", (id, ))
    return db_request_to_map(res.fetchall()[0])

def get_all():
    res = database.query_db("SELECT * FROM actions;")
    return db_request_to_list(res.fetchall())