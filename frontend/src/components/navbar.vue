<script setup>
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
  <header>
      <nav>
        <RouterLink to="/">Home  </RouterLink>
        <RouterLink to="/diashow">Diashow  </RouterLink>
        <RouterLink to="/settings">settings  </RouterLink>
        <RouterLink to="/landing">Landing</RouterLink>
      </nav>
  </header>
</template>

<style scoped>
  
</style>
