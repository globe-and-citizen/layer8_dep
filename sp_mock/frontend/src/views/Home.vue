<script setup>
import { computed, ref } from "vue";
import Navbar from "../components/Navbar.vue";

const token = ref(localStorage.getItem("token") || null);
const isLoggedIn = computed(() => token.value !== null);
let newPoem = ref("");

const userEmail = computed(() => {
    if (token.value && token.value.split(".").length > 1) {
        const payload = JSON.parse(atob(token.value.split(".")[1]));
        return payload.email;
    }
    return "";
});

</script>

<template>
    <Navbar></Navbar>
    <div id="app">
        <div v-if="isLoggedIn" class="welcome-container">
            <h1 style="color: rgb(136, 136, 136); font-weight: 600; padding-bottom: 2%">
                Welcome {{ userEmail }}!
            </h1>
            <div class="new-container">
                <button class="btn-secondary" @click="logoutUser">Logout</button>
            </div>
            <div id="poems-container-2" style="color: black">
                <div id="newPoem">
                    {{ newPoem }}
                </div>
            </div>
        </div>
    </div>
    <div></div>
</template>


<style scoped>
#app {
    font-family: "Arial", sans-serif;
    display: flex;
    justify-content: center;
    align-items: center;
    height: 80vh;
    width: 100vh;
    background-color: #f4f4f4;
}

.welcome-container {
    text-align: center;
    width: 100%;
}
</style>
