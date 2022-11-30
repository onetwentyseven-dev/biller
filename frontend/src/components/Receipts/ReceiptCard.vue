<template>
  <div class="card">
    <div class="card-header">Receipt Info</div>
    <div class="list-group list-group-flush">
      <div class="list-group-item d-flex justify-content-between" v-if="provider">
        <span>Provider</span>
        <span>{{ provider.name }}</span>
      </div>
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
</template>

<script lang="ts">
import type { IReceipt } from '@/api/types/receipt';
import type { IProvider } from '@/api/types/provider';
import { defineComponent, type PropType } from 'vue';

export default defineComponent({
  props: {
    receipt: {
      type: Object as PropType<IReceipt>,
      required: true,
    },
    provider: {
      type: Object as PropType<IProvider | undefined>,
    },
  },
  methods: {
    formatDate(v: string): string {
      return new Date(v).toDateString();
    },
  },
});
</script>
