<script>
export default {
    data() {
        return {
            Settings: [],
            hidden: "hidden",
            previewSrc: "./static/camera.png",
            imgUploaded: false,
            currentFolder: {},
            currentFolderIndex: null
        };
    },
    methods: {
        showPreview(event) {
            this.previewSrc = URL.createObjectURL(event.target.files[0]);
            this.imgUploaded = true;
        },
        async fetchSettings() {
            this.Settings = null;
            const res = await fetch(
                `../api/settings`
                );
            this.Settings = await res.json();
            console.log(this.Settings)
            this.getFolder()
        },
        getFolder(){
            this.Settings.AdditionalPics.forEach((folder, index) => {
                //console.log("/" + this.$route.params.folder)
                //console.log(folder.Path)
                if("/" + this.$route.params.folder == folder.Path){
                    this.currentFolder = folder
                    this.currentFolderIndex = index
                    this.fetchAdditioalContentPics()
                }
            })
        },
        async deletePicturefromExternContent(name){
      if(window.confirm("Sicher, dass Sie das Bild löschen wollen?") && this.currentFolderIndex != null){
        let formData = new URLSearchParams();
        formData.append('folder', this.currentFolderIndex);
        formData.append('picName', name);
        fetch(
          `../api/deletePictureFromAdditionalContent`,
          {
            method: 'post',
            body: formData
        })
        this.currentFolder.pics = this.currentFolder.pics.filter(pic => {return pic.Path != name})
      }
    },
    async fetchAdditioalContentPics(){
        if(this.currentFolderIndex != null){
            let formData = new URLSearchParams();
            formData.append('folder', this.currentFolderIndex);
            let result = await fetch(
            `../api/getExternContentPics`,
            {
                method: 'post',
                body: formData
            })
            this.currentFolder.pics = await result.json()
        }
    }
    },
    mounted() {
        this.fetchSettings();
    }
}
</script>

<template>
    <div v-if="currentFolder.Path && currentFolder.FreeUpload" id="everything">
    <h1>Upload to {{ currentFolder.Name }}</h1>
    <div class="boxes">
      <div>
          <h2>Navigation</h2>
          <ul>
            <li><a href="../">Home</a></li>
            <li><a href="../upload">Aktions Foto upload</a></li>
            <li><a href="../settings">Settings</a></li>
            <li><a href="../diashow">Diashow ohne Infos</a></li>
            <li><a href="../diashow2">Diashow mit Infos</a></li>
            <li><a href="../browse">Browse Aktionen</a></li>        
        </ul>
        </div>
      <div>
        <form enctype="multipart/form-data"
          action="../api/uploadToAdditionalContent"
          class="form-input"
          method="post">
          <img v-if="imgUploaded" id="previewPic" :src="previewSrc" alt="">
          <label v-if="!imgUploaded" for="imageFile">Lade ein Foto Hoch</label>
          <input  type="file" id="imageFile" name="imageFile" accept="image/jpg,image/png" @change="showPreview"/>
          <input type="hidden" name="folder" :value="currentFolderIndex">
          <input type="hidden" name="fromUploadPage" value="true">
          <br v-if="imgUploaded">
          <input class="button" v-if="imgUploaded" type="submit" value="upload">
        </form>
      </div>
      <div>
        <h2>Fotos verwalten:</h2>
        <div class="imgContainer" v-if="currentFolder.pics" >
            <div v-for="pic in currentFolder.pics">
              <a :href="'../externContent' + currentFolder.Path + '/' + pic.Path"><img id="actionpic" :src="'../externContent' + currentFolder.Path + '/' + pic.Path" alt=""></a>    
              <button @click="deletePicturefromExternContent(pic.Path)"><p>&#9587;</p></button>
            </div>
          </div>
      </div>
    </div>
  </div>
  <div v-else id="everything">
    <h1>Extern Content upload</h1>
    <div class="boxes">
        <div>
            <p>Leider konnte der angeforderte Ordner nicht gefunden werden oder ist nicht für jeden zum Upload freigeschalten.</p>
        </div>
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