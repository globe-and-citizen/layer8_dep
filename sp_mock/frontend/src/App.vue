<script setup>
import HelloWorld from './components/HelloWorld.vue'
import TheWelcome from './components/TheWelcome.vue'

async function testWASMHandler(){
  const res = await window.testWASM(42, "42")
  console.log(res)
}

async function pingProxyHandler(){
  try{
    const res = await window.pingProxy("http://localhost:5000/success")//proxy
    console.log("Response from JS: ", res)
  } catch(err){
    console.log("Err from JS:", err)
  }
}

async function genericPostHandler(){
  try{
    const res = await window.genericPost("http://localhost:5000/success", "This is my content!")
    console.log("Response from JS: ", res)
  } catch(err){
    console.log("Err from JS: ", err)
  }
}

</script>

<template>
  <header>
    <img alt="Vue logo" class="logo" src="./assets/logo.svg" width="125" height="125" />

    <div class="wrapper">
      <HelloWorld msg="We Got Poems Mock" />
      <button @click="testWASMHandler">Test WASM</button>
      <button @click="pingProxyHandler">Ping Proxy</button>
      <button @click="genericPostHandler">Trigger Post Handler</button>
    </div>
  </header>

  <!-- <main>
    <TheWelcome />
  </main> -->
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
