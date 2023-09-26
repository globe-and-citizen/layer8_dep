<script setup>
import InnerNav from "../components/InnerNav.vue";
import { useRoute } from "vue-router";
import { watch, onMounted, ref } from "vue";

const route = useRoute();

var poem = ref({
  title: "",
  body: "",
});

const getPoem = async () => {
  const response = await fetch(`http://localhost:3000/api/poems/${route.params.id}`);
  poem.value = await response.json();
  // replace newlines with <br>
  poem.value.body = poem.value.body.replace(/\n/g, "<br>");
};

// listen for route changes
watch(route, () => {
  getPoem();
});

onMounted(() => {
  getPoem();
});
</script>

<template>
  <main class="min-h-(screen-20) overflow-y-auto">
    <InnerNav />
    <h1 class="mb-4 text-2xl font-bold">{{ poem.title }}</h1>
    <p v-html="poem.body" class="mb-4"></p>
  </main>
</template>

<style scoped>
/* You don't need the style block for Tailwind CSS. */
</style>
