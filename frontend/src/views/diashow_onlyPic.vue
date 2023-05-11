<script>
export default {
  data() {
    return {
      diashow: null
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
  <div id="imgContainer">
    <img :src="diashow.Picture">
  </div>
</template>

<style>
img{
  height: 90%;
  width: 90%;
  object-fit: contain;
}
#imgContainer{
  display: flex;
  justify-content: space-around;
  align-items: center;
  margin: 0;
  height: 100vh;
  width: 100%;
  background-image: linear-gradient(30deg, #36d1dcbb, #5b86e5bb);
}
</style>