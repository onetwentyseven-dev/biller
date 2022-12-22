<script lang="ts">
import { ref, defineComponent } from 'vue';
import { RouterView } from 'vue-router';

import Navbar from './components/Navbar.vue';
import Loading from '@/components/Loading.vue';
import { WarmerRequest } from '@/api/index';
// import { useAuth0 } from '@auth0/auth0-vue';

export default defineComponent({
  components: {
    Loading,
    Navbar,
  },
  setup() {
    const loading = ref(true);
    // const auth0 = useAuth0();

    const warmup = async () => {
      await WarmerRequest.Get()
        .then(() => {
          loading.value = false;
        })
        .catch((e) => {
          console.log('caught ya bitch', e);
          setTimeout(warmup, 5000);
        });
    };

    warmup();

    // console.log(auth0.idTokenClaims.value);
    // auth0.getAccessTokenSilently().then((token) => {
    //   console.log(token);
    // });

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
