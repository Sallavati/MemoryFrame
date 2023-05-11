<script>
export default {
  data() {
    return {
      tagSearch: "",
      CurrentAktion: null,
      allTags: null,
      filteredTags: null,
      newTagWritten: "",
      trashIcon: "./static/trash.png"
    }
  },
  methods: {
    addTag(name) {
      let formData = new URLSearchParams();
      console.log(name)
      formData.append('tag', name);
      if(name == this.newTagWritten){
        this.newTagWritten = ""
      }else{
        this.tagSearch = ""
      }
      fetch(
        `./api/addTag`,
        {
          method: 'post',
          body: formData
        }).then(data => {
          this.fetchCurrentAktion();
          this.fetchAllTags();
        });          
    },
    rmTag(name) {
      let formData = new URLSearchParams();
      console.log(name)
      formData.append('tag', name);
      fetch(
        `./api/rmTag`,
        {
          method: 'post',
          body: formData
        }).then(data => {
          this.fetchCurrentAktion();
          this.fetchAllTags();
        });          
    },
    filterTags(){
      this.filteredTags = this.allTags.filter((t) => {
        let exists = false
         if(!t.Name.includes(this.tagSearch)){
            return false
        }
        this.CurrentAktion.Tags.forEach(element => {
          if(t.Name == element.Name){
            exists = true
          }
        });
        return !exists;
      });
      this.filteredTags = this.filteredTags.slice(0, 14)
    },
    async fetchCurrentAktion() {
      this.CurrentAktion = null
      const res = await fetch(
        `./api/currentAction`
      )
      this.CurrentAktion = await res.json()
      this.filterTags()
    },
    async fetchAllTags() {
      this.allTags = null
      const res = await fetch(
        `./api/tags`
      )
      this.allTags = await res.json()
      this.filterTags()
    }
  },
  mounted() {
    this.fetchCurrentAktion()
    this.fetchAllTags()
  }
}
</script>

<template>
  <div id="tags" v-if="!CurrentAktion.aktionNotLive">
    <h3>Fügen Sie einen Tags hinzu</h3>
    <h4>Aktuelle Tags:</h4>
    <div class="flexFelder PlatzNachOben">
        <div class="centeredLine" v-for="tag in CurrentAktion.Tags">{{ tag.Name }} <button @click="rmTag(tag.Name)" class="deleteButton opaccityAtHover"><img :src="trashIcon" alt="X"></button></div>
    </div>
    <h4>Auswählen von schon bekannten:</h4>
    <label for="searchTag">Suche Tag</label>
    <input @input="filterTags" type="text" id="searchTag" v-model="tagSearch">
    <div class="flexFelder PlatzNachOben">
      <button class="opaccityAtHover" v-for="tag in filteredTags" @click="addTag(tag.Name)">{{ tag.Name }}</button>
    </div>
    <h4>Neuen eingeben:</h4>
    <input type="text" v-model="newTagWritten">   
    <button class="button" @click="addTag(newTagWritten)">add Tag</button>
  </div>
</template>

<style>
#tags{
  display: block;
  text-align: left;
  margin-top: 1em;
}
button{
  margin-left: 1em;
}
select{
  margin-right: 1em;
  color-scheme: dark;
}
</style>