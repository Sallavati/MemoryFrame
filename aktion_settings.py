import database
import tags
import personen

def db_request_to_map(a):
    return {"id": a[0], 
            "title": a[1], 
            "upload_from_gallery_allowed": a[2], 
            "show_actions": a[3],
            "filter_with_tags": a[4],
            "filter_with_persons": a[5],
            "filter_with_time": a[6],
            "time_start": a[7],
            "time_end": a[8],
            "tags": tags.get_by_action_setting_id(a[0]),
            "persons": personen.get_by_action_setting_id(a[0])
            }

def db_request_to_list(diashow_filter):
    return [db_request_to_map(df) for df in diashow_filter]

def get_by_id(id):
    res = database.query_db_with_atributes("SELECT * FROM diashow_settings where id = ?;", (id, ))
    return db_request_to_map(res.fetchall()[0])

def get_all():
    res = database.query_db("SELECT * FROM diashow_settings;")
    return db_request_to_list(res.fetchall())