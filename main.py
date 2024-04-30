import database
import personen
import aktionen
import tags
import diashows

if __name__ == "__main__":
    database.setupDB()

    print(diashows.get_by_diashow_filter_id(2))
    print(aktionen.get_by_id(3))