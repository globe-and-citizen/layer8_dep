<script setup>
import { computed, ref } from "vue";
import Navbar from "../components/Navbar.vue";

const registerEmail = ref("");
const registerPassword = ref("");
const loginEmail = ref("");
const loginPassword = ref("");
const isRegister = ref(false);
const token = ref(localStorage.getItem("token") || null);

const isLoggedIn = computed(() => token.value !== null);

const registerUser = async () => {
  try {
    await layer8.fetch("http://localhost:5000/api/register", {
        method: "POST",  
        headers: {
          "Content-Type": "Application/Json"
        },
        body: JSON.stringify({
          email: registerEmail.value,
          password: registerPassword.value
        })
      });
    alert("Registration successful!");
  } catch (error) {
    console.log(error);
    alert("Registration failed!");
    isRegister.value = true
  }
};

const loginUser = async () => {
  try {
    const response_as_string = await layer8.fetch("http://localhost:5000/api/login", {
      method: "POST",
      headers: {
        "Content-Type": "Application/Json"
      },
      body: JSON.stringify({
        email: loginEmail.value,
        password: loginPassword.value
      })
    })      

    token.value =  response_as_string
    localStorage.setItem("token", response_as_string);
    alert("Login successful!");
  } catch (error) {
    console.error(error);
    alert("Login failed!");
  }
};

const logoutUser = () => {
  token.value = null;
  localStorage.removeItem("token");
};

const userEmail = computed(() => {
  if (token.value) {
    const payload = JSON.parse(atob(token.value.split(".")[1]));
    return payload.email;
  }
  return "";
});
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
      <h2>Welcome, {{ userEmail }}</h2>
      <button class="btn-primary" @click="logoutUser">Logout</button>
    </div>
  </div>
  <div>
    
  </div>
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
