<template>
  <div class="card">
    <div class="card-header d-flex justify-content-between">
      <span>Sheet Entries</span>
    </div>
    <table class="table table-bordered mb-0">
      <thead>
        <tr>
          <th>Bill Name</th>
          <th>Bill Dues</th>
          <th>Bill Payment Details</th>
          <th></th>
        </tr>
      </thead>
      <tbody v-for="entry in entries" :key="entry.entry_id">
        <tr>
          <td rowspan="2">{{ entry.bill_name }}</td>
          <td>
            <strong>Due Date:</strong> {{ formatDate(entry.date_due) }}<br />
            <strong>Amount Due:</strong> {{ formatAmount(entry.amount_due) }}
          </td>
          <td v-if="entry.amount_paid && entry.date_paid">
            <strong>Paid Date:</strong> {{ formatDate(entry.date_paid) }}<br />
            <strong>Amount Paid:</strong> {{ formatAmount(entry.amount_paid) }}
          </td>
          <td v-else>
            <strong>N/A</strong>
          </td>
          <td>
            <div>
              <RouterLink
                :to="{ name: 'receipt', params: { receiptID: entry.receipt_id } }"
                v-if="entry.receipt_id"
              >
                <div class="btn btn-sm btn-info me-1">
                  <font-awesome-icon icon="fa-solid fa-right-to-bracket" />
                </div>
              </RouterLink>
              <div class="btn btn-sm btn-info me-1" @click="modifyEntry(entry.entry_id)">
                <font-awesome-icon icon="fa-solid fa-pencil" />
              </div>
              <div class="btn btn-sm btn-danger">
                <font-awesome-icon icon="fa-solid fa-trash" @click="deleteEntry(entry.entry_id)" />
              </div>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    <div class="card-footer">
      <div class="row">
        <div class="col-lg-4 offset-4 d-grid gap-2">
          <button class="btn btn-info" @click="addEntry">Add New Entry</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { DateTime } from 'luxon';
import numeral from 'numeral';
import { defineComponent, type PropType, type Ref, ref } from 'vue';
import type { IBillSheetEntry } from '@/api/types/bill';
export default defineComponent({
  components: {},
  emits: ['addEntry', 'modifyEntry', 'deleteEntry'],
  props: {
    entries: {
      type: Array as PropType<IBillSheetEntry[]>,
      required: true,
    },
  },
  methods: {
    addEntry() {
      this.$emit('addEntry');
    },
    modifyEntry(entryID: string) {
      this.$emit('modifyEntry', entryID);
    },
    deleteEntry(entryID: string) {
      this.$emit('deleteEntry', entryID);
    },
    formatDate(date: Date): string {
      return DateTime.fromJSDate(date).toLocaleString(DateTime.DATE_MED_WITH_WEEKDAY);
    },
    formatAmount(amount: number): string {
      return numeral(amount).format('$0,0.00');
    },
  },
});
</script>
