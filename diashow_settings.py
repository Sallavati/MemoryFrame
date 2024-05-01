import database
import diashows

def db_request_to_map(df):
    dias = diashows.get_all()
    dias_to_show = [d["id"] for d in diashows.get_by_diashow_filter_id(df[0])]
    for d in dias:
        d["show"] = d["id"] in dias_to_show
        
    return {"id": df[0], 
            "title": df[1], 
            "editing_allowed": df[2], 
            "diashows": dias
            }

def db_request_to_list(diashow_filter):
    return [db_request_to_map(df) for df in diashow_filter]

def get_by_id(id):
    res = database.query_db_with_atributes("SELECT * FROM diashow_settings where id = ?;", (id, ))
    return db_request_to_map(res.fetchall()[0])

def get_all():
    res = database.query_db("SELECT * FROM diashow_settings;")
    return db_request_to_list(res.fetchall())