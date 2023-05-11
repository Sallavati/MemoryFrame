<script setup>
import Tags from '../components/tags.vue';
import Persons from '../components/personen.vue';
import Infotext from '../components/infotext.vue';
</script>

<script>
export default {
    data() {
        return {
            CurrentAktion: [],
            Settings: [],
            hidden: "hidden",
            previewSrc: "./static/camera.png",
            imgUploaded: false
        };
    },
    methods: {
        showPreview(event) {
            this.previewSrc = URL.createObjectURL(event.target.files[0]);
            this.imgUploaded = true;
        },
        changeTitle() {
            let formData = new URLSearchParams();
            if (this.newTitle != "") {
                formData.append("name", this.newTitle);
                this.CurrentAktion.Title = this.newTitle;
                this.newTitle = "";
                fetch(`./api/changeAktionName`, {
                    method: "post",
                    body: formData
                });
            }
        },
        async fetchCurrentAktion() {
            this.CurrentAktion = null;
            const res = await fetch(`./api/currentAction`);
            this.CurrentAktion = await res.json();
        },
        async fetchSettings() {
            this.Settings = null;
            const res = await fetch(`./api/settings`);
            this.Settings = await res.json();
        },
    },
    mounted() {
        this.fetchCurrentAktion();
        this.fetchSettings();
    },
    components: { Infotext }
}
</script>

<template>
  <div id="everything">
    <div v-if="CurrentAktion.aktionNotLive">
      <h1>Aktuell läuft keine Aktion.</h1> 
    </div>
    <div v-else>
      <div v-if="!CurrentAktion.Title" id="setTitle">
        <h1>Gib der Aktion einen Namen.</h1>
        <input type="text" v-model="newTitle">   
        <button class="button" @click="changeTitle">Set Title</button>
      </div>
      <h1 v-else>{{ CurrentAktion.Title }}</h1>
    </div>
    <div class="boxes">
      <div>
          <h2>Navigation</h2>
          <ul>
            <li><a href="./">Home</a></li>
            <li><a href="./settings">Settings</a></li>
            <li><a href="./diashow">Diashow ohne Infos</a></li>
            <li><a href="./diashow2">Diashow mit Infos</a></li>
            <li><a href="./browse">Browse Aktionen</a></li>
          </ul>
        </div>
      <div>
        <form enctype="multipart/form-data"
          action="./api/upload"
          class="form-input"
          method="post">
          <img v-if="imgUploaded" id="previewPic" :src="previewSrc" alt="">
          <label v-if="!imgUploaded" for="imageFile">Lade ein Foto Hoch</label>
          <input  type="file" id="imageFile" name="imageFile" capture="user" accept="image/jpg,image/png" @change="showPreview"/>
          <br v-if="imgUploaded">
          <input class="button" v-if="imgUploaded" type="submit" value="upload">
        </form>
      </div>
      <Tags v-if="Settings.AllowTagEditing"/>
      <Persons v-if="Settings.AllowPersonEditing"/>
      <Infotext v-if="Settings.AllowInfoTextEditing"/>
    </div>
  </div>
</template>

<style>
#setTitle{
  margin-left: 1em;
  text-align: left;
  margin-bottom: 2em;
}
#setTitle > h2{
  text-align: left;
}
#setTitle > input, #setTitle > button{
  height:2em
}

#setTitle > button{
  margin-left: 1em;
}

#everything{
  display: flex;
  flex-direction: column;
  align-items: center;
}

#previewPic{
  width: 80%;
  aspect-ratio: 1 / 1;
  object-fit: cover;
  border-radius: 1em;
  margin-bottom: 1em;
}

.form-input{
  display: flex;
  flex-direction: column;
  align-items: center;
}

.form-input input[type="file"] {
  display:none;
}

.form-input label{
  display:block;
  font-size: 1.5em;
  /* line-height:50px; */
  text-align:center;
  /* background-color: hsl(221, 27%, 36%); */
  cursor:pointer;
  border: none
}

.form-input input[type="submit"]{
  font-size: 1.3em;
  width: 80%;
}
</style>