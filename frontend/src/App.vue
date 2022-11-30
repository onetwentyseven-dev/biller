<script lang="ts">
import { ref, defineComponent } from 'vue';
import { RouterView } from 'vue-router';

import Navbar from './components/Navbar.vue';
import Loading from '@/components/Loading.vue';
import { WarmerRequest } from '@/api/index';

export default defineComponent({
  components: {
    Loading,
    Navbar,
  },
  setup() {
    const loading = ref(true);

    const warmup = async () => {
      await WarmerRequest.Get()
        .then(() => {
          loading.value = false;
        })
        .catch(() => {
          setTimeout(warmup, 5000);
        });
    };

    warmup();

    return {
      loading,
    };
  },
});
</script>

<template>
  <Navbar />
  <div class="container" v-if="loading">
    <div class="row">
      <div class="col d-flex justify-content-center align-items-center">
        <Loading message="Loading Application....Please Hang Tight" />
      </div>
    </div>
  </div>
  <RouterView v-else />
</template>

<style scoped></style>
