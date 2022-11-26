<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import { BillRequest, ProviderRequest } from '@/api/index';
import { reactive, ref } from 'vue';
import type { Ref } from 'vue';
import { useRouter } from 'vue-router';

import type { ICreateUpdateBill } from '@/api/types/bill';
import type { IProvider } from '@/api/types/provider';

export default {
  components: {
    Loading,
    DashboardLayout,
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const router = useRouter();
    const providers: Ref<IProvider[]> = ref([]);
    const bill: ICreateUpdateBill = reactive({
      name: '',
      provider_id: '',
    });

    const createBill = async () => {
      if (!bill) return;
      loading.value = true;

      await BillRequest.Create(bill).then((r) => {
        router.push({ name: 'bill', params: { billID: r.id } });
      });
    };

    Promise.all([
      ProviderRequest.List().then((r) => {
        providers.value = r;
      }),
    ]).then(() => {
      loading.value = false;
    });

    return {
      loading,
      bill,
      providers,
      createBill,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading message="Loading Bills" v-if="loading" />
    <div v-else>
      <h3>Provision A New Bill</h3>
      <hr />
      <div class="container">
        <div class="row">
          <div class="col-lg-6 offset-3">
            <div class="card">
              <div class="card-header">Provision a New Provider</div>
              <div class="card-body">
                <div class="mb-3">
                  <label for="name">Bill Name</label>
                  <input type="text" v-model="bill.name" class="form-control" />
                </div>
                <div class="mb-3">
                  <label for="name">Bill Provider</label>
                  <select v-model="bill.provider_id" class="form-select">
                    <option v-for="provider in providers" :key="provider.id" :value="provider.id">
                      {{ provider.name }}
                    </option>
                  </select>
                </div>
              </div>
              <div class="card-footer d-grid gap-2">
                <button class="btn btn-primary" @click="createBill">Create Bill</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
