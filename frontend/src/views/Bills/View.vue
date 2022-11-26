<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import { BillRequest, ProviderRequest } from '@/api/index';
import { ref } from 'vue';
import type { Ref } from 'vue';
import { useRoute } from 'vue-router';

import type { IBill } from '@/api/types/bill';
import type { IProvider } from '@/api/types/provider';

export default {
  components: {
    Loading,
    DashboardLayout,
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const bill: Ref<IBill | undefined> = ref();
    const provider: Ref<IProvider | undefined> = ref();
    const route = useRoute();
    const billID = route.params.billID as string;

    BillRequest.Get(billID)
      .then(async (r) => {
        bill.value = r;
        await ProviderRequest.Get(r.provider_id).then((r) => {
          provider.value = r;
        });
      })
      .then(() => {
        loading.value = false;
      });

    return {
      loading,
      bill,
      provider,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading message="Loading Bills" v-if="loading" />
    <div v-else-if="!bill || !provider">
      <div class="alert alert-danger">
        <h4>Bill Failed to Load. Check Console</h4>
      </div>
    </div>
    <div v-else>
      <h3>Bill - {{ bill.name }}</h3>
      <hr />
      <div class="container">
        <div class="row">
          <div class="col-lg-6">
            <div class="card">
              <div class="card-header">
                <span>Bill Info</span>
              </div>
              <div class="list-group list-group-flush">
                <RouterLink
                  :to="{ name: 'provider', params: { providerID: bill.provider_id } }"
                  class="list-group-item list-group-item-action"
                >
                  Provided By: {{ provider?.name }}
                </RouterLink>
              </div>
            </div>
          </div>
          <div class="col">
            <div class="list-group">
              <RouterLink :to="{}"> </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
