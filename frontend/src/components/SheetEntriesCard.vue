<template>
  <div class="card">
    <div class="card-header d-flex justify-content-between">
      <span>Sheet Entries</span>
    </div>
    <table class="table table-bordered mb-0">
      <tbody v-if="entries.length">
        <tr v-for="entry in entries" :key="entry.entry_id">
          <td>{{ entry.bill_name }}</td>
          <td>{{ entry.amount_due }}</td>
          <td>{{ entry.date_due }}</td>
          <td>{{ entry.amount_paid || 'N/A' }}</td>
          <td>{{ entry.date_paid || 'N/A' }}</td>
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
import type { IBillSheetEntry } from '@/api/types/bill';
import { defineComponent, type PropType } from 'vue';
export default defineComponent({
  emits: ['addEntry'],
  props: {
    entries: {
      type: Object as PropType<IBillSheetEntry[]>,
      required: true,
    },
  },
  methods: {
    addEntry() {
      this.$emit('addEntry');
    },
  },
});
</script>
