<script setup>
import { RouterView } from "vue-router";
import layer8_interceptor from 'layer8_interceptor'

try{
  layer8_interceptor.initEncryptedTunnel({
    ServiceProviderURL: 'container-service-3.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com', // TODO: this needs to become the same form as the VITE_BACKEND_URL
    Layer8Scheme: 'https',
    Layer8Host: 'aws-container-service-t1.gej3a3qi2as1a.ca-central-1.cs.amazonlightsail.com',
    Layer8Port: ''
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


