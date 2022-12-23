<script lang="ts">
import { defineComponent } from 'vue';
import { RouterView } from 'vue-router';

import Navbar from './components/Navbar.vue';
import Loading from '@/components/Loading.vue';

export default defineComponent({
  components: {
    Loading,
    Navbar,
  },
});
</script>

<template>
  <Navbar />
  <RouterView />
</template>

<style scoped></style>
