<script>
export default {
  data() {
    return {
      stats: null,
      seeMoreVisitors: false,
      seeMoreTags: false,
      speed: 2000,
      counts: [],
    }
  },
  methods: {
    async fetchStats() {
      this.stats = null
      const res = await fetch(
        `./api/stats`
      )
      this.stats = await res.json()
      this.counts = [
        {count: this.stats.ActionCount},
        {count: this.stats.VisitorCount},
        {count: this.stats.TagCount},
        {count: this.stats.ActionCount},
        {count: this.stats.ExternPictureCount},
        {count: this.stats.DataSize}
    ]
      this.countAllnumbers()
    },
    countAllnumbers(){
      this.stats.Tags.forEach((tag) => {
        tag.currentCount = 0
        tag.countHandler = setInterval(() => {
        tag.currentCount++;
        if(tag.currentCount >= tag.TimesUsed){
            clearInterval(tag.countHandler)
        }
        }, (this.speed/tag.TimesUsed))
      });
      this.stats.Visitors.forEach((visitor) => {
        visitor.currentCount = 0
        visitor.countHandler = setInterval(() => {
        visitor.currentCount++;
        if(visitor.currentCount >= visitor.TimesVisited){
            clearInterval(visitor.countHandler)
        }
        }, (this.speed/visitor.TimesVisited))
      });
      this.counts.forEach((counter) => {
        counter.currentCount = 0
        counter.countHandler = setInterval(() => {
        counter.currentCount++;
        if(counter.currentCount >= counter.count){
            clearInterval(counter.countHandler)
        }
        }, (this.speed/counter.count))
      });
    },
    recountTags(){
      this.stats.Tags.forEach((tag, index) => {
        if(index > 2){
          clearInterval(tag.countHandler)
          tag.currentCount = 0
          tag.countHandler = setInterval(() => {
          tag.currentCount++;
          if(tag.currentCount >= tag.TimesUsed){
              tag.currentCount = tag.TimesUsed
              clearInterval(tag.countHandler)
          }
          }, (this.speed/tag.TimesUsed))
        }
      });
    },
    recountVisitors(){
      this.stats.Visitors.forEach((visitor, index) => {
        if(index > 2){
          clearInterval(visitor.countHandler)
          visitor.currentCount = 0
          visitor.countHandler = setInterval(() => {
          visitor.currentCount++;
          if(visitor.currentCount >= visitor.TimesVisited){
              visitor.currentCount = visitor.TimesVisited
              clearInterval(visitor.countHandler)
          }
          }, (this.speed/visitor.TimesVisited))
        }
      });
    }
  },
  mounted() {
    this.fetchStats()
  }
}
</script>

<template>
  <div id="zahlen">
    <div>{{ counts[0].currentCount }} Aktionen</div>
    <div>{{ counts[1].currentCount }} Besucher</div>
    <div>{{ counts[2].currentCount }} Tags</div>
    <div>{{ counts[3].currentCount }} Fotos aus Aktionen</div>
    <div>{{ counts[4].currentCount }} Externe Fotos</div>
    <div>{{ counts[5].currentCount }} Mb Speicher verwendet</div>
  </div>
  <h3>Am meisten Besucht haben:</h3>
  <div id="visitors">
    <div v-for="person in stats.Visitors.slice(0,3)">
      <p>{{ person.Name }} ({{ person.currentCount }} mal)</p>
      <div :style="'height:' + 8*person.currentCount/stats.Visitors[0].TimesVisited + 'em'"></div>
    </div>
  </div>
  <button @click="seeMoreVisitors = !seeMoreVisitors; recountVisitors()" class="button PlatzNachOben">{{seeMoreVisitors ? "Weniger" : "Mehr"}} sehen</button>
  <div class="flexFelder PlatzNachOben" v-if="seeMoreVisitors">
    <div v-for="visitor in stats.Visitors.slice(3)">{{ visitor.Name }} ({{ visitor.currentCount }} mal)</div>
  </div>
  <br>
  <h3>Die meistbenutzten Tags sind:</h3>
  <div id="statTags" >
    <div v-for="tag in stats.Tags.slice(0,3)">
      <p>{{ tag.Name }} ({{ tag.currentCount }} mal)</p>
      <div :style="'height:' + 8*tag.currentCount/stats.Tags[0].TimesUsed + 'em'"></div>
    </div>
  </div>  
<button @click="seeMoreTags = !seeMoreTags; recountTags()" class="button PlatzNachOben">{{seeMoreTags ? "Weniger" : "Mehr"}} sehen</button>
<div class="flexFelder PlatzNachOben" v-if="seeMoreTags">
    <div v-for="tag in stats.Tags.slice(3)">{{ tag.Name }} ({{ tag.currentCount }} mal)</div>
</div>
</template>

<style>
h3{
  text-align: center;
}

#zahlen{
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr 1fr;
  gap: 1em;
  max-width: 80vw;
  margin: 0 auto;
  justify-content: space-between;
}
#zahlen > div{
  padding: 1em;
  background-image: linear-gradient(30deg, #36d1dcbb, #5b86e5bb);
  border-radius: .5em;
}

#visitors{
  display:flex;
  align-items: flex-end;
  justify-content: space-around;
  min-height: 12em;
}

#visitors > div > div{
  background-image: linear-gradient(30deg, #36d1dcbb, #5b86e5bb);
  border-top-right-radius: .5em;
  border-top-left-radius: .5em;
  margin: 0 .5em;
  width: 15vw;
}

#statTags > div, #visitors > div{
  display: flex;
  flex-direction: column;
  align-items: center;
}

#statTags{
  display:flex;
  align-items: flex-end;
  justify-content: space-around;
  min-height: 12em;
}

#statTags > div > div{
  background-image: linear-gradient(30deg, #36d1dcbb, #5b86e5bb);
  margin: 0 .5em;
  border-top-right-radius: .5em;
  border-top-left-radius: .5em;
  width: 15vw;
}
</style>