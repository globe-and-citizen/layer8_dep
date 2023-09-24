<script setup>
import { RouterLink } from "vue-router";
defineProps({
  msg: {
    type: String,
    required: true,
  },
});
async function testWASMHandler() {
  const res = await window.testWASM(42, "42");
  console.log(res);
}
</script>

<template>
  <div class="greetings">
    <h1 class="green font-bold text-3xl mb-4">{{ msg }}</h1>
    <h3 class="font-medium text-base mb-4">
      This is an mock service provider that integrates with the Layer8 end-to-end
      encryption proxy.
      <button @click="testWASMHandler" class="text-green-500">Test WASM</button>
    </h3>
    <div class="flex space-x-4 font-medium text-base">
      <RouterLink to="/">Home</RouterLink>
      <RouterLink to="/login">Login</RouterLink>
    </div>
  </div>
</template>

<style scoped></style>
