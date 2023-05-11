<script>
export default {
  data() {
    return {
      CurrentAktion: {InfoText: ""},
    }
  },
  methods: {
    actualizeText() {
      let formData = new URLSearchParams();
      formData.append('text', this.CurrentAktion.InfoText);
      fetch(
        `./api/editAktionInfotext`,
        {
          method: 'post',
          body: formData
        }).then(data => {
          this.fetchCurrentAktion();
        });          
    },
    async fetchCurrentAktion() {
      this.CurrentAktion = null
      const res = await fetch(
        `./api/currentAction`
      )
      this.CurrentAktion = await res.json()
    },
  },
  mounted() {
    this.fetchCurrentAktion()
    this.fetchAllTags()
  }
}
</script>

<template>
  <div v-if="!CurrentAktion.aktionNotLive">
    <h3>Passen Sie den InfoText an:</h3> 
    <textarea v-model="CurrentAktion.InfoText" name="" id="" cols="30" rows="10"></textarea>
    <button @click="actualizeText" class="button PlatzNachUnten">Submit Text Changes</button>
  </div>
</template>

<style>
  
</style>