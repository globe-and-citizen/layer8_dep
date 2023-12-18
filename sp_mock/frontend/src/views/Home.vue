<script setup>
import { computed, ref } from "vue";
import Navbar from "../components/Navbar.vue";
import { useRouter } from "vue-router";

const router = useRouter();
const token = ref(localStorage.getItem("token") || null);
const isLoggedIn = computed(() => token.value !== null);
let nextPoem = ref({});

const userName = computed(() => {
    if (token.value && token.value.split(".").length > 1) {
        const payload = JSON.parse(atob(token.value.split(".")[1]));
        return payload.username;
    }
    return "";
});

//console.log("payload", JSON.parse(atob(token.value.split(".")[1])));

const getPoem = async () => {
    try {
        console.log("going to try now 1...")   
        const resp = await layer8.fetch("https://layer8-mock-backend-gjzwz.ondigitalocean.app/nextpoem");
        console.log("going to try now 2...")  
        
        let poemObj = await resp.json()
        
        if (poemObj.title) {
            nextPoem.value = poemObj
        } else {
            console.error("The poemObj is malformed or other error....");
        }
    } catch (error) {
      console.error("error:", error);
    }
};


const logoutUser = () => {
    token.value = null;
    localStorage.removeItem("token");
    router.push({ name: "loginRegister" });
};

</script>

<template>
    <Navbar></Navbar>
    <div id="app">
        <div v-if="isLoggedIn" class="welcome-container">
            <h1 style="color: rgb(136, 136, 136); font-weight: 600; padding-bottom: 2%">
                Welcome {{ userName }}!
            </h1>
            <div class="new-container">
                <button @click="getPoem">
                    Get Next Poem
                </button>
                <button class="btn-secondary" @click="logoutUser">Logout</button>
            </div>
            <div id="poems-container-2" style="color: black">
                <br>
                <div id="newPoem">
                    <h3>Next Poem goes here: </h3>
                    <div>Title:</div>
                    <p style="font-weight: bold;">{{ nextPoem.title }}</p>
                    <div>Author:</div>
                    <p style="font-weight: bold;">{{ nextPoem.author }}</p>
                    <div>Body:</div>
                    <p style="font-weight: bold;">{{ nextPoem.body }}</p>
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
