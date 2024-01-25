<script setup>
import { computed, ref } from "vue";
import Navbar from "../components/Navbar.vue";
import { useRouter } from "vue-router";

const router = useRouter();
const registerEmail = ref("");
const registerPassword = ref("");
const loginEmail = ref("");
const loginPassword = ref("");
const profileImage = ref(null);
const isRegister = ref(false);
const isLoggedIn = computed(() => SpToken.value !== null);
const isContinueAnonymously = ref(false);
const SpToken = ref(localStorage.getItem("SP_TOKEN") || null);
const user = ref(localStorage.getItem("_user") ? JSON.parse(localStorage.getItem("_user")) : null);

const isLoaded = ref(false);

// ----
// TODO:
// Keep the backend URL in the .env file
// ----
//const BackendURL = "https://container-service-3.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com";
const BackendURL = "http://localhost:5001";


const registerUser = async () => {
  try {
    await layer8.fetch(BackendURL + "/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "Application/Json",
      },
      body: JSON.stringify({
        email: registerEmail.value,
        password: registerPassword.value,
        profile_image: profileImage.value,
      }),
    });
    alert("Registration successful!");
  } catch (error) {
    console.log(error);
    alert("Registration failed!");
    isRegister.value = true;
  }
};

const loginUser = async () => {
  if (loginEmail.value == "" || loginPassword.value == "") {
    console.log("login failed. Input needed");
    throw new Error("input needed");
    return;
  }

  try {
    const response = await layer8.fetch(BackendURL + "/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "Application/Json",
      },
      body: JSON.stringify({
        email: loginEmail.value,
        password: loginPassword.value,
      }),
    });

    const data = await response.json();
    if (response.status !== 200) {
      throw new Error(data.error);
    }
    SpToken.value = data.token;
    user.value = data.user;
    localStorage.setItem("SP_TOKEN", data.token);
    localStorage.setItem("_user", JSON.stringify(data.user));
    alert("Login successful!");
  } catch (error) {
    alert(error.message);
  }
};

const continueAnonymously = () => {
  isContinueAnonymously.value = true;
  alert("You are now logged in anonymously!");
  router.push({ name: "home" });
};

const logoutUser = () => {
  SpToken.value = null;
  localStorage.removeItem("SP_TOKEN");
  localStorage.removeItem("_user");
  isContinueAnonymously.value = false;
};

// Layer8 Components start here
const loginWithLayer8Popup = async () => {
  const response = await layer8.fetch(BackendURL + "/api/login/layer8/auth")
  const data = await response.json()
  // create opener window
  const popup = window.open(data.authURL, "Login with Layer8", "width=600,height=600");

  window.addEventListener("message", async (event) => {
    if (event.data.redr) {
      console.log("event.data.redr: ", event.data.redr)
      setTimeout(() => {
        layer8.fetch(BackendURL + "/api/login/layer8/auth", {
          method: "POST",
          headers: {
            "Content-Type": "Application/Json"
          },
          body: JSON.stringify({
            callback_url: event.data.redr,
          })
        })
          .then(res => res.json())
          .then(data => {
            localStorage.setItem("L8_TOKEN", data.token)
            router.push({ name: 'home' })
            popup.close();
          })
          .catch(err => console.log(err))
      }, 1000);
    }
  });
}

const uploadProfilePicture = async (e) => {
  e.preventDefault();

  isLoaded.value = true;

  const file = e.target.files[0];
  layer8.fetch(BackendURL + "/api/profile/upload", {
    method: "POST",
    headers: {
      "Content-Type": "application/layer8.buffer+json",
    },
    body: {
      name: file.name,
      size: file.size,
      type: file.type,
      buff: await file.arrayBuffer(),
    },
  })
    .then((res) => res.json())
    .then((data) => {
      profileImage.value = data.url;
    })
    .catch((err) => console.log(err))
    .finally(() => {
      isLoaded.value = false;
    });
};
// Layer8 Components end here
</script>

<template>
  <div class="h-screen bg-primary flex flex-col">
    <Navbar></Navbar>
    <div class="bg-primary-content w-full flex justify-center items-center p-4 flex-1">
      <div class="card w-auto bg-base-100 shadow-xl p-8 max-w-xs h-min" v-if="!isLoggedIn">
        <div v-if="isRegister" class="flex gap-3 flex-col">
          <h2 class="text-lg font-bold ">Register</h2>
          <input v-model="registerEmail" placeholder="Email" class="input input-bordered input-primary w-full max-w-xs"/>
          <input v-model="registerPassword" type="password" placeholder="Password"  class="input input-bordered input-primary w-full max-w-xs"/>
          <hr />
          <h1 class="text-dark pb-4 font-bold">Upload Profile Picture</h1>
          <input type="file" @change="uploadProfilePicture" />
          <div v-if="profileImage">
            <img :src="profileImage" />
          </div>
          <hr />
          <button class="btn btn-primary max-w-xs" @click="registerUser" :disabled="isLoaded">
            <div v-if="isLoaded" class="loading"></div>Register
          </button>
          <a class="block" @click="isRegister = false">Already registered? Login</a>
        </div>

        <div v-if="!isRegister"  class="flex gap-3 flex-col">
          <h2  class="text-lg font-bold">Login</h2>
          <input v-model="loginEmail" placeholder="Email" class="input input-bordered input-primary w-full max-w-xs"/>
          <input v-model="loginPassword" type="password" placeholder="Password" class="input input-bordered input-primary w-full max-w-xs" />
          <button class="btn btn-primary max-w-xs" @click="loginUser">Login</button>
          <a class="block" @click="isRegister = true">Don't have an account? Register</a>
        </div>
      </div>

      <div v-if="isLoggedIn" class="card w-auto bg-base-100 shadow-xl p-8 max-w-xs">
        <h1 class="text-dark pb-4 font-bold">
          Welcome {{ user?.email }}!
        </h1>
        <div v-if="user?.profile_image">
          <img :src="user?.profile_image" />
          <br />
          <hr />
          <br />
        </div>
        <div class="flex flex-col gap-4" v-if="!isContinueAnonymously">
          <button class="btn " @click="continueAnonymously">
            Login Anonymously
          </button>
          <button class="btn " @click="loginWithLayer8Popup">
            Login with Layer8
          </button>
          <button class="btn " @click="logoutUser">Logout</button>
        </div>
      </div>
    </div>
  </div>
</template>
