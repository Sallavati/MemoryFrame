<script>
export default {
    data() {
    return {
      trashIcon: "./static/trash.png",
      searchAktion: "",
      timeFilterStart: "",
      timeFilterStop: "",
      actions: [],
      years: [],
      edit: false,
      actionToedit: {'Tags': [], 'Visitors': []},
      allTags: [],
      tagSearch: "",
      filteredTags: null,
      personSearch: "",
      filteredPersons: null,
      allPersons: null,
      editDateStart: "",
      editDateEnd: "",
      editTimeStart: "",
      editTimeEnd: "",
    }
  },
  methods: {
    async fetchAllPersons() {
      const res = await fetch(
        `./api/persons`
      )
      this.allPersons = await res.json()
      this.filterPersons()
    },
    filterPersons(){
      this.filteredPersons = this.allPersons.filter((t) => {
        let exists = false
        if(!t.Name.includes(this.personSearch)){
            return false
        }
        this.actionToedit.Visitors.forEach(element => {
          if(t.Name == element.Name){
            exists = true
          }
        });
        return !exists;
      });
      this.filteredPersons = this.filteredPersons.slice(0, 14)
    },
    async fetchAllTags() {
      this.allTags = null
      const res = await fetch(
        `./api/tags`
      )
      this.allTags = await res.json()
      this.filterTags()
    },
    showAktionPics(action){
        action.showPics = true
        this.fetchPics(action.Folder)
    },
    async fetchPics(folder){
      let formData = new URLSearchParams();
      formData.append('Folder', folder);
      const res = await fetch(
        `./api/actionPictures`,
        {
          method: 'post',
          body: formData
      })
      let index;
      this.actions.forEach((e, i) => {
        if(e.Folder == folder){
            index = i;
        }
      })
      let pictures = await res.json();
      this.actions[index].pics = pictures

      if(this.actionToedit.Folder == folder){
        this.actionToedit.pics = pictures
      }
    },
    filterTags(){
      this.filteredTags = this.allTags.filter((t) => {
        let exists = false
         if(!t.Name.includes(this.tagSearch)){
            return false
        }
        this.actionToedit.Tags.forEach(element => {
          if(t.Name == element.Name){
            exists = true
          }
        });
        return !exists;
      });
      this.filteredTags = this.filteredTags.slice(0, 14)
    },
    async deletePicturefromAktion(folder, picpath){
      if(!window.confirm("Sicher, dass sie das Bild löschen wollen?")){
        return
      }
      let formData = new URLSearchParams();
      formData.append('folder', folder);
      formData.append('pic', picpath);
      const res = await fetch(
        `./api/deletePictureFromAction`,
        {
          method: 'post',
          body: formData
      })
      this.fetchPics(folder)
    },
    async deleteAction(){
      if(!window.confirm("sicher, dass sie die Aktion löschen wollen?")){
        return
      }
        let formData = new URLSearchParams();
        formData.append('folder', this.actionToedit.Folder);
        const res = await fetch(
        `./api/deleteAktion`,
        {
          method: 'post',
          body: formData
      })
      this.actions = this.actions.filter(action => action.Folder != this.actionToedit.Folder)
      this.edit = false
    },
    async actualizeAction(){
        //console.log(this.actionToedit)
        let formData = new URLSearchParams();
        this.actionToedit.Start = this.timeNDateStringToJson(this.editTimeStart, this.editDateStart)
        this.actionToedit.End = this.timeNDateStringToJson(this.editTimeEnd, this.editDateEnd)
        formData.append('Aktion', JSON.stringify(this.actionToedit));
        const res = await fetch(
        `./api/changeAction`,
        {
          method: 'post',
          body: formData
      })
      this.actions.forEach((action) => {
        if(action.Folder == this.actionToedit.Folder){
            action = Object.assign({}, this.actionToedit)
        }
      })
      await this.fetchActions()
      this.edit = false
    },
    addPersonToAction(name){
      this.actionToedit.Visitors.push({Name: name})
      this.filterPersons()
    },
    rmPersonFromAction(name){
      this.actionToedit.Visitors = this.actionToedit.Visitors.filter(visitor => visitor.Name != name)
      this.filterPersons()
    },
    addTagToAction(name){
      this.actionToedit.Tags.push({Name: name})
      this.filterTags()
    },
    rmTagFromAction(name){
      this.actionToedit.Tags = this.actionToedit.Tags.filter(tag => tag.Name != name)
      this.filterTags()
    },
    hidePics(folder){
        let index;
        this.actions.forEach((e, i) => {
            if(e.Folder == folder){
                index = i;
            }
        })
        this.actions[index].showPics = false;
    },
    async fetchActions() {
      this.allPersons = null
      const res = await fetch(
        `./api/actions`
      )
      this.actions = await res.json()
      this.actions.reverse()
      this.filterYears()
    },
    timeNDateStringToJson(timeString, dateString){
      return {
        Year: Number(dateString.slice(0, 4)), 
        Month: Number(dateString.slice(5, 7)), 
        Day: Number(dateString.slice(8, 10)),
        Hour: Number(timeString.slice(0, 2)),
        Minute: Number(timeString.slice(3, 5))
      }
    },
    timeJsonToString(timeJson){
      return `${timeJson.Hour < 10 ? "0" : ""}${timeJson.Hour}:${timeJson.Minute < 10 ? "0" : ""}${timeJson.Minute}`
    },
    dateJsonToString(dateJson){
      return `${dateJson.Year < 10 ? "000" : ""}${dateJson.Year}-${dateJson.Month < 10 ? "0" : ""}${dateJson.Month}-${dateJson.Day < 10 ? "0" : ""}${dateJson.Day}`
    },
    async editArticle(action){
        //console.log(action)
        await this.fetchAllTags()
        await this.fetchAllPersons()
        await this.fetchPics(action.Folder)
        this.actionToedit = Object.assign({}, action)
        this.filterPersons()
        this.filterTags()
        this.edit = true
        this.editDateEnd = this.dateJsonToString(action.End)
        this.editDateStart = this.dateJsonToString(action.Start)
        this.editTimeStart = this.timeJsonToString(action.Start)
        this.editTimeEnd = this.timeJsonToString(action.End)
    },
    stopEditing(){
        this.edit = false 
        this.actionToedit={Tags: [], Visitors: [], Folder: ""}
    },
    resetSearch(){
        this.searchAktion = ""
        this.timeFilterStart = ""
        this.timeFilterStop = ""
        this.filterYears()
    },
    aktionInTime(aktion){
        let aktionTimeStrength = aktion.Start.Year * 10000 + aktion.Start.Month * 100 + aktion.Start.Day
        let startTimeStrenght = 0
        let stopTimeStrenght = 999999999
        if(this.timeFilterStart != ""){
            startTimeStrenght = Number(this.timeFilterStart.slice(0,4)) * 10000 + Number(this.timeFilterStart.slice(5,7)) * 100 + Number(this.timeFilterStart.slice(8,10))
        }
        if(this.timeFilterStop != ""){
            stopTimeStrenght = Number(this.timeFilterStop.slice(0,4)) * 10000 + Number(this.timeFilterStop.slice(5,7)) * 100 + Number(this.timeFilterStop.slice(8,10))
        }
        return startTimeStrenght <= aktionTimeStrength && stopTimeStrenght >= aktionTimeStrength;
    },
    filterYears(){
        this.years = []
        this.actions.forEach(action => {
            let yearWritten = false
            if(!action.Title.includes(this.searchAktion)){
                yearWritten = true
            }
            if(!this.aktionInTime(action)){
                yearWritten = true
            }
            if(!yearWritten){
                this.years.forEach(year => {
                    if(year.year == action.Start.Year){
                        yearWritten = true;
                    }
                })
            }
            if(!yearWritten){
                this.years.push({year: action.Start.Year})
            }
        })
    },
    filteredActions(year){
        return this.actions.filter(action => action.Start.Year == year && action.Title.includes(this.searchAktion) && this.aktionInTime(action))
    }
  },
  mounted() {
    this.fetchActions()
    this.fetchAllTags()
    this.fetchAllPersons()
  }
}
</script>

<template>
    <h1>Browse Actions</h1>
    <div id="editAction" v-if="edit">
        <input id="titleInput" type="text" v-model="actionToedit.Title">
        <p>InfoText</p>
        <textarea v-model="actionToedit.InfoText"></textarea>
        <p>Start</p>
        <input type="date" v-model="editDateStart"/>
        <input type="time" v-model="editTimeStart">
        <p>Ende</p>
        <input type="date" v-model="editDateEnd"/>
        <input type="time" v-model="editTimeEnd">
        <div id="tags">
            <h3>Tags verwalten</h3>
            <h4>Aktuelle Tags:</h4>
            <div class="flexFelder PlatzNachOben">
                <div class="centeredLine" v-for="tag in actionToedit.Tags">{{ tag.Name }} <button @click="rmTagFromAction(tag.Name)" class="deleteButton opaccityAtHover"><img :src="trashIcon" alt="X"></button></div>
            </div>
            <h4>Auswählen von schon bekannten:</h4>
            <label for="searchTag">Suche Tag</label>
            <input @input="filterTags" type="text" id="searchTag" v-model="tagSearch">
            <div class="flexFelder PlatzNachOben">
            <button class="opaccityAtHover" v-for="tag in filteredTags" @click="addTagToAction(tag.Name)">{{ tag.Name }}</button>
            </div>
        </div>
        <div id="persons">
            <h3>Personen verwalten</h3> 
            <h4>Aktuelle Personen:</h4>
            <div class="flexFelder PlatzNachOben">
                <div class="centeredLine" v-for="Person in actionToedit.Visitors">{{ Person.Name }} <button @click="rmPersonFromAction(Person.Name)" class="deleteButton opaccityAtHover"><img :src="trashIcon" alt="X"></button></div>
            </div>
            <h4>Auswählen von schon bekannten:</h4>
            <label for="searchPerson">Suche Person</label>
            <input @input="filterPersons" type="text" id="searchPerson" v-model="personSearch">
            <div class="flexFelder PlatzNachOben">
            <button class="opaccityAtHover"  v-for="Person, i in filteredPersons"  @click="addPersonToAction(Person.Name)">{{ Person.Name }}</button>
            </div>
        </div>
        <br><br>
        <button @click="actualizeAction" class="button">Änderungen Übernehmen</button>
        <button @click="deleteAction" class="button PlatzNachLinks">Aktion löschen</button>
        <div class="imgContainer" v-if="actionToedit.pics">
            <div v-for="pic in actionToedit.pics">
                <a :href="'./aktionen' + actionToedit.Folder + pic.Path"><img :src="'./aktionen' + actionToedit.Folder + pic.Path" alt=""></a>
                <button @click="deletePicturefromAktion(actionToedit.Folder, pic.Path)"><p>&#9587;</p></button>
            </div>
        </div>
        <button @click="stopEditing" class="button PlatzNachOben">bearbeiten stoppen</button>
    </div>
    <div class="boxes">
        <div>
          <h2>Navigation</h2>
          <ul>
            <li><a href="./">Home</a></li>
            <li><a href="./settings">Settings</a></li>
            <li><a href="./diashow">Diashow ohne Infos</a></li>
            <li><a href="./diashow2">Diashow mit Infos</a></li>
            <li><a href="./upload">Foto Upload</a></li>
          </ul>
        </div>
        <div id="aktionFilter">
            <h2>Aktionen Filtern</h2>
            <label for="searchAktion">Suche nach Aktion</label>
            <input @input="filterYears" v-model="searchAktion" id="searchAktion" type="text">
            <br>
            <label for="startDate">Start Datum</label>
            <input @input="filterYears" id="startDate" v-model="timeFilterStart" type="date" name="">
            <br>
            <label for="endDate">End Datum</label>
            <input  @input="filterYears" id="endDate" v-model="timeFilterStop" type="date" name="">
            <br>
            <button @click="resetSearch" class="button">reset</button>
        </div>
        <template v-for="year in years">
            <h2>{{ year.year }}</h2>
            <div v-for="action in filteredActions(year.year)">
                <h2>{{ action.Title }}</h2>
                <p>{{ action.InfoText }}</p>
                <ul>
                    <li>Beginn: {{action.Start.Day}}.{{ action.Start.Month }}.{{ action.Start.Year }} {{ action.Start.Hour > 10 ? "" : "0" }}{{ action.Start.Hour }}:{{ action.Start.Minute > 10 ? "" : "0" }}{{ action.Start.Minute }}</li>
                    <li>Ende: {{action.End.Day}}.{{ action.End.Month }}.{{ action.End.Year }} {{ action.End.Hour > 10 ? "" : "0" }}{{ action.End.Hour }}:{{ action.End.Minute > 10 ? "" : "0" }}{{ action.End.Minute }}</li>
                    <li>Es waren da: 
                        <div class="flexFelder">
                            <div v-for="person in action.Visitors">{{ person.Name }}</div>
                        </div>
                    </li>
                    <li>Tags: 
                        <div class="flexFelder">
                            <div v-for="tag in action.Tags">{{ tag.Name }}</div>
                        </div>
                    </li>
                </ul>
                <button class="button" @click="editArticle(action)">Bearbeiten</button>
                <button v-if="!action.pics || !action.showPics" class="button PlatzNachLinks" @click="showAktionPics(action)">Bilder Laden</button>
                <button v-else class="button PlatzNachLinks" @click="hidePics(action.Folder)">Bilder Verbergen</button>
                <div class="imgContainer" v-if="action.pics && action.showPics">
                    <a v-for="pic in action.pics" :href="'./aktionen' + action.Folder + pic.Path"><img :src="'./aktionen' + action.Folder + pic.Path" alt=""></a>
                </div>
            </div>
        </template>
    </div>
</template>

<style>
#titleInput{
    font-size: 150%;
    max-width: 80%;
}
#editAction{
    z-index: 100;
    background-color: hsl(207, 22%, 18%);;
    position: fixed;
    width: calc(95vw - 2em);
    top: 2.5vh;
    left: 2.5vw;
    max-height: calc(95vh - 2em);
    padding: 1em;
    border-radius: 1.5em;
    box-shadow: 0 0 1em #ffffffaa;
    color: white;
    overflow-y: auto;
}
#editAction > input{
    /* color-scheme: light; */
    padding: 0.3em;
    border-radius: .5em;
}
#editAction > input[type="number"]{
    margin-right: 0;
    margin-bottom: 0.5em;
}
ul{
    padding-left: 1em;
}
#aktionFilter > input{
    margin-bottom: 0.5em;
}
</style>