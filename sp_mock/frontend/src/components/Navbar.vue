<script setup>
import {ref} from "vue"

const counter = ref(0)

async function testWASMLoadedHandler (){
  let res = await layer8.testWASMLoaded()
  console.log("Wasm Loaded")
  counter.value = res
}

</script>

<template>
<div id="navbar">
    <RouterLink to="/">Home</RouterLink>
    <RouterLink to="/stress-test">Stress test</RouterLink>
    <button @click="testWASMLoadedHandler">TestWASM</button>
    <span v-if="counter != 0 " >{{ counter }}</span >
    <br><hr><br>
</div>
</template>

<style>
</style>