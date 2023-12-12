<script setup>
import { computed, ref } from "vue";
import Navbar from "../components/Navbar.vue";
import { useRouter } from "vue-router";

const router = useRouter();
const registerEmail = ref("");
const registerPassword = ref("");
const loginEmail = ref("");
const loginPassword = ref("");
const isRegister = ref(false);
const isLoggedIn = computed(() => SpToken.value !== null);
const isContinueAnonymously = ref(false);
const SpToken = ref(localStorage.getItem("SP_TOKEN") || null);

const registerUser = async () => {
  try {
    // await layer8.fetch("http://localhost:5001/api/register", {
    await layer8.fetch("https://container-service-3.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com/api/register", {
      method: "POST",
      headers: {
        "Content-Type": "Application/Json",
      },
      body: JSON.stringify({
        email: registerEmail.value,
        password: registerPassword.value,
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
    // const response = await layer8.fetch("http://localhost:5001/api/login", {
      const response = await layer8.fetch("https://container-service-3.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com/api/login", {
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
    SpToken.value = data.token;
    localStorage.setItem("SP_TOKEN", data.token);
    alert("Login successful!");
  } catch (error) {
    console.error(error);
    alert("Login failed!");
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
  isContinueAnonymously.value = false;
};

const userName = computed(() => {
  if (SpToken.value && SpToken.value.split(".").length > 1) {
    const payload = JSON.parse(atob(SpToken.value.split(".")[1]));
    return payload.email;
  }
  return "";
});

const loginWithLayer8Popup = async () => {
  // const response = await layer8.fetch("http://localhost:8000/api/login/layer8/auth")
  const response = await layer8.fetch("https://container-service-3.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com/api/login/layer8/auth")
  const data = await response.json()

  //alert(data.authURL)
  // create opener window
  const popup = window.open(data.authURL, "Login with Layer8", "width=600,height=600");

  window.addEventListener("message", async (event) => {
    if (event.data.redr) {
      console.log("event.data.redr: ", event.data.redr)
      setTimeout(() => {
        // layer8.fetch("http://localhost:8000/api/login/layer8/auth", {
        layer8.fetch("https://container-service-3.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com/api/login/layer8/auth", {
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
</script>

<template>
  <Navbar></Navbar>
  <div id="app">
    <div class="container" v-if="!isLoggedIn">
      <div v-if="isRegister" class="form-container">
        <h2>Register</h2>
        <div class="input-group">
          <input v-model="registerEmail" placeholder="Email" />
        </div>
        <div class="input-group">
          <input v-model="registerPassword" type="password" placeholder="Password" />
        </div>
        <button class="btn-primary" @click="registerUser">Register</button>
        <a style="display: block" @click="isRegister = false">Already registered? Login</a>
      </div>

      <div v-if="!isRegister" class="form-container">
        <h2>Login</h2>
        <div class="input-group">
          <input v-model="loginEmail" placeholder="Email" />
        </div>
        <div class="input-group">
          <input v-model="loginPassword" type="password" placeholder="Password" />
        </div>
        <button class="btn-primary" @click="loginUser">Login</button>
        <a style="display: block" @click="isRegister = true">Don't have an account? Register</a>
      </div>
    </div>

    <div v-if="isLoggedIn" class="welcome-container">
      <h1 style="color: rgb(136, 136, 136); font-weight: 600; padding-bottom: 2%">
        Welcome {{ userName }}!
      </h1>
      <div class="new-container" v-if="!isContinueAnonymously">
        <button class="btn-secondary" @click="continueAnonymously">
          Login Anonymously
        </button>
        <button class="btn-secondary" @click="loginWithLayer8Popup">
          Login with Layer8
        </button>
        <button class="btn-secondary" @click="logoutUser">Logout</button>
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

.container {
  display: flex;
  justify-content: space-around;
  width: 50%;
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.form-container {
  width: 100%;
}

.input-group {
  margin-bottom: 15px;
}

.btn-primary {
  background-color: #4caf50;
  color: white;
  border: none;
  padding: 10px 20px;
  cursor: pointer;
  border-radius: 5px;
  font-size: 16px;
  transition: background-color 0.3s;
}

.btn-primary:hover {
  background-color: #45a049;
}

.welcome-container {
  text-align: center;
  width: 100%;
}
</style>
