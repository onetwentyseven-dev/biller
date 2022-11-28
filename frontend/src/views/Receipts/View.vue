<script lang="ts">
import { ReceiptRequest } from '@/api';
import type { IReceipt } from '@/api/types/receipt';
import type { IProvider } from '@/api/types/provider';
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import { ref, type Ref } from 'vue';
import { useRoute } from 'vue-router';

export default {
  components: {
    Loading,
    DashboardLayout,
  },
  methods: {
    formatDate(v: string): string {
      return new Date(v).toDateString();
    },
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const receipt: Ref<IReceipt | undefined> = ref();
    const provider: Ref<IProvider | undefined> = ref();
    const receiptFileData: Ref<string | undefined> = ref();

    const route = useRoute();
    const receiptID = route.params.receiptID as string;

    Promise.all([
      ReceiptRequest.Get(receiptID).then((r) => {
        receipt.value = r;
      }),
      ReceiptRequest.GetFile(receiptID)
        .then((r) => {
          receiptFileData.value = URL.createObjectURL(r);
          console.log(receiptFileData.value);
        })
        .catch((e: Error) => {
          console.log(e.message);
        }),
    ]).then(() => (loading.value = false));

    return {
      loading,
      receipt,
      receiptFileData,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading message="Loading Receipt" v-if="loading" />
    <div v-else-if="!receipt">
      <div class="alert alert-danger">
        <h4>Receipt Failed to load. Check Console</h4>
      </div>
    </div>
    <div v-else>
      <h3>Viewing {{ receipt.label }} Receipt</h3>
      <hr />
      <div class="row">
        <div class="col">
          <div class="card">
            <div class="card-header">Receipt Info</div>
            <div class="list-group list-group-flush">
              <div class="list-group-item d-flex justify-content-between">
                <span>Amount Paid</span>
                <span>{{ '$' + receipt.amount_paid }}</span>
              </div>
              <div class="list-group-item d-flex justify-content-between">
                <span>Date Paid</span>
                <span>{{ formatDate(receipt.date_paid) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="row mt-4" v-if="receiptFileData">
        <div class="col">
          <div class="card">
            <div class="card-header">Receipt File</div>
            <div class="card-body">
              <iframe :src="receiptFileData" width="100%" height="600" frameborder="0"></iframe>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
