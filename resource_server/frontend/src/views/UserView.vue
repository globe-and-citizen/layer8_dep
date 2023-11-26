<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const token = ref(localStorage.getItem("token") || null);

let user = ref({
  email: "",
  username: "",
  first_name: "",
  last_name: "",
  display_name: "",
  country: "",
  email_verified: false,
});

const getUserDetails = async () => {
  try {
    const resp = await window.fetch("http://localhost:3050/api/v1/profile", {
      method: "GET",
      headers: {
        "Content-Type": "Application/Json",
        Authorization: `Bearer ${token.value}`,
      },
    });
    const data = await resp.json();
    user.value = data;
  } catch (error) {
    console.error(error);
  }
};

// Call decodeToken when the component is mounted
onMounted(() => {
  getUserDetails();
});

const logoutUser = () => {
  token.value = null;
  localStorage.removeItem("token");
  router.push("/");
};

const verifyEmail = async () => {
  try {
    const resp = await window.fetch(
      "http://localhost:3050/api/v1/verify-email",
      {
        method: "POST",
        headers: {
          "Content-Type": "Application/Json",
          Authorization: `Bearer ${token.value}`,
        },
      }
    );
    const data = await resp.json();
    if (data.message === "OK!") {
      alert("Email verified!");
      user.value.email_verified = true;
    } else {
      alert("Email verification failed!");
    }
  } catch (error) {
    console.error(error);
  }
};

const changeDisplayName = async () => {
  try {
    const newDisplayName = prompt("Enter new display name:");
    if (newDisplayName == "") {
      alert("Please enter a display name!");
      return;
    }
    const resp = await window.fetch(
      "http://localhost:3050/api/v1/change-display-name",
      {
        method: "POST",
        headers: {
          "Content-Type": "Application/Json",
          Authorization: `Bearer ${token.value}`,
        },
        body: JSON.stringify({
          display_name: newDisplayName,
        }),
      }
    );
    const data = await resp.json();
    if (data.message === "OK!") {
      alert("Display name changed!");
      user.value.display_name = newDisplayName;
    } else {
      alert("Display name change failed!");
    }
  } catch (error) {
    console.error(error);
  }
};
</script>

<template>
  <div class="header">
    <img src="../assets/L8Logo.png" alt="Layer8" width="500" height="100" />
  </div>
  <div id="app">
    <div class="container">
      <div class="form-container">
        <h1
          style="
            font-weight: bold;
            padding-left: 10%;
            padding-bottom: 3%;
            color: black;
          "
        >
          User Profile
        </h1>
        <p class="box">Username: {{ user.username }}</p>
        <hr class="line" />
        <p class="box">First Name: {{ user.first_name }}</p>
        <hr class="line" />
        <p class="box">Last Name: {{ user.last_name }}</p>
        <hr class="line" />
        <p class="box">Email: {{ user.email }}</p>
        <hr class="line" />
        <p class="box">Country: {{ user.country }}</p>
        <hr class="line" />
        <p class="box">Display Name: {{ user.display_name }}</p>
        <button
          class="btn-primary"
          style="margin-left: 16%"
          @click="changeDisplayName"
        >
          Change Display Name
        </button>
        <hr class="line" />
        <p class="box">Email Verified: {{ user.email_verified }}</p>
        <button
          class="btn-primary"
          style="margin-left: 28%"
          v-if="!user.email_verified"
          @click="verifyEmail"
        >
          Verify Email
        </button>
        <hr class="line" />
        <button
          class="btn-primary-2"
          style="margin-top: 4%"
          @click="logoutUser"
        >
          Logout
        </button>
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

.box {
  padding: 10px;
  color: #000000;
  font-size: 14px;
  font-family: monospace;
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

.btn-primary {
  background-color: #4caf50;
  color: white;
  border: none;
  padding: 4px 7px;
  cursor: pointer;
  border-radius: 5px;
  font-size: 10px;
  font-family: monospace;
  margin-bottom: 4%;
}

.btn-primary:hover {
  background-color: #45a049;
  transition-duration: 0.5s;
}

.btn-primary-2 {
  background-color: #1b54b1;
  color: white;
  border: none;
  padding: 8px 10px;
  cursor: pointer;
  border-radius: 5px;
  font-size: 1rem;
  margin-left: 34%;
  font-family: monospace;
}

.line {
  width: 100%;
  color: black;
  height: 4px;
  background-color: black;
}
</style>
