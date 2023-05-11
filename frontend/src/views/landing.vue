<script setup>
    import statsVue from '../components/stats.vue';
</script>

<template>
    <div id="content">
        <h1>Wilkommen beim MemoryFrame</h1>
        <div class="boxes">
            <div>
                <h2>About</h2>
                <p>
                    Bei MemoryFrame handelt es sich um ein Web basiertes 
                    Bilderrahmen-System, um bei Feiern, einem entspannten Abend 
                    oder jeglicher anderen Form von Aktion nichts zu verpassen 
                    und eine coole Diashow zu bieten. Dafür können bei Aktionen live 
                    Fotos hochgeladen werden, als auch externe Bilder in extra Ordnern 
                    gespeuichert werden. In den Zwei Diashow Ansichten, kann man 
                    sich diese dann nach belieben gefiltert über weitere Geräte 
                    entspannt anschauen. Also warte nicht lang und lade dein
                    erstes Foto hoch.
                </p>
            </div>
            <div>
                <h2>Navigation</h2>
                <ul>
                    <li><a href="./upload">Foto Upload</a></li>
                    <li><a href="./diashow">Diashow ohne Infos</a></li>
                    <li><a href="./diashow2">Diashow mit Infos</a></li>
                    <li><a href="./settings">Settings</a></li>
                    <li><a href="./browse">Browse Aktionen</a></li>
                </ul>
            </div>
            <div>
                <h2>Statistiken</h2>
                <statsVue/>
            </div>
        </div>
    </div>
</template>

<style>
#content{
    max-width: 90vw;
}
 #diashow-views{
    display: flex;
    flex-direction: column;
    max-width: 80vw;
 }
 #diashow-views > div {
    height: 80vw;
    border-radius: 2em;
 }
</style>