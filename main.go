package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	godotenv "github.com/joho/godotenv"
	qrcode "github.com/skip2/go-qrcode"
)

type Person struct {
	Name         string
	TimesVisited int
}

type Tag struct {
	Name      string
	TimesUsed int
}

type onlyName struct {
	Name string
}

type onlyPath struct {
	Path string
}

type TimeforJson struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

type Aktion struct {
	Title    string
	Visitors []onlyName
	Tags     []onlyName
	Start    TimeforJson
	End      TimeforJson
	Folder   string
	InfoText string
}

type additionalFolder struct {
	Name       string
	InfoText   string
	Path       string
	Show       bool
	FreeUpload bool
}

type DiashowfilterStruct struct {
	OnlyExternContent  bool
	SortAfterTimeStart bool
	SortAfterTimeEnd   bool
	TimeStart          TimeforJson
	TimeEnd            TimeforJson
	SortAfterTags      bool
	Tags               []onlyName
	SortAfterPersons   bool
	Persons            []onlyName
}

type Settings struct {
	DiashowSpeed         int
	AutoActionStop       int
	AllowPersonEditing   bool
	AllowTagEditing      bool
	AllowInfoTextEditing bool
	AdditionalPics       []additionalFolder
	Diashowfilter        DiashowfilterStruct
}

type Diashow struct {
	FromAktion   bool
	DiashowSpeed int
	Picture      string
	Aktion       Aktion
	Folder       additionalFolder
}

type Statistic struct {
	Visitors           []Person
	Tags               []Tag
	ActionCount        int
	ExternPictureCount int
	ActionPictureCount int
	VisitorCount       int
	TagCount           int
	LastAktion         Aktion
	DataSize           int64
}

var aktionLive bool = false

var currentActionID int = 2

var LastTimeUploaded time.Time

func main() {
	http.HandleFunc("/api/upload", uploadFile)
	http.HandleFunc("/api/uploadToAdditionalContent", uploadToAdditionalContent)
	http.HandleFunc("/api/rmTag", rmTag)
	http.HandleFunc("/api/addTag", addTag)
	http.HandleFunc("/api/rmPerson", rmPerson)
	http.HandleFunc("/api/addPerson", addPerson)
	http.HandleFunc("/api/editAktionInfotext", editAktionInfotext)
	http.HandleFunc("/api/changeAktionName", changeAktionName)
	http.HandleFunc("/api/currentAction", returnCurrentAction)
	http.HandleFunc("/api/tags", returnTags)
	http.HandleFunc("/api/settings", returnSettings)
	http.HandleFunc("/api/actions", returnActions)
	http.HandleFunc("/api/actionPictures", returnPictures)
	http.HandleFunc("/api/persons", returnPersons)
	http.HandleFunc("/api/diashow", returnDiashow)
	http.HandleFunc("/api/stats", returnStats)
	http.HandleFunc("/api/stopAction", stopAktionWeb)
	http.HandleFunc("/api/getExternContentPics", getExternPics)
	http.HandleFunc("/api/toggleSortAfterTags", toggleSortAfterTags)
	http.HandleFunc("/api/toggleSortAfterPersons", toggleSortAfterPersons)
	http.HandleFunc("/api/toggleSortAfterTimeEnd", toggleSortAfterTimeEnd)
	http.HandleFunc("/api/toggleSortAfterTimeStart", toggleSortAfterTimeStart)
	http.HandleFunc("/api/setFilterTimeStart", setFilterTimeStart)
	http.HandleFunc("/api/setFilterTimeEnd", setFilterTimeEnd)
	http.HandleFunc("/api/toggleOnlyAdditionalContent", toggleOnlyAdditionalContent)
	http.HandleFunc("/api/editExternContentInfoText", editExternContentInfoText)
	http.HandleFunc("/api/addAdditionalContent", addAdditionalContent)
	http.HandleFunc("/api/changeAction", changeAktion)
	http.HandleFunc("/api/deletePictureFromAction", deletePicFromAction)
	http.HandleFunc("/api/deleteAdditionalContent", deleteAdditionalContent)
	http.HandleFunc("/api/deleteAktion", deleteAction)
	http.HandleFunc("/api/deletePictureFromAdditionalContent", deletePicFromAdditionalContent)
	http.HandleFunc("/api/toggleShowExternContent", toggleShowExternContent)
	http.HandleFunc("/api/toggleFreeUpload", toggleFreeUpload)
	http.HandleFunc("/api/toggleAllowTagEditing", toggleAllowTagEditing)
	http.HandleFunc("/api/toggleAllowInfoTextEditing", toggleAllowInfoTextEditing)
	http.HandleFunc("/api/toggleAllowPersonEditing", toggleAllowPersonEditing)
	http.HandleFunc("/api/ignoreTag", ignoreTag)
	http.HandleFunc("/api/ignorePerson", ignorePerson)
	http.HandleFunc("/api/setTag", setTag)
	http.HandleFunc("/api/setDiashowSpeed", setDiashowSpeed)
	http.HandleFunc("/api/setAutoActionStop", setAutoActionStop)
	http.HandleFunc("/api/setPerson", setPerson)

	fs := http.FileServer(http.Dir("./externContent"))
	http.Handle("/externContent/", http.StripPrefix("/externContent", fs))

	vueRoutes := []string{"diashow", "settings", "upload", "diashow2", "browse", "uploadToExtern/"}
	for _, route := range vueRoutes {
		http.HandleFunc("/"+route, returnIndexHTML)
	}

	hostFrontend()
	hostStatic()
	hostActions()

	godotenv.Load(".env")
	var IpAdress string = os.Getenv("IpAddress")
	var Port string = os.Getenv("Port")

	link := IpAdress
	if Port != ":80" {
		link += Port
	}

	link += "/upload"
	fmt.Println(link)

	qrcode.WriteColorFile(link, qrcode.Medium, 256, color.White, color.Black, "./static/qr.png")

	rand.Seed(int64(time.Now().Nanosecond()))

	// Start the server.
	fmt.Println("Server listening on port " + Port)
	log.Panic(
		http.ListenAndServe(Port, nil),
	)
}

func addTag(w http.ResponseWriter, r *http.Request) {
	if aktionLive {
		newTag := Tag{
			Name:      r.FormValue("tag"),
			TimesUsed: 1,
		}

		if newTag.Name != "" {
			fmt.Printf("New Tag: %+v", newTag)
			tagAdded := false //Ist der Tag bereits hinzugefügt worden?
			aktionenJson, _ := ioutil.ReadFile("aktionen.json")
			var aktionen []Aktion
			json.Unmarshal([]byte(aktionenJson), &aktionen)
			currentAction := aktionen[currentActionID]
			for _, tag := range currentAction.Tags {
				if newTag.Name == tag.Name {
					tagAdded = true
					break
				}
			}

			tagExists := false //Ist der Tag bereits verwendet worden?
			var tagID int
			var timesUsed int = 0
			tagsJson, _ := ioutil.ReadFile("tags.json")
			var tags []Tag
			json.Unmarshal([]byte(tagsJson), &tags)
			for index, tag := range tags {
				if newTag.Name == tag.Name {
					tagExists = true
					timesUsed = tag.TimesUsed
					tagID = index
					break
				}
			}

			if !tagAdded {
				aktionen[currentActionID].Tags = append(aktionen[currentActionID].Tags, onlyName{Name: newTag.Name})
				timesUsed++

				newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
				_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)

				if tagExists {
					tags[tagID].TimesUsed = timesUsed
				} else {
					tags = append(tags, newTag)
				}
				newTagJson, _ := json.MarshalIndent(tags, "", " ")
				_ = ioutil.WriteFile("tags.json", newTagJson, 0644)
			}

			fmt.Fprint(w, "Succsess")
		} else {
			fmt.Fprint(w, "ERROR - No Tag submitted")
		}
	} else {
		fmt.Fprint(w, "ERROR - No Aktion")
	}
}

func rmTag(w http.ResponseWriter, r *http.Request) {
	if aktionLive {
		tagToRemove := onlyName{
			Name: r.FormValue("tag"),
		}

		if tagToRemove.Name != "" {
			fmt.Printf("Tag to Delete: %+v", tagToRemove)
			tagExisted := false //Ist der Tag bereits hinzugefügt worden?
			aktionenJson, _ := ioutil.ReadFile("aktionen.json")
			var aktionen []Aktion
			json.Unmarshal([]byte(aktionenJson), &aktionen)
			filteredTags := make([]onlyName, 0)
			for _, tag := range aktionen[currentActionID].Tags {
				if tagToRemove.Name != tag.Name {
					filteredTags = append(filteredTags, tag)
				} else {
					tagExisted = true
				}
			}

			aktionen[currentActionID].Tags = filteredTags

			if tagExisted {
				tagsJson, _ := ioutil.ReadFile("tags.json")
				var tags []Tag
				json.Unmarshal([]byte(tagsJson), &tags)
				for i, tag := range tags {
					if tagToRemove.Name == tag.Name {
						if tag.TimesUsed > 0 {
							tags[i].TimesUsed -= 1
						}
						break
					}
				}

				newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
				_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)

				newTagJson, _ := json.MarshalIndent(tags, "", " ")
				_ = ioutil.WriteFile("tags.json", newTagJson, 0644)

				fmt.Fprint(w, "Succsess")
			} else {
				fmt.Fprint(w, "ERROR - Tag did not even exist")
			}
		} else {
			fmt.Fprint(w, "ERROR - No Tag submitted")
		}
	} else {
		fmt.Fprint(w, "ERROR - No Aktion")
	}
}

func editAktionInfotext(w http.ResponseWriter, r *http.Request) {
	if aktionLive {
		infotext := r.FormValue("text")

		aktionenJson, _ := ioutil.ReadFile("aktionen.json")
		var aktionen []Aktion
		json.Unmarshal([]byte(aktionenJson), &aktionen)
		aktionen[currentActionID].InfoText = infotext

		newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
		_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)
	}

}

func addPerson(w http.ResponseWriter, r *http.Request) {
	if aktionLive {
		newVisitor := Person{
			Name:         r.FormValue("person"),
			TimesVisited: 1,
		}

		if newVisitor.Name != "" {
			fmt.Printf("New Visitor: %+v", newVisitor)

			visitorAdded := false //Ist der Tag bereits hinzugefügt worden?
			aktionenJson, _ := ioutil.ReadFile("aktionen.json")
			var aktionen []Aktion
			json.Unmarshal([]byte(aktionenJson), &aktionen)
			currentAction := aktionen[currentActionID]
			for _, tag := range currentAction.Visitors {
				if newVisitor.Name == tag.Name {
					visitorAdded = true
					break
				}
			}

			visitorExists := false //Ist der Tag bereits verwendet worden?
			var visitorID int
			var timesVisited int = 0
			tagsJson, _ := ioutil.ReadFile("personen.json")
			var visitors []Person
			json.Unmarshal([]byte(tagsJson), &visitors)
			for index, person := range visitors {
				if newVisitor.Name == person.Name {
					visitorExists = true
					timesVisited = person.TimesVisited
					visitorID = index
					break
				}
			}

			if !visitorAdded {
				aktionen[currentActionID].Visitors = append(aktionen[currentActionID].Visitors, onlyName{Name: newVisitor.Name})
				timesVisited++

				newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
				_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)

				if visitorExists {
					visitors[visitorID].TimesVisited = timesVisited
				} else {
					visitors = append(visitors, newVisitor)
				}
				newTagJson, _ := json.MarshalIndent(visitors, "", " ")
				_ = ioutil.WriteFile("personen.json", newTagJson, 0644)
			}

			fmt.Fprint(w, "Succsess")
		} else {
			fmt.Fprint(w, "ERROR - No Tag submitted")
		}
	} else {
		fmt.Fprint(w, "ERROR - No Aktion")
	}
}

func rmPerson(w http.ResponseWriter, r *http.Request) {
	if aktionLive {
		visitorToDelete := onlyName{
			Name: r.FormValue("person"),
		}

		if visitorToDelete.Name != "" {
			fmt.Printf("Delete Visitor: %+v", visitorToDelete)

			visitorExisted := false //Ist der Tag bereits hinzugefügt worden?
			aktionenJson, _ := ioutil.ReadFile("aktionen.json")
			var aktionen []Aktion
			json.Unmarshal([]byte(aktionenJson), &aktionen)
			filteredPersons := make([]onlyName, 0)
			for _, person := range aktionen[currentActionID].Visitors {
				if visitorToDelete.Name != person.Name {
					filteredPersons = append(filteredPersons, person)
				} else {
					visitorExisted = true
				}
			}

			aktionen[currentActionID].Visitors = filteredPersons

			if visitorExisted {
				personsJson, _ := ioutil.ReadFile("personen.json")
				var visitors []Person
				json.Unmarshal([]byte(personsJson), &visitors)
				for i, person := range visitors {
					if visitorToDelete.Name == person.Name {
						if person.TimesVisited > 0 {
							visitors[i].TimesVisited--
						}
						break
					}
				}

				newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
				_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)

				newVisitorsJson, _ := json.MarshalIndent(visitors, "", " ")
				_ = ioutil.WriteFile("personen.json", newVisitorsJson, 0644)

				fmt.Fprint(w, "Succsess")
			} else {
				fmt.Fprint(w, "ERROR - Tag dit not exist")
			}
		} else {
			fmt.Fprint(w, "ERROR - No Tag submitted")
		}
	} else {
		fmt.Fprint(w, "ERROR - No Aktion")
	}
}

func changeAktionName(w http.ResponseWriter, r *http.Request) {
	if aktionLive {
		if r.FormValue("name") != "" {
			aktionenJson, _ := ioutil.ReadFile("aktionen.json")
			var aktionen []Aktion
			json.Unmarshal([]byte(aktionenJson), &aktionen)
			aktionen[currentActionID].Title = r.FormValue("name")
			newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
			_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)
		} else {
			fmt.Fprint(w, "ERROR - No Tag submitted")
		}
	} else {
		fmt.Fprint(w, "ERROR - No Aktion")
	}
}

func toggleSortAfterTags(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.Diashowfilter.SortAfterTags = !settings.Diashowfilter.SortAfterTags

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func toggleSortAfterPersons(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.Diashowfilter.SortAfterPersons = !settings.Diashowfilter.SortAfterPersons

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func toggleSortAfterTimeEnd(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.Diashowfilter.SortAfterTimeEnd = !settings.Diashowfilter.SortAfterTimeEnd

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func toggleSortAfterTimeStart(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.Diashowfilter.SortAfterTimeStart = !settings.Diashowfilter.SortAfterTimeStart

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func toggleOnlyAdditionalContent(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.Diashowfilter.OnlyExternContent = !settings.Diashowfilter.OnlyExternContent

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func toggleAllowTagEditing(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.AllowTagEditing = !settings.AllowTagEditing

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
}

func toggleAllowInfoTextEditing(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.AllowInfoTextEditing = !settings.AllowInfoTextEditing

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
}

func toggleAllowPersonEditing(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	settings.AllowPersonEditing = !settings.AllowPersonEditing

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
}

func toggleFreeUpload(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	index, err := strconv.Atoi(r.FormValue("index"))
	if err != nil {
		fmt.Println(err)
	}

	settings.AdditionalPics[index].FreeUpload = !settings.AdditionalPics[index].FreeUpload

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
}

func toggleShowExternContent(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	index, err := strconv.Atoi(r.FormValue("index"))
	if err != nil {
		fmt.Println(err)
	}

	settings.AdditionalPics[index].Show = !settings.AdditionalPics[index].Show

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
}

func editExternContentInfoText(w http.ResponseWriter, r *http.Request) {

	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	index, err := strconv.Atoi(r.FormValue("index"))

	if err != nil {
		fmt.Println(err)
	}

	infotext := r.FormValue("text")
	if infotext != "" {
		settings.AdditionalPics[index].InfoText = infotext

		newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
		_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
	}
}

func setFilterTimeStart(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	var newTime TimeforJson
	if r.FormValue("Year") != "" {
		newTime.Year, _ = strconv.Atoi(r.FormValue("Year"))
		if r.FormValue("Month") != "" {
			newTime.Month, _ = strconv.Atoi(r.FormValue("Month"))
		} else {
			newTime.Month = 1
		}
		if r.FormValue("Day") != "" {
			newTime.Day, _ = strconv.Atoi(r.FormValue("Day"))
		} else {
			newTime.Day = 1
		}

		settings.Diashowfilter.TimeStart = newTime
		newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
		_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
	}
}

func setFilterTimeEnd(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	var newTime TimeforJson
	if r.FormValue("Year") != "" {
		newTime.Year, _ = strconv.Atoi(r.FormValue("Year"))
		if r.FormValue("Month") != "" {
			newTime.Month, _ = strconv.Atoi(r.FormValue("Month"))
		} else {
			newTime.Month = 1
		}
		if r.FormValue("Day") != "" {
			newTime.Day, _ = strconv.Atoi(r.FormValue("Day"))
		} else {
			newTime.Day = 1
		}

		settings.Diashowfilter.TimeEnd = newTime
		newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
		_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
	}
}

func addAdditionalContent(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	var newAdditionalContent additionalFolder
	if r.FormValue("Title") != "" {
		newAdditionalContent.Name = r.FormValue("Title")
		newAdditionalContent.InfoText = r.FormValue("InfoText")
		newAdditionalContent.Show = false
		newAdditionalContent.FreeUpload = false

		newAdditionalContent.Path = "/" + strings.ReplaceAll(newAdditionalContent.Name, " ", "")
		os.Mkdir("./externContent"+newAdditionalContent.Path, 0750)

		settings.AdditionalPics = append(settings.AdditionalPics, newAdditionalContent)
		newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
		_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
	}

}

func deleteAdditionalContent(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.FormValue("folder"))
	if err != nil {
		fmt.Print(err)
		return
	}

	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	os.RemoveAll("./externContent" + settings.AdditionalPics[index].Path)

	var newAdditionalContents []additionalFolder

	for i, folder := range settings.AdditionalPics {
		if i != index {
			newAdditionalContents = append(newAdditionalContents, folder)
		}
	}

	settings.AdditionalPics = newAdditionalContents

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
}

func deletePicFromAdditionalContent(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.FormValue("folder"))
	if err != nil {
		fmt.Print(err)
		return
	}

	picName := r.FormValue("picName")

	if picName != "" {
		fmt.Println(picName)
		settingsJson, _ := ioutil.ReadFile("settings.json")
		var settings Settings
		json.Unmarshal([]byte(settingsJson), &settings)
		// Achtung sicherheitsLücke, Überprüfen, ob Bild existiert
		os.Remove("./externContent" + settings.AdditionalPics[index].Path + "/" + picName)
	}
}

func deletePicFromAction(w http.ResponseWriter, r *http.Request) {
	folder := r.FormValue("folder")

	fmt.Println(folder)

	actionJson, _ := ioutil.ReadFile("aktionen.json")
	var aktionen []Aktion
	json.Unmarshal([]byte(actionJson), &aktionen)

	var pictureAktion Aktion
	actionExists := false
	for _, action := range aktionen {
		if action.Folder == folder {
			actionExists = true
			pictureAktion = action
		}
	}

	picName := r.FormValue("pic")
	fmt.Println(actionExists)

	if picName != "" && actionExists {
		fmt.Println(picName)
		// Achtung sicherheitsLücke, Überprüfen, ob Bild existiert
		os.Remove("./aktionen/" + pictureAktion.Folder + picName)
	}
}

func ignoreTag(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	name := r.FormValue("name")
	var filteredTags []onlyName
	if name != "" {
		for _, tag := range settings.Diashowfilter.Tags {
			if tag.Name != name {
				filteredTags = append(filteredTags, tag)
			}
		}
	}
	if len(filteredTags) == 0 {
		filteredTags = []onlyName{}
	}
	settings.Diashowfilter.Tags = filteredTags
	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func ignorePerson(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	name := r.FormValue("name")
	var filteredPersons []onlyName
	if name != "" {
		for _, person := range settings.Diashowfilter.Persons {
			if person.Name != name {
				filteredPersons = append(filteredPersons, person)
			}
		}
	}
	if len(filteredPersons) == 0 {
		filteredPersons = []onlyName{}
	}
	settings.Diashowfilter.Persons = filteredPersons
	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)
}

func setTag(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	name := r.FormValue("name")
	if name != "" {
		var newTag onlyName
		newTag.Name = name
		settings.Diashowfilter.Tags = append(settings.Diashowfilter.Tags, newTag)
	}

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func setPerson(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)
	name := r.FormValue("name")
	if name != "" {
		var newPerson onlyName
		newPerson.Name = name
		settings.Diashowfilter.Persons = append(settings.Diashowfilter.Persons, newPerson)
	}

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func setDiashowSpeed(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)
	diashowSpeed, err := strconv.Atoi(r.FormValue("DiashowSpeed"))
	if err != nil {
		fmt.Println(err)
	}
	if r.FormValue("DiashowSpeed") != "" {
		settings.DiashowSpeed = diashowSpeed
	}

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func setAutoActionStop(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	autoActionStop, err := strconv.Atoi(r.FormValue("AutoActionStop"))
	if err != nil {
		fmt.Println(err)
	}
	if r.FormValue("AutoActionStop") != "" {
		settings.AutoActionStop = autoActionStop
	}

	newSettingsJson, _ := json.MarshalIndent(settings, "", " ")
	_ = ioutil.WriteFile("settings.json", newSettingsJson, 0644)

}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	uhrzeit := time.Now()
	timeDiff := uhrzeit.Sub(LastTimeUploaded)
	if int(timeDiff.Minutes()) > settings.AutoActionStop {
		stopAktion()
	}

	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("imageFile")
	if err == nil {
		defer file.Close()

		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		var endung string
		if strings.Index(handler.Filename, ".png") > -1 {
			endung = ".png"
		} else if strings.Index(handler.Filename, ".jpg") > -1 {
			endung = ".jpg"
		}

		if !aktionLive {
			startAktion()
		}

		uhrzeit := time.Now()
		LastTimeUploaded = uhrzeit

		var filename string = strconv.Itoa(uhrzeit.Year()) + "_" + strconv.Itoa(int(uhrzeit.Month())) + "_" + strconv.Itoa(uhrzeit.Day()) + "_" + strconv.Itoa(uhrzeit.Hour()) + "_" + strconv.Itoa(uhrzeit.Minute())

		aktionenJson, _ := ioutil.ReadFile("aktionen.json")
		var aktionen []Aktion
		json.Unmarshal([]byte(aktionenJson), &aktionen)

		tempFile, err := ioutil.TempFile("./aktionen"+aktionen[currentActionID].Folder, filename+"-*"+endung)
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()

		// read all of the contents of our uploaded file into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
		}
		// write this byte array to our temporary file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		http.Redirect(w, r, "../upload", http.StatusSeeOther)
	} else {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
	}

}

func uploadToAdditionalContent(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	index, err := strconv.Atoi(r.FormValue("folder"))
	if err != nil {
		fmt.Println(err)
	}

	if len(settings.AdditionalPics) > index {
		folder := "/externContent" + settings.AdditionalPics[index].Path

		file, handler, err := r.FormFile("imageFile")
		if err == nil {
			defer file.Close()

			fmt.Printf("Uploaded File: %+v\n", handler.Filename)
			fmt.Printf("File Size: %+v\n", handler.Size)
			fmt.Printf("MIME Header: %+v\n", handler.Header)

			tempFile, err := ioutil.TempFile("."+folder, "*"+handler.Filename)
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()

			// read all of the contents of our uploaded file into a
			// byte array
			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}
			// write this byte array to our temporary file
			tempFile.Write(fileBytes)
			// return that we have successfully uploaded our file!
			fromUploadPage := r.FormValue("fromUploadPage")
			fmt.Println(fromUploadPage)
			if fromUploadPage == "true" {
				http.Redirect(w, r, "../uploadToExtern"+settings.AdditionalPics[index].Path, http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "../settings", http.StatusSeeOther)
			}
		} else {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
		}
	}

}

func returnDiashow(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	file, _ := ioutil.ReadFile("aktionen.json")
	var actions []Aktion
	json.Unmarshal([]byte(file), &actions)

	var diashow Diashow
	diashow.DiashowSpeed = settings.DiashowSpeed

	var involvedFolders []additionalFolder
	for _, folder := range settings.AdditionalPics {
		if folder.Show {
			involvedFolders = append(involvedFolders, folder)
		}
	}

	involvedAktions := actions
	var aktionsBuffer []Aktion

	if settings.Diashowfilter.SortAfterPersons {
		aktionsBuffer = involvedAktions
		involvedAktions = []Aktion{}
		for _, aktion := range aktionsBuffer {
			hasPerson := false
			for _, person := range settings.Diashowfilter.Persons {
				for _, visitor := range aktion.Visitors {
					if visitor.Name == person.Name {
						hasPerson = true
						break
					}
					if hasPerson {
						break
					}
				}
				if hasPerson {
					break
				}
			}
			if hasPerson {
				involvedAktions = append(involvedAktions, aktion)
			}
		}
	}

	if settings.Diashowfilter.SortAfterTags {
		aktionsBuffer = involvedAktions
		involvedAktions = []Aktion{}
		for _, aktion := range aktionsBuffer {
			hasTag := false
			for _, tagToSearch := range settings.Diashowfilter.Tags {
				for _, tag := range aktion.Tags {
					if tag.Name == tagToSearch.Name {
						hasTag = true
						break
					}
					if hasTag {
						break
					}
				}
				if hasTag {
					break
				}
			}
			if hasTag {
				involvedAktions = append(involvedAktions, aktion)
			}
		}
	}

	if settings.Diashowfilter.SortAfterTimeEnd {
		aktionsBuffer = involvedAktions
		involvedAktions = []Aktion{}
		for _, aktion := range aktionsBuffer {
			toLate := false
			if aktion.Start.Year > settings.Diashowfilter.TimeEnd.Year {
				toLate = true
			} else if aktion.Start.Year == settings.Diashowfilter.TimeEnd.Year {
				if aktion.Start.Month > settings.Diashowfilter.TimeEnd.Month {
					toLate = true
				} else if aktion.Start.Month == settings.Diashowfilter.TimeEnd.Month {
					if aktion.Start.Day > settings.Diashowfilter.TimeEnd.Day {
						toLate = true
					}
				}
			}
			if !toLate {
				involvedAktions = append(involvedAktions, aktion)
			}
		}
	}

	if settings.Diashowfilter.SortAfterTimeStart {
		aktionsBuffer = involvedAktions
		involvedAktions = []Aktion{}
		for _, aktion := range aktionsBuffer {
			toEarly := false
			if aktion.End.Year < settings.Diashowfilter.TimeStart.Year {
				toEarly = true
			} else if aktion.End.Year == settings.Diashowfilter.TimeStart.Year {
				if aktion.End.Month < settings.Diashowfilter.TimeStart.Month {
					toEarly = true
				} else if aktion.End.Month == settings.Diashowfilter.TimeStart.Month {
					if aktion.End.Day < settings.Diashowfilter.TimeStart.Day {
						toEarly = true
					}
				}
			}
			if !toEarly {
				involvedAktions = append(involvedAktions, aktion)
			}
		}
	}

	var folder string

	if (len(involvedFolders) > 0) || ((len(involvedAktions) > 0) && !settings.Diashowfilter.OnlyExternContent) {
		randomInt := rand.Int()
		if settings.Diashowfilter.OnlyExternContent {
			index := randomInt % len(involvedFolders)
			diashow.Folder = involvedFolders[index]
			folder = "./externContent" + diashow.Folder.Path
			diashow.FromAktion = false
		} else {
			index := randomInt % (len(involvedAktions) + len(involvedFolders))
			if index < len(involvedFolders) {
				diashow.Folder = involvedFolders[index]
				folder = "./externContent" + diashow.Folder.Path
				diashow.FromAktion = false
			} else {
				diashow.Aktion = involvedAktions[index-(len(involvedFolders))]
				folder = "./aktionen" + diashow.Aktion.Folder
				diashow.FromAktion = true
			}
		}

		pics, _ := ioutil.ReadDir(folder)
		if len(pics) > 0 {
			picID := rand.Int() % len(pics)
			diashow.Picture = folder + "/" + pics[picID].Name()
		} else {
			diashow.FromAktion = true
			diashow.Aktion.Title = "Empty Folder"
			diashow.Picture = "./static/camera.png"
			diashow.DiashowSpeed = 0
		}

	} else {
		diashow.FromAktion = false
		diashow.Folder.Name = "Empty Folder"
		diashow.Folder.InfoText = "Es scheint keine Bilder zum anzeigen zu geben, checke deine Diashow Filter oder lade Bilder Hoch."
		diashow.Picture = "./static/camera.png"
		diashow.DiashowSpeed = 15
	}

	diashowJson, err := json.Marshal(diashow)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprint(w, string(diashowJson))
}

func getExternPics(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.Atoi(r.FormValue("folder"))

	if err != nil {
		log.Panic(err)
	}

	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	if index >= len(settings.AdditionalPics) {
		log.Panic("Content doesn't exist")
		return
	}

	pics, _ := ioutil.ReadDir("./externContent" + settings.AdditionalPics[index].Path)

	var picPaths []onlyPath
	for _, pic := range pics {
		picPaths = append(picPaths, onlyPath{Path: pic.Name()})
	}

	//fmt.Println(picPaths)

	myJson, err := json.Marshal(picPaths)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprint(w, string(myJson))
}

func returnTags(w http.ResponseWriter, r *http.Request) {
	tagJson, _ := ioutil.ReadFile("tags.json")
	var tags []Tag
	json.Unmarshal([]byte(tagJson), &tags)

	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i].TimesUsed > tags[j].TimesUsed
	})

	myJson, err := json.Marshal(tags)
	if err != nil {
		log.Panic(err)
	}

	fmt.Fprintf(w, string(myJson))
}

func returnSettings(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	myJson, err := json.Marshal(settings)
	if err != nil {
		log.Panic(err)
	}

	fmt.Fprintf(w, string(myJson))
}

func returnActions(w http.ResponseWriter, r *http.Request) {
	aktionJson, _ := ioutil.ReadFile("aktionen.json")
	var actions []Aktion
	json.Unmarshal([]byte(aktionJson), &actions)

	myJson, err := json.Marshal(actions)
	if err != nil {
		log.Panic(err)
	}

	fmt.Fprintf(w, string(myJson))
}

func returnPictures(w http.ResponseWriter, r *http.Request) {
	folder := r.FormValue("Folder")
	if folder != "" {
		aktionJson, _ := ioutil.ReadFile("aktionen.json")
		var actions []Aktion
		json.Unmarshal([]byte(aktionJson), &actions)

		pics, _ := ioutil.ReadDir("./aktionen" + folder)
		var picsToReturn []onlyPath

		for _, pic := range pics {
			var newPic onlyPath
			newPic.Path = pic.Name()
			picsToReturn = append(picsToReturn, newPic)
		}

		myJson, err := json.Marshal(picsToReturn)
		if err != nil {
			log.Panic(err)
		}

		fmt.Fprintf(w, string(myJson))
	} else {
		fmt.Fprintf(w, "Aktion Folder not found")
	}

}

func returnPersons(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("personen.json")
	var persons []Person
	json.Unmarshal([]byte(file), &persons)

	myJson, err := json.Marshal(persons)
	if err != nil {
		log.Panic(err)
	}

	sort.SliceStable(persons, func(i, j int) bool {
		return persons[i].TimesVisited > persons[j].TimesVisited
	})

	fmt.Fprintf(w, string(myJson))
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func returnStats(w http.ResponseWriter, r *http.Request) {
	var stat Statistic

	personFile, _ := ioutil.ReadFile("personen.json")
	json.Unmarshal([]byte(personFile), &stat.Visitors)

	sort.SliceStable(stat.Visitors, func(i, j int) bool {
		return stat.Visitors[i].TimesVisited > stat.Visitors[j].TimesVisited
	})

	stat.VisitorCount = len(stat.Visitors)

	tagFile, _ := ioutil.ReadFile("tags.json")
	json.Unmarshal([]byte(tagFile), &stat.Tags)

	sort.SliceStable(stat.Tags, func(i, j int) bool {
		return stat.Tags[i].TimesUsed > stat.Tags[j].TimesUsed
	})

	stat.TagCount = len(stat.Tags)

	aktionenJson, _ := ioutil.ReadFile("aktionen.json")
	var aktionen []Aktion
	json.Unmarshal([]byte(aktionenJson), &aktionen)

	stat.ActionCount = len(aktionen)
	stat.ActionPictureCount = 0
	for _, aktion := range aktionen {
		pics, _ := ioutil.ReadDir("./aktionen" + aktion.Folder)
		stat.ActionPictureCount += len(pics)
	}

	stat.DataSize, _ = DirSize("./")
	stat.DataSize /= 1048576

	stat.LastAktion = aktionen[len(aktionen)-1]

	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	stat.ExternPictureCount = 0
	for _, addCnt := range settings.AdditionalPics {
		pics, _ := ioutil.ReadDir("./externContent" + addCnt.Path)
		stat.ExternPictureCount += len(pics)
	}

	myJson, err := json.Marshal(stat)
	if err != nil {
		log.Panic(err)
	}
	fmt.Fprintf(w, string(myJson))
}

func returnCurrentAction(w http.ResponseWriter, r *http.Request) {
	settingsJson, _ := ioutil.ReadFile("settings.json")
	var settings Settings
	json.Unmarshal([]byte(settingsJson), &settings)

	uhrzeit := time.Now()
	timeDiff := uhrzeit.Sub(LastTimeUploaded)
	if int(timeDiff.Minutes()) > settings.AutoActionStop {
		stopAktion()
	}

	if aktionLive {
		file, _ := ioutil.ReadFile("aktionen.json")
		var actions []Aktion
		json.Unmarshal([]byte(file), &actions)

		//fmt.Printf("Aktionen 187:%+v", actions)

		aktionJson, err := json.Marshal(actions[currentActionID])
		if err != nil {
			log.Panic(err)
		}
		fmt.Fprintf(w, string(aktionJson))
	} else {
		aktionNotLiveJson, _ := json.Marshal(map[string]bool{"aktionNotLive": true})
		fmt.Fprintf(w, string(aktionNotLiveJson))
	}
}

func deleteAction(w http.ResponseWriter, r *http.Request) {
	folder := r.FormValue("folder")

	actionsJson, _ := ioutil.ReadFile("aktionen.json")
	var actions []Aktion
	json.Unmarshal([]byte(actionsJson), &actions)

	actionExists := false
	var newActions []Aktion
	for _, action := range actions {
		if action.Folder == folder {
			actionExists = true
		} else {
			newActions = append(newActions, action)
		}
	}

	if actionExists {
		os.RemoveAll("./aktionen" + folder)
		newActionssJson, _ := json.MarshalIndent(newActions, "", " ")
		_ = ioutil.WriteFile("aktionen.json", newActionssJson, 0644)
	}
}

func changeAktion(w http.ResponseWriter, r *http.Request) {
	actionStr := r.FormValue("Aktion")
	//fmt.Println(actionStr)
	var actionToChange Aktion
	json.Unmarshal([]byte(actionStr), &actionToChange)

	file, _ := ioutil.ReadFile("aktionen.json")
	var actions []Aktion
	json.Unmarshal([]byte(file), &actions)

	var index int
	for i, action := range actions {
		if action.Folder == actionToChange.Folder {
			index = i
		}
	}

	var persons []Person
	personFile, _ := ioutil.ReadFile("personen.json")
	json.Unmarshal([]byte(personFile), &persons)

	for _, visitor := range actions[index].Visitors {
		visitorDeleted := true
		for _, person := range actionToChange.Visitors {
			if person.Name == visitor.Name {
				visitorDeleted = false
			}
		}
		if visitorDeleted {
			for i, person := range persons {
				if person.Name == visitor.Name {
					if person.TimesVisited > 0 {
						persons[i].TimesVisited = persons[i].TimesVisited - 1
					} else {
						persons[i].TimesVisited = 0
					}
				}
			}
		}
	}

	for _, visitor := range actionToChange.Visitors {
		visitorAdded := true
		for _, person := range actions[index].Visitors {
			if person.Name == visitor.Name {
				visitorAdded = false
			}
		}
		if visitorAdded {
			for i, person := range persons {
				if person.Name == visitor.Name {
					persons[i].TimesVisited = persons[i].TimesVisited + 1
				}
			}
		}
	}

	newPersonJson, _ := json.MarshalIndent(persons, "", " ")
	_ = ioutil.WriteFile("personen.json", newPersonJson, 0644)

	var tags []Tag
	tagFile, _ := ioutil.ReadFile("tags.json")
	json.Unmarshal([]byte(tagFile), &tags)

	for _, maybeDeletedTag := range actions[index].Tags {
		tagDeleted := true
		for _, tag := range actionToChange.Tags {
			if tag.Name == maybeDeletedTag.Name {
				tagDeleted = false
			}
		}
		if tagDeleted {
			for i, tag := range tags {
				if maybeDeletedTag.Name == maybeDeletedTag.Name {
					if tag.TimesUsed > 0 {
						tags[i].TimesUsed--
					} else {
						tags[i].TimesUsed = 0
					}
				}
			}
		}
	}

	for _, maybeAddedTag := range actionToChange.Tags {
		tagAdded := true
		for _, tag := range actions[index].Tags {
			if tag.Name == maybeAddedTag.Name {
				tagAdded = false
			}
		}
		if tagAdded {
			for i, tag := range tags {
				if tag.Name == maybeAddedTag.Name {
					tags[i].TimesUsed = tag.TimesUsed + 1
				}
			}
		}
	}

	//fmt.Println(tags)
	//fmt.Println(persons)

	newTagJson, _ := json.MarshalIndent(tags, "", " ")
	_ = ioutil.WriteFile("tags.json", newTagJson, 0644)

	actions[index] = actionToChange
	newActionJson, _ := json.MarshalIndent(actions, "", " ")
	_ = ioutil.WriteFile("aktionen.json", newActionJson, 0644)
}

func returnIndexHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/dist/index.html")
}

func hostFrontend() {
	// Serve static files from the frontend/dist directory.
	fs := http.FileServer(http.Dir("./frontend/dist"))
	http.Handle("/", fs)
}

func hostStatic() {
	// Serve static files from the frontend/dist directory.
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
}

func hostActions() {
	// Serve static files from the frontend/dist directory.
	fs := http.FileServer(http.Dir("./aktionen"))
	http.Handle("/pics/", http.StripPrefix("/pics", fs))
	http.Handle("/aktionen/", http.StripPrefix("/aktionen", fs))
}

func startAktion() {
	uhrzeit := time.Now()
	LastTimeUploaded = uhrzeit

	var newAction Aktion
	newAction.Start.Year = uhrzeit.Year()
	newAction.Start.Month = int(uhrzeit.Month())
	newAction.Start.Day = uhrzeit.Day()
	newAction.Start.Hour = uhrzeit.Hour()
	newAction.Start.Minute = uhrzeit.Minute()
	emptyOnlyname := []onlyName{}
	newAction.Visitors = emptyOnlyname
	newAction.Tags = emptyOnlyname
	newAction.InfoText = ""

	newAction.Folder = "/" + strconv.Itoa(uhrzeit.Year()) + "/" + strconv.Itoa(int(uhrzeit.Month())) + "/" + strconv.Itoa(int(uhrzeit.Day())) + "-" + strconv.Itoa(int(rand.Int()%10000)) + "/"
	os.Mkdir("./aktionen"+"/"+strconv.Itoa(uhrzeit.Year()), 0750)
	os.Mkdir("./aktionen"+"/"+strconv.Itoa(uhrzeit.Year())+"/"+strconv.Itoa(int(uhrzeit.Month())), 0750)
	os.Mkdir("./aktionen"+newAction.Folder, 0750)
	fmt.Println(newAction.Folder)

	aktionenJson, _ := ioutil.ReadFile("aktionen.json")
	var aktionen []Aktion
	json.Unmarshal([]byte(aktionenJson), &aktionen)

	aktionen = append(aktionen, newAction)
	newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
	_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)

	currentActionID = len(aktionen) - 1
	aktionLive = true
}

func stopAktion() {
	if aktionLive {
		aktionenJson, _ := ioutil.ReadFile("aktionen.json")
		var aktionen []Aktion
		json.Unmarshal([]byte(aktionenJson), &aktionen)

		aktionen[currentActionID].End.Year = LastTimeUploaded.Year()
		aktionen[currentActionID].End.Month = int(LastTimeUploaded.Month())
		aktionen[currentActionID].End.Day = LastTimeUploaded.Day()
		aktionen[currentActionID].End.Hour = LastTimeUploaded.Hour()
		aktionen[currentActionID].End.Minute = LastTimeUploaded.Minute()

		newAktionJson, _ := json.MarshalIndent(aktionen, "", " ")
		_ = ioutil.WriteFile("aktionen.json", newAktionJson, 0644)

		aktionLive = false
	}
}

func stopAktionWeb(w http.ResponseWriter, r *http.Request) {
	print("stopped over web")
	stopAktion()
}
