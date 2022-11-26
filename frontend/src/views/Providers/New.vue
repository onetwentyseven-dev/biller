<script lang="ts">
import { ProviderRequest } from '@/api';
import type { ICreateUpdateProvider } from '@/api/types/provider';
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import DashboardLayout from '../layouts/DashboardLayout.vue';
import Loading from '@/components/Loading.vue';

export default {
  name: 'ProviderNew',
  setup() {
    const loading = ref<Boolean>(false);
    const loadingMessage = ref<String>('Creating Provider');
    const errorMessage = ref<String | undefined>();
    const provider: ICreateUpdateProvider = reactive({ name: '', web_address: '' });
    const router = useRouter();

    const createProvider = async () => {
      if (!provider) return;
      errorMessage.value = undefined;
      loading.value = true;

      await ProviderRequest.Create(provider).then((r) => {
        router.push({ name: 'provider', params: { providerID: r.id } });
      });
    };

    return {
      loading,
      loadingMessage,
      errorMessage,
      provider,
      createProvider,
    };
  },
  components: { DashboardLayout, Loading },
};
</script>

<template>
  <DashboardLayout>
    <h3>Providers</h3>
    <hr />
    <Loading :message="loadingMessage.toString()" v-if="loading" />

    <div class="container" v-else>
      <div class="row" v-if="errorMessage">
        <div class="col">
          <div class="alert alert-danger">
            <h4>Error Encountered Creating Provider</h4>
            {{ errorMessage }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="col-lg-6 offset-3">
          <div class="card">
            <div class="card-header">Provision a New Provider</div>
            <div class="card-body">
              <div class="mb-3">
                <label for="name">Provider Name</label>
                <input type="text" v-model="provider.name" class="form-control" />
              </div>
              <div class="mb-3">
                <label for="name">Provider Web Address</label>
                <input type="text" v-model="provider.web_address" class="form-control" />
              </div>
            </div>
            <div class="card-footer d-grid gap-2">
              <button class="btn btn-primary" @click="createProvider">Create Provider</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
