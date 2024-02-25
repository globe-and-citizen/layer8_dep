<script setup>
import { RouterView } from 'vue-router';
import layer8 from "layer8_interceptor";
import { onMounted } from 'vue';
import emitter from '@/plugins/mitt';

onMounted(async () => {
  await layer8.initEncryptedTunnel({
    ServiceProviderURL: 'localhost:6001',
    Layer8Scheme: 'http',
    Layer8Host: 'localhost',
    Layer8Port: '5001'
  });
});

const uploadFile = (e) => {
  const file = e.target.files[0];
  const formData = new FormData();
  formData.append('file', file);

  layer8.fetch('http://localhost:6001/api/upload', {
    method: 'POST',
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    body: formData
  })
    .then((response) => response.json())
    .then(() => {
      emitter.emit('reload_images');
    });
}
</script>

<template>
  <header>
    <h1>ImSharer</h1>
    <input type="file" name="upload" class="hidden" ref="uploadFile" @change="uploadFile" />
    <input type="button" value="Upload" @click="$refs.uploadFile.click()" />
  </header>

  <RouterView />
</template>
