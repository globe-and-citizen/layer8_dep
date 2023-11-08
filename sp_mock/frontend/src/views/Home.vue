<script setup>
import { computed, ref } from "vue";
import Navbar from "../components/Navbar.vue";
import { useRouter } from "vue-router";

const router = useRouter();
const token = ref(localStorage.getItem("token") || null);
const isLoggedIn = computed(() => token.value !== null);
let newPoem = ref("");
const gotAPoem = ref(false);

const userEmail = computed(() => {
    if (token.value && token.value.split(".").length > 1) {
        const payload = JSON.parse(atob(token.value.split(".")[1]));
        return payload.email;
    }
    return "";
});

const getPoems = async () => {
    // try {
    //   const poem = await layer8.getPoem({
    //     method: "GET",
    //     headers: {
    //       "Content-Type": "application/json",
    //       "Authorization": `${token.value}`,
    //     }
    //   });
    //   if (poem) {
    //     showPoem(poem);
    //     gotAPoem.value = true;
    //   } else {
    //     console.error("No poem content received.");
    //   }
    // } catch (error) {
    //   console.error(error);
    //   alert("Poems failed!");
    // }

    // Since poems are not implemented in the new mock, just using a placeholder here
    const poem =
        "Roses are red, violets are blue, I'm a placeholder, and so are you.";
    showPoem(poem);
    gotAPoem.value = true;
};

const showPoem = (poemText) => {
    const poemContainer = document.getElementById("poems-container-2");
    const newPoem = document.createElement("div");
    newPoem.id = "newPoem";
    newPoem.textContent = poemText;
    poemContainer.appendChild(newPoem);
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
                Welcome {{ userEmail }}!
            </h1>
            <div class="new-container">
                <button @click="getPoems">
                    Get Poems
                </button>
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
