<script lang="ts">
import { ref } from 'vue';
import type { Ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';

import { SheetRequest, EntryRequest } from '@/api/index';
import type { IBillSheet, IBillSheetEntry } from '@/api/types/bill';

export default {
  components: {
    Loading,
    DashboardLayout,
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const loadingMessage: Ref<String> = ref('Loading Bill Sheet');
    const isEditing: Ref<Boolean> = ref(false);
    const sheet: Ref<IBillSheet | undefined> = ref(undefined);
    const entries: Ref<IBillSheetEntry[]> = ref([]);

    const router = useRouter();
    const route = useRoute();
    const sheetID = route.params.sheetID as string;

    const updateSheet = () => {
      if (!sheet.value) return;
      loading.value = true;
      SheetRequest.Update(sheet.value.id, sheet.value).then(() => {
        isEditing.value = false;
        loading.value = false;
      });
    };

    const deleteSheet = () => {
      if (!sheet.value) return;
      loading.value = true;
      SheetRequest.Delete(sheet.value.id).then(() => {
        router.push({ name: 'sheets' });
      });
    };

    Promise.all([
      SheetRequest.Get(sheetID).then((r) => {
        sheet.value = r;
      }),
      EntryRequest.ListBySheetID(sheetID).then((r) => {
        entries.value = r;
      }),
    ]).then(() => {
      loading.value = false;
    });

    return {
      loading,
      loadingMessage,
      isEditing,
      sheet,
      entries,
      updateSheet,
      deleteSheet,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading :message="loadingMessage.toString()" v-if="loading" />
    <div v-else-if="!sheet">
      <div class="alert alert-danger">
        <h4>Sheet Failed to Load. Check Console</h4>
      </div>
    </div>

    <div v-else>
      <h3>{{ sheet.name }} Bill Sheet</h3>
      <hr />
      <div class="row">
        <div class="col-lg-8 offset-2">
          <div class="card">
            <div class="card-header d-flex justify-content-between">
              <span> Bill Sheet Details </span>
              <div>
                <div class="btn btn-sm btn-info me-1" @click="isEditing = !isEditing">
                  <font-awesome-icon icon="fa-solid fa-pencil" />
                </div>
                <div class="btn btn-sm btn-danger" @click="deleteSheet">
                  <font-awesome-icon icon="fa-solid fa-trash" />
                </div>
              </div>
            </div>
            <div v-if="isEditing">
              <div class="card-body">
                <div class="mb-3">
                  <label for="name">Provider Name</label>
                  <input type="text" v-model="sheet.name" class="form-control" />
                </div>
              </div>
              <div class="card-footer">
                <button class="btn btn-primary" @click="updateSheet">Update Sheet</button>
              </div>
            </div>
            <div v-else>
              <div class="list-group list-group-flush">
                <div class="list-group-item d-flex justify-content-between">
                  <span>Amount Paid</span>
                  <span>{{ sheet?.amount_paid || '$0.00' }}</span>
                </div>
                <div class="list-group-item d-flex justify-content-between">
                  <span>Amount Due</span>
                  <span>{{ sheet?.amount_due || '$0.00' }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="row mt-3">
        <div class="col">
          <div class="card">
            <div class="card-header d-flex justify-content-between">
              <span>Sheet Entries</span>
            </div>
            <table class="table table-bordered mb-0">
              <tbody v-if="entries.length">
                <tr v-for="entry in entries" :key="entry.entry_id">
                  <td>{{ entry.bill_name }}</td>
                  <td>{{ entry.amount_due }}</td>
                  <td>{{ entry.date_due.toDateString() }}</td>
                </tr>
              </tbody>
            </table>

            <!-- <div class="list-group flush">
              <RouterLink
                class="list-group-item list-group-item-action"
                :to="{ name: 'bill', params: { billID: bill.id } }"
                v-for="bill in bills"
                :key="bill.id"
              >
                {{ bill.name }}
              </RouterLink>
            </div> -->
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
