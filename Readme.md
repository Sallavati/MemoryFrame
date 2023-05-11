# MemoryFrame

## Grundkonzept
Bei MemoryFrame handelt es sich um ein Web basiertes Bilderrahmen-System, um bei Feiern, einem entspannten Abend oder jeglicher anderen Form von Aktion nichts zu verpassen und eine coole Diashow zu bieten. Dazu ist das System in 3 Hauptbestandteile aufgeteilt:

### Aktionen
Das Herzsstück des Systems sind die Aktionen, dieses sind jeweils Zeitphasen, in welchen Teilnehmer live Fotos, über die Upload-Seite hochladen können, welche dann im dazugehörigen Ordner gespeichert werden. Zusätzlich kann man der Aktion auch Teilnehmer, Tags, einen Infotext und einen Titel zuweisen. Gestartet werden Aktionen durch den ersten Foto upload, gestoppt entweder Manuell über die Settingsseite oder automatisch, wenn eine bestimmte Zeit lang (bei Settings kofigurierbar) kein Bild mehr hochgeladen wurde. Im nachhinein, kann man bei "browse Actions" alle Aktionen betrachten und bearbeiten.

### Externe Contents
Damit man nicht nur Fotos betrachten kann, welche vor Ort entstanden sind, gibt es die Möglichkeit externe Fotos über Extern Content - Ordner einzubinden. Diese kann man in den Settings einrichten, verwalten und auch einen extra Link für Foto Management eines bestimmten Ordners Freischalten. Zusatz Informationen die man zu solchen Ordnern angeben kann/muss sind ein Infotext und ein Titel.

### Diashow
Was wäre ein Bilderramen ohne Bilderanzeige. Für die Diashow gibt es zwei mögliche Ansichten, entweder mit Informationen zur jeweiligen Aktion/Externen Ordner aus dem das aktuelle Bild stammt und eine ohne. Zudem kann man in den Einstellungen selbstverständlich anpassen was man sehen möchte. So kann man zum Beispiel einstellen, welche Externen Bilder eingebunden werden sollen, ob vielleicht sogar nur externe Bilder angezeigt werden sollen. Oder Aktionen nach verschiedenen Parametern sortieren.

## Inbetriebnahme
### Vorraussetzungen:
* golang instaliert
  
### Schritte
* Standart Datenbank aus dem Ordner defaultData ins Hauptverzeichnis kopieren
* .env Datei anlegen und Folgende Variablen angeben:
  * IpAddress --> Adresse, über welche die Website erreichbar sein soll.
  * Port --> Port auf welchen die Website geöffnet werden soll.
* starten: go run main.go