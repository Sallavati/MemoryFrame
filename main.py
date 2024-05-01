import database
import aktionen
import diashow_settings
import webserver

if __name__ == "__main__":
    database.setupDB()

    print(diashow_settings.get_all())
    print(aktionen.get_by_id(3))

    webserver.app.run(debug=True, host='0.0.0.0')