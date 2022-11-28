<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import { ReceiptRequest } from '@/api/index';
import { ref } from 'vue';
import type { Ref } from 'vue';
import type { IReceipt } from '@/api/types/receipt';

export default {
  name: 'ReceiptList',
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
    const receipts: Ref<IReceipt[]> = ref([]);

    ReceiptRequest.List().then((r) => {
      receipts.value = r;
      loading.value = false;
    });

    return {
      loading,
      receipts,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading message="Loading Receipts" v-if="loading" />
    <div v-else>
      <div class="d-flex justify-content-between align-items-center">
        <span class="h3">My Receipts</span>
        <RouterLink :to="{ name: 'create-receipt' }" class="btn btn-outline-primary">
          <font-awesome-icon icon="fa-solid fa-plus" />
        </RouterLink>
      </div>
      <hr />
      <div class="container">
        <div class="col">
          <div class="list-group">
            <RouterLink
              :to="{ name: 'receipt', params: { receiptID: receipt.id } }"
              v-for="receipt in receipts"
              :key="receipt.id"
              class="list-group-item list-group-item-action"
            >
              {{ formatDate(receipt.date_paid.toString()) }} - {{ receipt.label }}
            </RouterLink>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
