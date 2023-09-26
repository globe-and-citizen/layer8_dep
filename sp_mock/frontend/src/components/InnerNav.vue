<script setup lang="ts">
import { RouterLink } from "vue-router";
import { ref } from "vue";

const poems = ref([
  {
    id: 1,
    name: "Poem 1",
    toName: "/poems/1",
  },
  {
    id: 2,
    name: "Poem 2",
    toName: "/poems/2",
  },
  {
    id: 3,
    name: "Poem 3",
    toName: "/poems/3",
  },
  {
    id: 4,
    name: "Poem 4",
    toName: "/poems/4",
  },
  {
    id: 5,
    name: "Poem 5",
    toName: "/poems/5",
  },
]);
</script>

<template>
  <nav class="flex justify-center mb-16">
    <div v-for="(i, index) in poems" :key="index">
      <RouterLink
        :to="i.toName"
        class="cursor-pointer px-4 py-2 hover:bg-green-300 hover:scale-105 hover:text-white select-none"
        >{{ i.name }}</RouterLink
      >
    </div>
  </nav>
</template>

<style scoped></style>
