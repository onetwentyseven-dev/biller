<script lang="ts">
import { ref } from 'vue';
import type { Ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';

import { ProviderRequest } from '@/api/index';
import type { IProvider } from '@/api/types/provider';
import type { IBill } from '@/api/types/bill';

export default {
  components: {
    Loading,
    DashboardLayout,
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const provider: Ref<IProvider | undefined> = ref(undefined);
    const bills: Ref<IBill[]> = ref([]);
    const router = useRouter();
    const route = useRoute();
    const providerID = route.params.providerID as string;

    const deleteProvider = () => {
      if (!provider.value) return;
      ProviderRequest.Delete(provider.value.id).then(() => {
        router.push({ name: 'providers' });
      });
    };

    const updateProvider = () => {
      if (!provider.value) return;
      // const update:I
      ProviderRequest.Update(provider.value.id, provider.value).then(() => {});
    };

    Promise.all([
      ProviderRequest.Get(providerID).then((r) => {
        provider.value = r;
      }),
      ProviderRequest.GetProviderBills(providerID).then((r) => {
        bills.value = r;
      }),
    ]).then(() => (loading.value = false));

    return {
      loading,
      provider,
      bills,
      deleteProvider,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading message="Loading Provider" v-if="loading" />
    <div v-else-if="!provider">
      <div class="alert alert-danger">
        <h4>Provider Failed to Load. Check Console</h4>
      </div>
    </div>

    <div v-else>
      <h3>{{ provider.name }}</h3>
      <hr />
      <div class="row">
        <div class="col-lg-6">
          <div class="card">
            <div class="card-header d-flex justify-content-between">
              <span> Provider Info </span>
              <div class="btn btn-sm btn-danger" @click="deleteProvider">
                <font-awesome-icon icon="fa-solid fa-trash" />
              </div>
            </div>

            <div class="list-group flush">
              <a
                :href="provider.web_address"
                v-if="provider.web_address"
                class="list-group-item list-group-item-action"
                target="_blank"
              >
                Website (Open's In New Tab)
              </a>
            </div>
          </div>
        </div>
        <div class="col-lg-6">
          <div class="card">
            <div class="card-header">Provider Bills</div>

            <div class="list-group flush">
              <RouterLink
                class="list-group-item list-group-item-action"
                :to="{ name: 'bill', params: { billID: bill.id } }"
                v-for="bill in bills"
                :key="bill.id"
              >
                {{ bill.name }}
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
