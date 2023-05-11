<script>
export default {
  data() {
    return {
      newTitle: "",
    }
  },
  methods: {
    changeTitle () {
      let formData = new URLSearchParams();
      if(this.newTitle != ""){
        formData.append('name', this.newTitle);
        this.newTitle = ""
        fetch(
          `http://localhost:3000/api/changeAktionName`,
          {
            method: 'post',
            body: formData
          })//.then(location.reload())
      }
    }
  }
}
</script>

<template>
    Was für einen Titel soll die Aktion bekommen?
    <input type="text" v-model="newTitle">   
    <button @click="changeTitle">Set Title</button>
</template>