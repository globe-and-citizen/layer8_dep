<script setup>
import {ref} from "vue"

const counter = ref(0)

async function testWASMLoadedHandler (){
  let res = await layer8.testWASMLoaded()
  console.log("Wasm Loaded")
  counter.value = res
}


async function ping8000 () {
  try {
    let response = await layer8.fetch("http://localhost:8000/", {
        method: "POST",  
        headers: {
          "Content-Type": "Application/Json"
        },
        body: JSON.stringify({
          email: "registerEmail.value",
          password: "registerPassword.value"
        })
      });
      console.log("Ping to 8000 succeeded from navbar: ", await response.text())
  } catch (error) {
    console.log("Ping to 8000 failed from navbar: ", error);
  }
};



</script>

<template>
<div id="navbar">
    <RouterLink to="/">Home</RouterLink>
    <RouterLink to="/stress-test">Stress test</RouterLink>
    <button @click="ping8000">Ping 8000</button>
    <button @click="testWASMLoadedHandler">TestWASM</button>
    <span v-if="counter != 0 " >{{ counter }}</span >
    <br><hr><br>
</div>
</template>

<style>
</style>