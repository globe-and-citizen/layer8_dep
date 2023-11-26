<script setup>
import { computed, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const registerEmail = ref("");
const registerUsername = ref("");
const registerPassword = ref("");
const registerFirstName = ref("");
const registerLastName = ref("");
const registerDisplayName = ref("");
const loginUsername = ref("");
const loginPassword = ref("");
const isRegister = ref(false);
const token = ref(localStorage.getItem("token") || null);

const isLoggedIn = computed(() => token.value !== null);

const registerUser = async () => {
  try {
    if (registerUsername.value == "" || registerPassword.value == "") {
      alert("Please enter a username and password!");
      return;
    }
    await window.fetch("http://localhost:3050/api/v1/register-user", {
      method: "POST",
      headers: {
        "Content-Type": "Application/Json",
      },
      body: JSON.stringify({
        email: registerEmail.value,
        username: registerUsername.value,
        first_name: registerFirstName.value,
        last_name: registerLastName.value,
        password: registerPassword.value,
        display_name: registerDisplayName.value,
      }),
    });
    alert("Registration successful!");
    isRegister.value = true;
  } catch (error) {
    console.error(error);
    alert("Registration failed!");
  }
};

const loginUser = async () => {
  try {
    if (loginUsername.value == "" || loginPassword.value == "") {
      alert("Please enter a username and password!");
      return;
    }
    const respOne = await window.fetch(
      "http://localhost:3050/api/v1/login-precheck",
      {
        method: "POST",
        headers: {
          "Content-Type": "Application/Json",
        },
        body: JSON.stringify({
          username: loginUsername.value,
        }),
      }
    );
    const responseOne = await respOne.json();
    console.log(responseOne.salt);
    const respTwo = await window.fetch(
      "http://localhost:3050/api/v1/login-user",
      {
        method: "POST",
        headers: {
          "Content-Type": "Application/Json",
        },
        body: JSON.stringify({
          username: loginUsername.value,
          password: loginPassword.value,
          salt: responseOne.salt,
        }),
      }
    );
    const responseTwo = await respTwo.json();
    console.log(responseTwo);
    if (responseTwo.token) {
      token.value = responseTwo.token;
      localStorage.setItem("token", responseTwo.token);
      alert("Login successful!");
      router.push("/user");
    } else {
      console.error("No token received.");
    }
  } catch (error) {
    console.error(error);
    alert("Login failed!");
  }
};
</script>

<template>
  <div class="header">
    <img src="../assets/L8Logo.png" alt="Layer8" width="500" height="100" />
  </div>
  <div id="app">
    <div class="container" v-if="!isLoggedIn">
      <div v-if="isRegister" class="form-container">
        <h2
          style="
            margin-left: 22%;
            margin-bottom: 8%;
            font-weight: normal;
            color: black;
            font-family: monospace;
          "
        >
          Register
        </h2>
        <div class="input-group">
          <input
            v-model="registerEmail"
            placeholder="Email"
            class="input-button"
          />
        </div>
        <div class="input-group">
          <input
            v-model="registerUsername"
            placeholder="Username"
            class="input-button"
          />
        </div>
        <div class="input-group">
          <input
            v-model="registerFirstName"
            placeholder="First Name"
            class="input-button"
          />
        </div>
        <div class="input-group">
          <input
            v-model="registerLastName"
            placeholder="Last Name"
            class="input-button"
          />
        </div>
        <div class="input-group">
          <input
            v-model="registerDisplayName"
            placeholder="Display Name"
            class="input-button"
          />
        </div>
        <div class="input-group">
          <input
            v-model="registerPassword"
            type="password"
            placeholder="Password"
            class="input-button"
          />
        </div>
        <button class="btn-primary" @click="registerUser">Register</button>
        <a
          class="text"
          style="display: block; cursor: pointer"
          @click="isRegister = false"
          >Already registered? Login</a
        >
      </div>

      <div v-if="!isRegister" class="form-container">
        <h2
          style="
            margin-left: 34%;
            margin-bottom: 8%;
            font-weight: normal;
            color: black;
            font-family: monospace;
          "
        >
          Login
        </h2>
        <div class="input-group">
          <input
            v-model="loginUsername"
            placeholder="Username"
            class="input-button"
          />
        </div>
        <div class="input-group">
          <input
            v-model="loginPassword"
            type="password"
            placeholder="Password"
            class="input-button"
          />
        </div>
        <button class="btn-primary-2" @click="loginUser">Login</button>
        <a
          class="text"
          style="display: block; cursor: pointer"
          @click="isRegister = true"
          >Don't have an account? Register</a
        >
      </div>
    </div>
  </div>
</template>

<style scoped>
#app {
  display: flex;
  justify-content: center;
  padding: 10%;
  background-color: #ffffff;
}

.header {
  display: flex;
  justify-content: space-around;
  background-color: rgb(255, 255, 255);
  padding-top: 1%;
}

.container {
  display: flex;
  justify-content: space-around;
  background-color: rgb(255, 255, 255);
  padding: 3rem;
  border-radius: 1rem;
  border-color: #000000;
  border-style: solid;
  border-width: 5px;
}

.form-container {
  width: 100%;
}

.input-group {
  margin-bottom: 1rem;
}

.input-button {
  width: 100%;
  padding: 6px 15px;
  border-radius: 5px;
  border: #adadad;
  border-style: solid;
  border-width: 2px;
  font-size: 13px;
  font-family: monospace;
}

.btn-primary {
  background-color: #3440ab;
  color: white;
  border: none;
  padding: 8px 10px;
  cursor: pointer;
  border-radius: 5px;
  font-size: 1rem;
  margin-left: 25%;
  font-family: monospace;
}

.btn-primary-2 {
  background-color: #4caf50;
  color: white;
  border: none;
  padding: 8px 10px;
  cursor: pointer;
  border-radius: 5px;
  font-size: 1rem;
  margin-left: 34%;
  font-family: monospace;
}

.btn-primary:hover {
  background-color: #190e96;
  transition-duration: 0.5s;
}

.btn-primary-2:hover {
  background-color: #45a049;
  transition-duration: 0.5s;
}

.text {
  color: rgb(37, 37, 37);
  font-size: 13px;
  text-decoration: underline;
  margin-top: 10px;
  text-align: center;
}

.text:hover {
  color: rgb(35, 108, 38);
  transition-duration: 0.5s;
}
</style>
