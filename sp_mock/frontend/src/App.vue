<script setup>
import WelcomeMock from "./components/WelcomeMock.vue";
import { RouterView } from "vue-router";
</script>

<template>
  <div class="grid grid-cols-2 gap-x-2">
    <header>
      <div>
        <WelcomeMock msg="Welcome! We Got Poems Mock" />
      </div>
    </header>
    <RouterView />
  </div>
</template>

<style scoped>
header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
