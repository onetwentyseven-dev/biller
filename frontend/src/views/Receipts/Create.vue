<script lang="ts">
import { ProviderRequest, ReceiptRequest } from '@/api';
import type { IProvider } from '@/api/types/provider';
import type { ICreateUpdateReceipt } from '@/api/types/receipt';
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import DashboardLayout from '../layouts/DashboardLayout.vue';
import Loading from '@/components/Loading.vue';
import Datepicker from 'vue3-datepicker';

export default {
  name: 'CreateReceipt',
  setup() {
    const loading = ref<Boolean>(true);
    const loadingMessage = 'Creating Receipt';
    const errorMessage = ref<String | undefined>('');
    const receipt: ICreateUpdateReceipt = reactive({
      label: '',
      amount_paid: 0,
      date_paid: new Date(),
    });
    const providers = ref<IProvider[]>([]);
    const router = useRouter();
    const createReceipt = async () => {
      if (!receipt) return;
      errorMessage.value = undefined;
      loading.value = true;
      await ReceiptRequest.Create(receipt).then((r) => {
        router.push({ name: 'receipts', params: { receiptID: r.id } });
      });
    };

    ProviderRequest.List().then((r) => {
      providers.value = r;
      loading.value = false;
    });

    return {
      loading,
      loadingMessage,
      errorMessage,
      receipt,
      providers,
      createReceipt,
    };
  },
  components: { DashboardLayout, Loading, Datepicker },
};
</script>

<template>
  <DashboardLayout>
    <h3>Create A Receipt</h3>
    <hr />
    <Loading :message="loadingMessage" v-if="loading" />
    <div class="container" v-else>
      <div class="class" v-if="errorMessage">
        <div class="col">
          <div class="alert alert-danger">
            <h4>Error Encountered Uploading Receipt</h4>
            {{ errorMessage }}
          </div>
        </div>
      </div>
      <div class="row" v-else>
        <div class="col-lg-6 offset-3">
          <div class="card">
            <div class="card-header">Receipt Details</div>
            <div class="card-body">
              <div class="mb-3">
                <label for="label">Receipt Label</label>
                <input type="text" v-model="receipt.label" class="form-control" />
              </div>
              <div class="mb-3">
                <label for="provider_id">Provider</label>
                <select
                  name="provider_id"
                  id="provider_id"
                  v-model="receipt.provider_id"
                  class="form-select"
                >
                  <option v-for="provider in providers" :key="provider.id" :value="provider.id">
                    {{ provider.name }}
                  </option>
                </select>
              </div>
              <div class="mb-3">
                <label for="amount_due">Amount Due</label>
                <input
                  type="text"
                  id="amount_paid"
                  v-model.number="receipt.amount_paid"
                  class="form-control"
                />
              </div>
              <div class="mb-3">
                <label for="date_paid">Date Paid</label>
                <Datepicker v-model="receipt.date_paid" class="form-control" />
              </div>
              <!-- TODO: Add Filepicker so that file can be uploaded at the sametime -->
            </div>
            <div class="card-footer d-grid gap-2">
              <button class="btn btn-primary" @click="createReceipt">Create Receipt</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
