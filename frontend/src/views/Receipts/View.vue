<script lang="ts">
import { ref, type Ref } from 'vue';
import { useRoute } from 'vue-router';

import ReceiptCard from '@/components/Receipts/ReceiptCard.vue';
import ReceiptFileCard from '@/components/Receipts/ReceiptFileCard.vue';
import { ProviderRequest, ReceiptRequest } from '@/api';
import type { IReceipt } from '@/api/types/receipt';
import type { IProvider } from '@/api/types/provider';
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';

export default {
  components: {
    Loading,
    DashboardLayout,
    ReceiptCard,
    ReceiptFileCard,
  },
  methods: {
    formatDate(v: string): string {
      return new Date(v).toDateString();
    },
    uploadReceipt(id: string, file: any): void {
      this.loading = false;
    },
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const receipt: Ref<IReceipt | undefined> = ref();
    const provider: Ref<IProvider | undefined> = ref();

    const route = useRoute();
    const receiptID = route.params.receiptID as string;

    Promise.all([
      ReceiptRequest.Get(receiptID).then(async (r) => {
        receipt.value = r;
        if (r.provider_id) {
          await ProviderRequest.Get(r.provider_id).then((r) => {
            provider.value = r;
          });
        }
      }),
    ]).then(() => (loading.value = false));

    return {
      loading,
      receipt,
      provider,
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
          <ReceiptCard :receipt="receipt" :provider="provider" />
        </div>
      </div>
      <div class="row mt-4">
        <div class="col">
          <ReceiptFileCard :receipt="receipt" />
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
