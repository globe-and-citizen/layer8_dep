<script setup>
import { ref } from "vue";
import { RouterLink } from "vue-router";

const requestsSent = ref(0);
const totalTimeSpent = ref(0);
const numberOfRequest = ref(0)

async function testWASMLoadedHandler (){
  let res = await layer8.testWASMLoaded()
  console.log(res)
}


async function testWASMHandler() {
  const startTime = performance.now();
  for (let i = 0; i < numberOfRequest.value; i++) {
    const res = await layer8.testWASM(i, "42");
    requestsSent.value++;
    console.log(res);
  }
  const endTime = performance.now();
  totalTimeSpent.value = endTime - startTime;
  console.log("Total request sent: ", requestsSent.value)
  console.log("Total time spent: ", totalTimeSpent.value, "ms")
}
</script>

<template>
    <div id="navbar">
    <RouterLink to="/">Home</RouterLink>
    <RouterLink to="/stress-test">Stress test</RouterLink>
    <button @click="testWASMLoadedHandler">TestWASM</button>
    <br><hr><br>
  </div>
  <div class="greetings">
    <div>
      <label for="">Number of request</label>
      <input type="text" v-model="numberOfRequest">
      <button @click="testWASMHandler" class="text-green-500">Execute</button>

      <div>
        Total request sent: {{ requestsSent }} times
      </div>

      <div>
        Total time spent: {{ totalTimeSpent }} ms
      </div> 
    </div>
  </div>
</template>

<style scoped></style>
