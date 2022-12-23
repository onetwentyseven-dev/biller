<template>
  <div class="card">
    <div class="card-header">Modifying Entry {{ entry.bill_name }}</div>
    <div class="card-body">
      <div class="row">
        <div class="col-lg-6 offset-3">
          <div class="mb-3">
            <label for="receipt_id">Receipt:</label>
            <select id="receipt_id" class="form-select" v-model="entry.receipt_id">
              <option v-for="receipt in receipts" :value="receipt.id">
                {{ receipt.label }}
              </option>
            </select>
          </div>
          <div class="mb-3">
            <label for="amount_paid">Amount Paid:</label>
            <input
              type="text"
              id="amount_paid"
              v-model.number="entry.amount_paid"
              class="form-control"
            />
          </div>
          <div class="mb-3">
            <label for="date_paid">Date Paid:</label>
            <Datepicker id="date_paid" v-model="entry.date_paid" class="form-control" />
          </div>
        </div>
      </div>
    </div>
    <div class="card-footer">
      <div class="row">
        <div class="col-lg-4 offset-4">
          <button class="btn btn-primary" @click="modifyEntry">Update Entry</button>
          <button class="btn btn-danger ms-2" @click="cancelModify">Nevermind</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Datepicker from 'vue3-datepicker';
import { defineComponent, type PropType, type Ref, ref } from 'vue';

import type { IBillSheetEntry } from '@/api/types/bill';
import type { IReceipt } from '@/api/types/receipt';

export default defineComponent({
  emits: ['modifyEntry', 'cancelModify'],
  methods: {
    modifyEntry() {
      this.$emit('modifyEntry', this.entry);
    },
    cancelModify() {
      this.$emit('cancelModify');
    },
  },
  components: {
    Datepicker,
  },
  props: {
    entry: {
      type: Object as PropType<IBillSheetEntry>,
      required: true,
    },
    receipts: {
      type: Array as PropType<IReceipt[]>,
      required: true,
    },
  },
});
</script>
