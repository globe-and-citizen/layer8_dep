<script setup>
import { ref } from "vue";

const poems = ref([
  {
    id: 1,
    name: "Poem 1",
  },
  {
    id: 2,
    name: "Poem 2",
  },
  {
    id: 3,
    name: "Poem 3",
  },
  {
    id: 4,
    name: "Poem 4",
  },
  {
    id: 5,
    name: "Poem 5",
  },
]);
</script>

<template>
  <div>
    <div class="flex justify-center mb-16">
      <div
        v-for="(i, index) in poems"
        :key="i"
        class="cursor-pointer px-4 py-2 hover:bg-green-300 hover:scale-105 hover:text-white select-none"
      >
        {{ i.name }}
      </div>
    </div>
    <div>
      <h1 class="font-extrabold text-2xl text-center mb-8">Poem title</h1>
      <div>
        <p class="font-normal text-md mx-16">
          Poem body Lorem ipsum dolor, sit amet consectetur adipisicing elit. Harum
          obcaecati aperiam ullam modi suscipit eos quam amet explicabo perspiciatis? Vel
          tempora doloremque esse dicta? Porro doloribus itaque culpa magnam omnis.
        </p>
      </div>
    </div>
  </div>
</template>
