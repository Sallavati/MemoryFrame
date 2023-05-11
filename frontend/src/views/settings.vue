<script setup>
import Tags from '../components/tags.vue';
import Infotext from '../components/infotext.vue';
import Persons from '../components/personen.vue';
import Navbar from '../components/navbar.vue';
import statsVue from '../components/stats.vue';
</script>

<script>
export default {
  data() {
    return {
      CurrentAktion: [],
      newTagSelected: "",
      newPersonSelected: "",
      personFilterSearch: "",
      tagFilterSearch: "",
      filteredTags: [],
      filteredPersons: [],
      allPersons: [],
      allTags: [],
      Settings: {},
      startTime: "",
      endTime: "",
      newTitle: "",
      qrLink: "./static/qr.png",
      switchOffLink: "./static/switch_off.png",
      switchOnLink: "./static/switch_on.png",
      trashIcon: "./static/trash.png",
      newAdditionalContent: {Title: "", InfoText: ""},
      ip: new URL("./", import.meta.url)
    }
  },
  methods: {
    changeTitle () {
      let formData = new URLSearchParams();
      if(this.newTitle != ""){
        formData.append('name', this.newTitle);
        this.CurrentAktion.Title = this.newTitle
        this.newTitle = ""
        fetch(
          `./api/changeAktionName`,
          {
            method: 'post',
            body: formData
          })
      }
    },
    filterTags(){
      this.filteredTags = this.allTags.filter((t) => {
        let exists = false
        this.Settings.Diashowfilter.Tags.forEach(element => {
          if(t.Name == element.Name){
            exists = true
          }
        });
        return !exists;
      });
      this.filteredTags = this.filteredTags.filter(tag => tag.Name.includes(this.tagFilterSearch))
    },
    filterPersons(){
      this.filteredPersons = this.allPersons.filter((t) => {
        let exists = false
        this.Settings.Diashowfilter.Persons.forEach(element => {
          if(t.Name == element.Name){
            exists = true
          }
        });
        return !exists;
      });
      this.filteredPersons = this.filteredPersons.filter(person => person.Name.includes(this.personFilterSearch)) 
    },
    async fetchCurrentAktion() {
      this.CurrentAktion = null
      const res = await fetch(
        `./api/currentAction`
      )
      this.CurrentAktion = await res.json()
    },
    async fetchAllPersons() {
      this.allPersons = null
      const res = await fetch(
        `./api/persons`
      )
      this.allPersons = await res.json()
      this.filterPersons()
    },
    async fetchAdditioalContentPics(index){
      let formData = new URLSearchParams();
      if(index < this.Settings.AdditionalPics.length){
        formData.append('folder', index);
        let result = await fetch(
          `./api/getExternContentPics`,
          {
            method: 'post',
            body: formData
          })
          this.Settings.AdditionalPics[index].pics = await result.json()
      }
    },
    async fetchAllTags() {
      this.allTags = null
      const res = await fetch(
        `./api/tags`
      )
      this.allTags = await res.json()
      this.filterTags()
    },
    async fetchSettings() {
      this.Settings = null
      const res = await fetch(
        `./api/settings`
      )
      this.Settings = await res.json()
      this.startTime = `${this.Settings.Diashowfilter.TimeStart.Year}-${this.Settings.Diashowfilter.TimeStart.Month > 9 ? this.Settings.Diashowfilter.TimeStart.Month : "0" + this.Settings.Diashowfilter.TimeStart.Month}-${this.Settings.Diashowfilter.TimeStart.Day > 9 ? this.Settings.Diashowfilter.TimeStart.Day : "0" + this.Settings.Diashowfilter.TimeStart.Day}`
      this.endTime = `${this.Settings.Diashowfilter.TimeEnd.Year}-${this.Settings.Diashowfilter.TimeEnd.Month > 9 ? this.Settings.Diashowfilter.TimeEnd.Month : "0" + this.Settings.Diashowfilter.TimeEnd.Month}-${this.Settings.Diashowfilter.TimeEnd.Day > 9 ? this.Settings.Diashowfilter.TimeEnd.Day : "0" + this.Settings.Diashowfilter.TimeEnd.Day}`
    }, 
    actualizeContentText(text, id) {
      let formData = new URLSearchParams();
      formData.append('text', text);
      formData.append('index', id);
      fetch(
        `./api/editExternContentInfoText`,
        {
          method: 'post',
          body: formData
        }).then(data => {
          this.fetchSettings();
        });          
    },
    async fetchCurrentAktion() {
      this.CurrentAktion = null
      const res = await fetch(
        `./api/currentAction`
      )
      this.CurrentAktion = await res.json()
    },
    toggleSortAfterTags(){
      fetch(`./api/toggleSortAfterTags`)
      this.Settings.Diashowfilter.SortAfterTags = !this.Settings.Diashowfilter.SortAfterTags
    },
    toggleSortAfterPersons(){
      fetch(`./api/toggleSortAfterPersons`)
      this.Settings.Diashowfilter.SortAfterPersons = !this.Settings.Diashowfilter.SortAfterPersons
    },
    toggleSortAfterTimeStart(){
      fetch(`./api/toggleSortAfterTimeStart`)
      this.Settings.Diashowfilter.SortAfterTimeStart = !this.Settings.Diashowfilter.SortAfterTimeStart
    },
    toggleSortAfterTimeEnd(){
      fetch(`./api/toggleSortAfterTimeEnd`)
      this.Settings.Diashowfilter.SortAfterTimeEnd = !this.Settings.Diashowfilter.SortAfterTimeEnd
    },
    toggleOnlyAdditionalContent(){
      fetch(`./api/toggleOnlyAdditionalContent`)
      this.Settings.Diashowfilter.OnlyExternContent = !this.Settings.Diashowfilter.OnlyExternContent
    },
    toggleAllowPersonEditing(){
      fetch(`./api/toggleAllowPersonEditing`)
      this.Settings.AllowPersonEditing = !this.Settings.AllowPersonEditing
    },
    toggleAllowInfoTextEditing(){
      fetch(`./api/toggleAllowInfoTextEditing`)
      this.Settings.AllowInfoTextEditing = !this.Settings.AllowInfoTextEditing
    },
    toggleAllowTagEditing(){
      fetch(`./api/toggleAllowTagEditing`)
      this.Settings.AllowTagEditing = !this.Settings.AllowTagEditing
    },
    toggleFreeUpload(index){
      let formData = new URLSearchParams();
      formData.append('index', String(index));
      fetch(
        `./api/toggleFreeUpload`,
        {
          method: 'post',
          body: formData
        })
        this.Settings.AdditionalPics[index].FreeUpload = !this.Settings.AdditionalPics[index].FreeUpload
    },
    toggleShowExternContent(index){
      let formData = new URLSearchParams();
      formData.append('index', String(index));
      this.Settings.AdditionalPics[index].Show = !this.Settings.AdditionalPics[index].Show
      fetch(
        `./api/toggleShowExternContent`,
        {
          method: 'post',
          body: formData
        })
    },
    ignorePerson(name){
      let formData = new URLSearchParams();
      formData.append('name', name);
      fetch(
        `./api/ignorePerson`,
        {
          method: 'post',
          body: formData
        })
        this.Settings.Diashowfilter.Persons = this.Settings.Diashowfilter.Persons.filter((t) => {
          let samePerson = false
          if(t.Name == name){
            samePerson = true
          }else{
            return !samePerson
          }
        })
        this.filterPersons()
    },
    ignoreTag(name){
      let formData = new URLSearchParams();
      formData.append('name', name);
      fetch(
        `./api/ignoreTag`,
        {
          method: 'post',
          body: formData
        })
        this.Settings.Diashowfilter.Tags = this.Settings.Diashowfilter.Tags.filter((t) => {
          let sameTag = false
          if(t.Name == name){
            sameTag = true
          }else{
            return !sameTag
          }
        })
        this.filterTags()
    },
    stopAktion(){
      fetch(`./api/stopAction`)
      this.CurrentAktion.aktionNotLive = true;
    },
    setTag(name){
      let formData = new URLSearchParams();
      if(name != ""){
        formData.append('name', name);
        fetch(
          `./api/setTag`,
          {
            method: 'post',
            body: formData
          })
          this.Settings.Diashowfilter.Tags.push({Name: name})
          this.newTagSelected = ""
          this.filterTags()
      }
    },
    setPerson(name){
      let formData = new URLSearchParams();
        formData.append('name', name);
        fetch(
        `./api/setPerson`,
        {
          method: 'post',
          body: formData
        })
        this.Settings.Diashowfilter.Persons.push({Name: name})
        this.filterPersons()
    },
    setFilterTimeEnd(){
      if(this.endTime != ""){
      let formData = new URLSearchParams();
      formData.append('Year', this.endTime.slice(0, 4));
      formData.append('Month', this.endTime.slice(5, 7));
      formData.append('Day', this.endTime.slice(8, 10));
      fetch(
        `./api/setFilterTimeEnd`,
        {
          method: 'post',
          body: formData
        })
      }
    },
    setFilterTimeStart(){
      let formData = new URLSearchParams();
      if(this.startTime != ""){
      formData.append('Year', this.startTime.slice(0, 4));
      formData.append('Month', this.startTime.slice(5, 7));
      formData.append('Day', this.startTime.slice(8, 10));
      fetch(
        `./api/setFilterTimeStart`,
        {
          method: 'post',
          body: formData
        })
      }
    },
    setDiashowSpeed(){
      let formData = new URLSearchParams();
      formData.append('DiashowSpeed', this.Settings.DiashowSpeed);
      fetch(
        `./api/setDiashowSpeed`,
        {
          method: 'post',
          body: formData
        })
    },
    setAutoActionStop(){
      let formData = new URLSearchParams();
      formData.append('AutoActionStop', this.Settings.AutoActionStop);
      fetch(
        `./api/setAutoActionStop`,
        {
          method: 'post',
          body: formData
      })
    },
    addAdditionalContent(){
      let nameUsed = false
      this.Settings.AdditionalPics.forEach(folder => {
        if(folder.Name.replaceAll(" ", "") == this.newAdditionalContent.Title.replaceAll(" ", "")){
          nameUsed = true
          alert("Name ist bereits vergeben")
        }
      })
      if(nameUsed){
        return
      }
      let formData = new URLSearchParams();
      if(this.newAdditionalContent.Title != ""){
        formData.append('Title', this.newAdditionalContent.Title);
        formData.append('InfoText', this.newAdditionalContent.InfoText);
        fetch(
          `./api/addAdditionalContent`,
          {
            method: 'post',
            body: formData
        })
        this.Settings.AdditionalPics.push({
          "Name": this.newAdditionalContent.Title,
          "InfoText": this.newAdditionalContent.InfoText,
          "Path": "Reload Page to see",
          'Show': false
        })
        this.newAdditionalContent.InfoText = ""
        this.newAdditionalContent.Title = ""
      }
    },
    async deleteAdditionalContent(index){
      if(window.confirm("Sicher, dass Sie den Ordner löschen wollen? Alle Bilder werden verloren gehen")){
        let formData = new URLSearchParams();
        formData.append('folder', index);
        fetch(
          `./api/deleteAdditionalContent`,
          {
            method: 'post',
            body: formData
        })
      }
      this.Settings.AdditionalPics = this.Settings.AdditionalPics.filter(folder => folder.Name != this.Settings.AdditionalPics[index].Name)      
    },
    async deletePicturefromExternContent(index, name){
      if(window.confirm("Sicher, dass Sie das Bild löschen wollen?")){
        let formData = new URLSearchParams();
        formData.append('folder', index);
        formData.append('picName', name);
        fetch(
          `./api/deletePictureFromAdditionalContent`,
          {
            method: 'post',
            body: formData
        })
        console.log(this.Settings.AdditionalPics[index].pics)
        this.Settings.AdditionalPics[index].pics = this.Settings.AdditionalPics[index].pics.filter(pic => {return pic.Path != name})
      }
    },
  },
  mounted() {
    this.fetchCurrentAktion()
    this.fetchSettings()
    this.fetchAllPersons()
    this.fetchAllTags()
  }
}
</script>

<template>
    <h1>Einstellungen</h1>
    <div class="boxes">
    <div>
      <h2>Navigation</h2>
      <ul>
        <li><a href="./">Home</a></li>
        <li><a href="./upload">Foto Upload</a></li>
        <li><a href="./diashow">Diashow ohne Infos</a></li>
        <li><a href="./diashow2">Diashow mit Infos</a></li>
        <li><a href="./browse">Browse Aktionen</a></li>
      </ul>
    </div>
    <div v-if="!CurrentAktion.aktionNotLive">
      <h2>Aktuelle Aktion:</h2>
      <Persons/>
      <Tags/>
      <Infotext/>
      <button class="button" @click="stopAktion">Stop Aktion</button>
    </div>
    <div>
      <h2>Generelle Einstellungen</h2>
      <div>
        <label for="DiashowSpeed">Diashow Geschwindigkeit</label>
        <input @input="setDiashowSpeed" type="number" name="DiashowSpeed" v-model="Settings.DiashowSpeed">
      </div>
      <div class="space">
        <label for="AutoActionStop">Auto Aktions Stop in min</label>
        <input @input="setAutoActionStop" type="number" name="AutoActionStop" v-model="Settings.AutoActionStop">
      </div>
      <p>
        Tag Bearbeitung erlauben?
        <button class="switchIconButton" @click="toggleAllowTagEditing"><img :src="Settings.AllowTagEditing ? switchOnLink : switchOffLink"/></button>
      </p>
      <p>
        Besucher Bearbeitung erlauben?
        <button class="switchIconButton" @click="toggleAllowPersonEditing"><img :src="Settings.AllowPersonEditing ? switchOnLink : switchOffLink"/></button>
      </p>
      <p>
        InfoText Bearbeitung erlauben?
        <button class="switchIconButton" @click="toggleAllowInfoTextEditing"><img :src="Settings.AllowInfoTextEditing ? switchOnLink : switchOffLink"/></button>
      </p>
    </div>

    <div>
      <h2>Diashow Filter</h2>
      <h3>
        Nur Externen Content anzeigen
        <button class="switchIconButton" @click="toggleOnlyAdditionalContent"><img :src="Settings.Diashowfilter.OnlyExternContent ? switchOnLink : switchOffLink"/></button>
      </h3>
      <h3>
        Nach Tags sortieren
        <button class="switchIconButton" @click="toggleSortAfterTags"><img :src="Settings.Diashowfilter.SortAfterTags ? switchOnLink : switchOffLink"/></button>
      </h3>
      <div v-if="Settings.Diashowfilter.SortAfterTags">
        <div class="flexFelder PlatzNachOben">
          <div class="centeredLine" v-for="tag in Settings.Diashowfilter.Tags">{{ tag.Name }} <button @click="ignoreTag(tag.Name)" class="deleteButton opaccityAtHover"><img :src="trashIcon" alt="X"></button></div>
        </div>
        <h4>Tag zum Filtern Auswählen:</h4>
        <label for="searchTag">Suche Tag</label>
        <input @input="filterTags" type="text" id="searchTag" v-model="tagFilterSearch">
        <div class="flexFelder PlatzNachOben">
          <button class="opaccityAtHover"  v-for="tag in filteredTags"  @click="setTag(tag.Name)">{{ tag.Name }}</button>
        </div>
      </div>

      <h3>
        Nach Personen sortieren
        <button class="switchIconButton" @click="toggleSortAfterPersons"><img :src="Settings.Diashowfilter.SortAfterPersons ? switchOnLink : switchOffLink"/></button>
      </h3>
      <div v-if="Settings.Diashowfilter.SortAfterPersons">
        <div class="flexFelder PlatzNachOben">
          <div class="centeredLine" v-for="Person in Settings.Diashowfilter.Persons">{{ Person.Name }} <button @click="ignorePerson(Person.Name)" class="deleteButton opaccityAtHover"><img :src="trashIcon" alt="X"></button></div>
        </div>
        <h4>Person zum Filtern Auswählen:</h4>
        <label for="searchPerson">Suche Person</label>
        <input @input="filterPersons" type="text" id="searchPerson" v-model="personFilterSearch">
        <div class="flexFelder PlatzNachOben">
          <button class="opaccityAtHover"  v-for="person in filteredPersons"  @click="setPerson(person.Name)">{{ person.Name }}</button>
        </div>
      </div>
      
      <h3>Nach Zeitraum sortieren:</h3>
      <h4>
        Start
        <button class="switchIconButton" @click="toggleSortAfterTimeStart"><img :src="Settings.Diashowfilter.SortAfterTimeStart ? switchOnLink : switchOffLink"/></button>
      </h4>
      <input v-model="startTime" type="date" @input="setFilterTimeStart">
      <h4>
        Ende
        <button class="switchIconButton" @click="toggleSortAfterTimeEnd"><img :src="Settings.Diashowfilter.SortAfterTimeEnd ? switchOnLink : switchOffLink"/></button>
      </h4>
      <input v-model="endTime" type="date" @input="setFilterTimeEnd">
    </div>

    <div>
      <h2>Externer Content</h2>
      <div v-for="(folder, index) in Settings.AdditionalPics">
        <h3>
          {{ folder.Name }}
          <button class="switchIconButton" @click="toggleShowExternContent(index)"><img :src="Settings.AdditionalPics[index].Show ? switchOnLink : switchOffLink"/></button>
          <button @click="folder.show = !folder.show" class="button">{{ folder.show ? "Hide" : "Edit"}}</button>
        </h3>
        <div v-if="folder.show">
          <p>Speicherort: {{ folder.Path }}</p>
          <h4>Upload Picture into Folder</h4>
          <p>
            Foto Verwaltungs Link Freigeschalten
            <button class="switchIconButton" @click="toggleFreeUpload(index)"><img :src="folder.FreeUpload ? switchOnLink : switchOffLink"/></button>
          </p>
          <p><a v-if="folder.FreeUpload" :href="'./uploadToExtern' +  folder.Path">{{ ip + 'uploadToExtern' +  folder.Path }}</a></p>
          <form action="./api/uploadToAdditionalContent" enctype="multipart/form-data" method="post">
            <input id="picUpoad" class="PlatzNachUnten" type="file" name="imageFile" accept="image/jpg,image/png">
            <input type="hidden" name="folder" :value="index">
            <input class="button" type="submit" value="upload">
          </form>
          <h4>Infotext:</h4> 
          <textarea v-model="folder.InfoText" name="" id="" cols="30" rows="10"></textarea>
          <button @click="actualizeContentText(folder.InfoText, index)" class="button PlatzNachUnten">Submit Text Changes</button>
          <br>
          <button @click="folder.showPics = !folder.showPics; fetchAdditioalContentPics(index)" class="button">Bilder {{ !folder.showPics ? "verwalten" : "verstecken"}}</button>
          <div class="imgContainer" v-if="folder.showPics" >
            <div v-for="pic in folder.pics">
              <a :href="'./externContent' + folder.Path + '/' + pic.Path"><img id="actionpic" :src="'./externContent' + folder.Path + '/' + pic.Path" alt=""></a>    
              <button @click="deletePicturefromExternContent(index, pic.Path)"><p>&#9587;</p></button>
            </div>
          </div>
          <br>
          <button class="button PlatzNachOben" @click="deleteAdditionalContent(index)">Delete Folder</button>
          <br>
          <br>
        </div>
      </div>
      <br>
      <h3>Neuen Ordner für Zusätzlichen Content einrichten</h3>
      <div class="space">
        <label for="Title">Titel</label>
        <input type="text" v-model="newAdditionalContent.Title" name="Title">
      </div>
      <div class="space">
        <label for="InfoText">InfoText</label>
        <input type="text" v-model="newAdditionalContent.InfoText" name="Infotext">
      </div>
      <button class="space button" @click="addAdditionalContent">Add</button>
    </div>
    <div>
      <h2>Statistiken</h2>
      <statsVue/>
    </div>

    <div>
      <h2>Qr-Code zum Fotohochladen</h2>
      <img id="qrCode" :src="qrLink" alt="qr-code">
    </div>  
  </div>
</template>

<style>
label{
  margin-right: 1em;
}
.space{
  margin-top: .5em;
}
button{
  margin: 0;
}
input{
  margin-right: 1em;
}
#qrCode{
  max-width: 300px;
  aspect-ratio: 1 / 1;
  border-radius: 1em;
  margin-bottom: 2em;
  margin-left: 1em;
}
.switchIconButton{
  background-color: #00000000;
  border: none;
  height: 1.5em;
  object-fit: contain;
}
h3, h4{
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: .5em; 
  text-align: left;
  margin-bottom: .5em;
  margin-left: 0;
}
.DateInput{
  display: flex;
  flex-direction: column;
}
.DateInput > *{
  margin-bottom: 0.5em;
}
.DateInput > div > input{
  margin-left: .5em;
}
.DateInput > button{
  max-width: 10ch;
}
input[type="number"]{
  max-width: 10ch;
}
</style>



