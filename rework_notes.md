## Settings Json

- Colorsheme
- Passwort Schutz aktiv
- Passwort
- Fotobox Aktiv
- Fotobox Kamera Auswahl
- Foto komprimierung
- Current action settings
- Diashow geschwindigkeit
- Aktions auto stop
- aktuelle action / diashow settings ID

## Datenbank

- Diashow Settings
    - Id
    - Alias
    - Anzeigen
    - Editing allowed
    - Tag Filter active
- Diashow Tag Filter
    - Diashow Settings Id
- Aktionen Settings
    - Id
    - Alias
    - Editing allowed
    - Anzeigen
    - Tag Filter active
    - Person Filter active
    - Time Filter active
    - start/stop time
- Aktionen Tag Filter
    - Tag Id
    - Aktion Settings ID
- Aktionen Personen Filter
    - Person ID
    - Aktionen Settings Id
- Tags
    - Id
    - Name
- Personen
    - Id
    - Name
- Aktion
    - Id
    - Name
    - Start
    - Infotext
    - Folder
- Tags in Aktion
    - Tag ID
    - Aktion ID
- Personen in Aktion
    - Person Id
    - Aktion ID
- Diashow
    - Id
    - Folder
    - Infotext
- Tags in Diashow
    - Tag ID
    - Diashow Id
- Login Token
    - Token
    - Created time

## Websites

- Home
    - Über Text
    - Links
- Settings
    - Diashow Settings
        - Presets auswählen
        - Presets hinzufügen / Verwalten
            - Einzelne seite
    - Aktions Settings
        - Wie Diashow Settings
    - Personen verwalten
        - Einzelne Seite
    - Tags verwalten
        - Einzelne seite
    - Design settings
        - Farben
    - Pin Schutz
        - Aktivieren / Deaktivieren
        - Pin ändern
    - Fotobox Settings
- Diashows verwalten
    - Diashows auflisten für Einstellungen einzelne Seite
    - Tags verwalten
    - Bilder einzeln hinzufügen/ Löschen
    - Zip export / import
    - Bilder Bearbeitung mit Auswahl Funktion
- Aktions Verwaltung wie Diashow Verwaltung
- Fotobox seite
    - Foto machen —> Countdown
    - Foto anzeigen, Frage Upload/verwerfen
    - 
- Aktuelle Aktion
    - Aktionsverwaltung (Titel, Persone, Tags, infotext)
    - Fotoupload
    - Foto verwaltung

## Backend File Struckture

- migrate_old_db.py
  - migrates old data base model into new Database

- main.py 
  - initialisation (Databank)
  - finding variables like Ip address
  - starting server

-  database.py
   - initialisation script for db
   - trackes db changes for effiency
   - db access abstruction function

- person/tags.py
  - create tags / persons
  - changes tags / persons
  - tag / person api
  - tag / person access in general

- diashow / action.py
  - diashow Action changes
  - settings changes

- pictures.py
  - diashow api
  - picture upload endpoint
  - pictures to zip
  - zip upload
  - fotobox script
  - picture comprimisation

- webserver.py
  - init webserver
  - has webserver object
  - hosts frontend

- settings.py
  - manages settings.json