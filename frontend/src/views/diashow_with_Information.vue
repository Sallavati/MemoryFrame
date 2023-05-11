<script>
export default {
  data() {
    return {
      diashow: null,
      qrLink: "./static/qr.png"
    }
  },
  methods: {
    async fetchDiashow() {
      this.diashow = null
      const res = await fetch(
        `./api/diashow`
      )
      this.diashow = await res.json()
      setTimeout(this.fetchDiashow, this.diashow.DiashowSpeed*1000)
      console.log("timer set")
    }
  },
  mounted() {
    this.fetchDiashow()
  }
}
</script>

<template>
  <div id="stuff">
    <img id="diashowpic" :src="diashow.Picture">
    <div class="infoBox" v-if="diashow.FromAktion">
        <h2>{{ diashow.Aktion.Title }}</h2>
        <p>Von {{ diashow.Aktion.Start.Day }}.{{ diashow.Aktion.Start.Month }}.{{ diashow.Aktion.Start.Year }} {{ diashow.Aktion.Start.Hour > 10 ? "" : "0" }}{{ diashow.Aktion.Start.Hour }}:{{ diashow.Aktion.Start.Minute > 10 ? "" : "0" }}{{ diashow.Aktion.Start.Minute}}</p>
        <p>Bis {{ diashow.Aktion.End.Day }}.{{ diashow.Aktion.End.Month }}.{{ diashow.Aktion.End.Year }} {{ diashow.Aktion.End.Hour > 10 ? "" : "0"  }}{{ diashow.Aktion.End.Hour }}:{{ diashow.Aktion.End.Minute > 10 ? "" : "0"  }}{{ diashow.Aktion.End.Minute }}</p>
        <h3 v-if="diashow.Aktion.Visitors.length > 0" >Besucher:</h3>
        <div v-if="diashow.Aktion.Visitors.length > 0" class="flexFelder PlatzNachOben">
            <div class="allblack" v-for="visitor in diashow.Aktion.Visitors">{{ visitor.Name }}</div>
        </div>
        <h3 v-if="diashow.Aktion.Tags.length > 0">Tags:</h3>
        <div v-if="diashow.Aktion.Tags.length > 0" class="flexFelder PlatzNachOben">
            <div class="allblack" v-for="tag in diashow.Aktion.Tags">{{ tag.Name }}</div>
        </div>
        <h3 v-if="diashow.Aktion.InfoText.length > 0">InfoText:</h3>
        <p v-if="diashow.Aktion.InfoText.length > 0">{{ diashow.Aktion.InfoText }}</p>
        <h3>Lade ein Bild Hoch:</h3>
        <img class="qrCode" :src="qrLink" alt="qr-code">
    </div>
    <div class="infoBox" v-else>
      <h2>{{ diashow.Folder.Name }}</h2>
      <p>{{ diashow.Folder.InfoText }}</p>
      <h3>Lade ein Bild Hoch:</h3>
      <img class="qrCode" :src="qrLink" alt="qr-code">
    </div>
  </div>
</template>

<style>
#diashowpic{
  display: block;
  max-width: 65vw;
  max-height: 90vh;
  object-fit: contain;
}
#stuff{
  display: flex;
  justify-content: space-around;
  align-items: center;
  margin: 0;
  padding: 0;
  height: 100vh;
  width: 100%;
  overflow: hidden;
  background-image: linear-gradient(30deg, #36d1dcbb, #5b86e5bb);
}
.infoBox{
  display:block;
  margin-right: 3em;
  border: solid 2px black;
  padding: 2.5vh 2vw;
  width: 21vw;
  max-height: 95vh;
  overflow-x:hidden;
}
.qrCode{
  max-width: 15vw;
  aspect-ratio: 1 / 1;
  margin-left: 3vw;
  margin-top: 1em;
}
</style>