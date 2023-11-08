<script setup>
import { computed, ref } from "vue";
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const code = ref(new URLSearchParams(window.location.search).get("code"))
const token = ref(localStorage.getItem("token") || null)

onMounted(() => {
    setTimeout(() => {
        layer8.fetch("http://localhost:8000/api/login/layer8/auth", {
            method: "POST",
            headers: {
                "Content-Type": "Application/Json"
            },
            body: JSON.stringify({
                callback_url: window.location.href,
            })
        })
            .then(res => res.json())
            .then(data => {
                localStorage.setItem("token", data.token)
                router.push({ name: 'home' })
            })
            .catch(err => console.log(err))
    }, 1000);
})
</script>

<template>
    <div>
        <h1>Login with layer8...</h1>
    </div>
</template>

