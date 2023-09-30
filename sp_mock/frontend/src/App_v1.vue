<script setup>
import HelloWorld from './components/HelloWorld.vue'
import {ref} from "vue"

let username = ref("chester")
let password = ref("tester")

async function login(){
  console.log(username.value, password.value)
  try{
    const res = await window.genericGetRequest(`http://localhost:5000/login?username=${username.value}&password=${password.value}`)//proxy
    console.log("Response from JS: ", res)
  } catch(err){
    console.log("Err from JS:", err)
  }
}

async function register(){
  const genericObject = {
    username: username.value,
    password: password.value
  }


  try{
    const res = await window.genericPost("http://localhost:5000/register", JSON.stringify(genericObject))
    console.log("Response from JS: ", res)
  } catch(err){
    console.log("Err from JS: ", err)
  }
}


async function testWASMHandler(){
  const res = await window.testWASM(42, "42")
  console.log(res)
}

async function genericGetRequest(){
  try{
    const res = await window.genericGetRequest("http://localhost:5000/success")//proxy
    console.log("Response from JS: ", res)
  } catch(err){
    console.log("Err from JS:", err)
  }
}

async function genericPostHandler(){
  const genericObject = {
    one: "1",
    two: false,
    three: {
      key: "value"
    },
    four: [1,2,3]
  }


  try{
    const res = await window.genericPost("http://localhost:5000/success", JSON.stringify(genericObject))
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
      <button @click="genericGetRequest">Generic GET Request</button>
      <button @click="genericPostHandler">Generic POST Request (Sends Object)</button>
    </div>
  </header>


  <hr>
  <br>

  <div>
    <h2>login</h2>
    <div>
      <label for="username">username </label>
      <input type="text" id="username" v-model="username">
    </div>
    <div>
      <label for="password">password </label>
      <input type="text" id="password" v-model="password">
    </div>
    <button @click="login">login</button>
    <button @click="register">register</button>
    <!-- <butto></button> -->
  </div>

  <hr>


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