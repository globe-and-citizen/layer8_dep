<script setup>
function loginWithLayer8() {
  console.log("Login with Layer 8");
}
</script>

<template>
  <div class="flex items-center justify-center">
    <div class="bg-white p-8 rounded-lg shadow-md w-96">
      <h1 class="text-3xl font-semibold mb-4">Login</h1>
      <button
        @click="loginWithLayer8"
        class="bg-[hsla(160,100%,37%,1)] hover:bg-[#28c993] flex space-x-2 items-center text-white text-base py-2 px-4 rounded-md transition duration-300 ease-in-out transform hover:scale-105"
      >
        <img class="w-5 h-5" src="@/assets/images/key-minimalistic-square.svg" alt="" />
        <span>Login with Layer 8</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
/* Add any additional custom styles you want here */
</style>
