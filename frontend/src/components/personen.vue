<script>
export default {
  data() {
    return {
      personSearch: "",
      CurrentAktion: null,
      allPersons: null,
      filteredPersons: null,
      newPersonSelected: "",
      newPersonWritten: "",
      trashIcon: "./static/trash.png"
    }
  },
  methods: {
    addPersonWithName(name = "") {
      let formData = new URLSearchParams();
      formData.append('person', name);
      if(name == this.newPersonWritten){
        this.newPersonWritten = ""
      }else{
        this.personSearch = ""
      }
      fetch(
        `./api/addPerson`,
        {
          method: 'post',
          body: formData
        }).then(data => {
          this.fetchCurrentAktion();
          this.fetchAllPersons();
        });
  },
  rmPerson(name = "") {
      let formData = new URLSearchParams();
      formData.append('person', name);
      fetch(
        `./api/rmPerson`,
        {
          method: 'post',
          body: formData
        }).then(data => {
          this.fetchCurrentAktion();
          this.fetchAllPersons();
        });
  },
    filterPersons(){
      this.filteredPersons = this.allPersons.filter((t) => {
        let exists = false
        if(!t.Name.includes(this.personSearch)){
            return false
        }
        this.CurrentAktion.Visitors.forEach(element => {
          if(t.Name == element.Name){
            exists = true
          }
        });
        return !exists;
      });
      this.filteredPersons = this.filteredPersons.slice(0, 14)
    },
    async fetchCurrentAktion() {
      this.CurrentAktion = null
      const res = await fetch(
        `./api/currentAction`
      )
      this.CurrentAktion = await res.json()
      this.filterPersons()
    },
    async fetchAllPersons() {
      this.CurrentAktion = null
      const res = await fetch(
        `./api/persons`
      )
      this.allPersons = await res.json()
      this.filterPersons()
    },
    logSearch(){
      console.log(this.personSearch)
      let testString = "TestDies Das"
      console.log(testString.includes(this.personSearch))
      //
    }
  },
  mounted() {
    this.fetchCurrentAktion()
    this.fetchAllPersons()
  }
}
</script>

<template>
  <div id="persons" v-if="!CurrentAktion.aktionNotLive">
    <h3>Fügen Sie Personen hinzu</h3> 
    <h4>Aktuelle Personen:</h4>
    <div class="flexFelder PlatzNachOben">
        <div class="centeredLine" v-for="Person in CurrentAktion.Visitors">{{ Person.Name }} <button @click="rmPerson(Person.Name)" class="deleteButton opaccityAtHover"><img :src="trashIcon" alt="X"></button></div>
    </div>
    <h4>Auswählen von schon bekannten:</h4>
    <label for="searchPerson">Suche Person</label>
    <input @input="filterPersons" type="text" id="searchPerson" v-model="personSearch">
    <div class="flexFelder PlatzNachOben">
      <button class="opaccityAtHover"  v-for="Person, i in filteredPersons"  @click="addPersonWithName(Person.Name)">{{ Person.Name }}</button>
    </div>
    <h4>Neuen eingeben:</h4>
    <input type="text" v-model="newPersonWritten">   
    <button class="button" @click="addPersonWithName(newPersonWritten)">add Person</button>
  </div>
</template>

<style>
#persons{
  text-align: left;
  margin-top: 1em;
}

button{
  margin-left: 1em;
}
select{
  color-scheme: dark;
  margin-right: 1em;
}
</style>