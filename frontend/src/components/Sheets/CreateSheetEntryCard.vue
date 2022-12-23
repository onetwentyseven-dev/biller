<template>
  <div class="card">
    <div class="card-header">Add Entry To Sheet</div>
    <div class="card-body">
      <div class="row">
        <div class="col-lg-6 offset-3">
          <div class="mb-3">
            <label for="bill_id">Bill:</label>
            <select id="bill_id" class="form-select" v-model="entry.bill_id">
              <option v-for="bill in bills" :value="bill.id">
                {{ bill.name }} - {{ bill.provider_name }}
              </option>
            </select>
          </div>
          <div class="mb-3">
            <label for="amount_due">Amount Due:</label>
            <input
              type="text"
              id="amount_due"
              v-model.number="entry.amount_due"
              class="form-control"
            />
          </div>
          <div class="mb-3">
            <label for="date_due">Date Due:</label>
            <Datepicker v-model="entry.date_due" class="form-control" />
          </div>
        </div>
      </div>
    </div>
    <div class="card-footer">
      <div class="row">
        <div class="col-lg-2 offset-4 d-grid gap-2">
          <button class="btn btn-primary" @click="submitEntry">Create Entry</button>
        </div>
        <div class="col-lg-2 d-grid gap-2">
          <button class="btn btn-danger" @click="cancelAdd">Nevermind</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, type PropType, type Ref } from 'vue';
import Datepicker from 'vue3-datepicker';

import type { IBill, ICreateUpdateBillSheetEntry } from '@/api/types/bill';

export default defineComponent({
  components: {
    Datepicker,
  },
  emits: ['cancelAdd', 'addEntry'],
  props: {
    bills: {
      type: Array as PropType<IBill[]>,
      required: true,
    },
  },
  methods: {
    cancelAdd() {
      this.$emit('cancelAdd');
    },
    submitEntry() {
      this.$emit('addEntry', this.entry);
    },
  },
  setup() {
    const entry: Ref<ICreateUpdateBillSheetEntry> = ref({
      bill_id: '',
      amount_due: 0,
      date_due: new Date(),
    });

    return {
      entry,
    };
  },
});
</script>
