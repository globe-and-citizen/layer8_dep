<script setup>
import { RouterView } from "vue-router";
import layer8_interceptor from 'layer8_interceptor'

try{
  layer8_interceptor.initEncryptedTunnel({
    ServiceProviderURL: 'localhost:8000', // TODO: this needs to become the same form as the VITE_BACKEND_URL
    Layer8Scheme: 'http',
    Layer8Host: 'localhost',
    Layer8Port: '5001'
  })
}catch(err){
  console.log(".initEncryptedTunnel error: ", err)
}
</script>

<template>
  <div>
    <RouterView />
  </div>
</template>

<style scoped>
</style>


