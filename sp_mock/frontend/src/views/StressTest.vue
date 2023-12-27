<script setup>
import { ref } from "vue";
import Navbar from "../components/Navbar.vue";

const requestsSent = ref(0);
const totalTimeSpent = ref(0);
const numberOfRequest = ref(0)

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
  <Navbar></Navbar>
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
