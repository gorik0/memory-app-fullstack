<script >

import { defineComponent, ref, onMounted } from 'vue';

// Wrapping exported object in define component
// gives us typing help! Woot!
export default defineComponent({
  name: 'App',
  setup() {
    const errorCode = ref(null);
    const errorMessage = ref(null);

    onMounted(async () => {
      const response = await fetch('/api/account/me', {
        method: 'GET',
      });

      const body = await response.json();

      errorCode.value = response.status;
      errorMessage.value = body.error.message;
    });

    return {
      errorCode,
      errorMessage,
    };
  },
});
</script>

<template>
  <h1>Helo GORIK</h1>
  <h1>Scaffolded App Works Well!</h1>
  <h3 v-if="errorCode">Error code: {{ errorCode }}</h3>
  <h3 v-if="errorMessage">{{ errorMessage }}</h3>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
