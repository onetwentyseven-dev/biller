<script lang="ts">
import { ref } from 'vue';
import type { Ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import numeral from 'numeral';

import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import CreateSheetEntryCard from '@/components/Sheets/CreateSheetEntryCard.vue';
import ModifyEntryCard from '@/components/Sheets/ModifyEntryCard.vue';
import SheetEntriesCard from '@/components/Sheets/SheetEntriesCard.vue';

import { SheetRequest, EntryRequest, BillRequest, ReceiptRequest } from '@/api/index';
import type { IBill, IBillSheet, IBillSheetEntry } from '@/api/types/bill';
import type { IReceipt } from '@/api/types/receipt';

export default {
  components: {
    Loading,
    DashboardLayout,
    CreateSheetEntryCard,
    ModifyEntryCard,
    SheetEntriesCard,
  },
  methods: {
    formatAmount(amount: number): string {
      return numeral(amount).format('$0,0.00');
    },
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const loadingMessage = 'Loading Bill Sheet';
    const isEditing: Ref<Boolean> = ref(false);
    const isAddingEntry: Ref<boolean> = ref(false);
    const modifiableEntry: Ref<IBillSheetEntry | undefined> = ref();
    const sheet: Ref<IBillSheet | undefined> = ref(undefined);
    const entries: Ref<IBillSheetEntry[]> = ref([]);
    const receipts: Ref<IReceipt[]> = ref([]);

    const bills: Ref<IBill[]> = ref([]);

    const router = useRouter();
    const route = useRoute();
    const sheetID = route.params.sheetID as string;

    const modifyEntry = (entryID: string) => {
      const entry = entries.value.find((e) => e.entry_id === entryID);
      if (!entry) return;
      console.log(entry);
      modifiableEntry.value = entry;
    };

    const createEntry = async (entry: IBillSheetEntry) => {
      loading.value = true;
      try {
        await EntryRequest.CreateBySheetID(sheetID, entry);
        await EntryRequest.ListBySheetID(sheetID).then((r) => {
          entries.value = r;
          isAddingEntry.value = false;
          loading.value = false;
        });
        await SheetRequest.Get(sheetID).then((r) => {
          sheet.value = r;
        });
      } catch (e) {
        console.log(e);
      }
    };

    const updateEntry = async (entry: IBillSheetEntry) => {
      loading.value = true;
      try {
        await EntryRequest.UpdateBySheetID(sheetID, entry.entry_id, entry);
        await EntryRequest.ListBySheetID(sheetID).then((r) => {
          entries.value = r;
        });
        await SheetRequest.Get(sheetID).then((r) => {
          sheet.value = r;
        });
      } catch (e) {
        console.log(e);
      }
      modifiableEntry.value = undefined;
      loading.value = false;
    };

    const deleteEntry = async (entryID: string) => {
      loading.value = true;
      try {
        await EntryRequest.DeleteByEntryID(sheetID, entryID);
        await EntryRequest.ListBySheetID(sheetID).then((r) => {
          entries.value = r;
        });
        await SheetRequest.Get(sheetID).then((r) => {
          sheet.value = r;
        });
      } catch (e) {
        console.log(e);
      }
      loading.value = false;
    };

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
      BillRequest.List().then((r) => {
        bills.value = r;
      }),
      SheetRequest.Get(sheetID).then((r) => {
        sheet.value = r;
      }),
      EntryRequest.ListBySheetID(sheetID).then((r) => {
        entries.value = r;
      }),
      ReceiptRequest.List().then((r) => {
        receipts.value = r;
      }),
    ]).then(() => {
      loading.value = false;
    });

    return {
      loading,
      loadingMessage,
      isEditing,
      isAddingEntry,
      sheet,
      bills,
      receipts,
      entries,
      modifiableEntry,
      createEntry,
      modifyEntry,
      updateEntry,
      deleteEntry,
      updateSheet,
      deleteSheet,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading :message="loadingMessage" v-if="loading" />
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
                  <label for="name">Sheet Name</label>
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
                  <span>Amount Due</span>
                  <span>{{ sheet.amount_due ? formatAmount(sheet.amount_due) : '$0.00' }}</span>
                </div>
                <div class="list-group-item d-flex justify-content-between">
                  <span>Amount Paid</span>
                  <span>{{ sheet.amount_paid ? formatAmount(sheet.amount_paid) : '$0.00' }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="row mt-3">
        <div class="col">
          <CreateSheetEntryCard
            v-if="isAddingEntry"
            :bills="bills"
            @cancel-add="() => (isAddingEntry = false)"
            @add-entry="createEntry"
          />
          <ModifyEntryCard
            v-else-if="modifiableEntry"
            :entry="modifiableEntry"
            :receipts="receipts"
            @cancel-modify="() => (modifiableEntry = undefined)"
            @modify-entry="updateEntry"
          />
          <SheetEntriesCard
            v-else
            :entries="entries"
            :receipts="receipts"
            @add-entry="() => (isAddingEntry = true)"
            @modify-entry="modifyEntry"
            @delete-entry="deleteEntry"
          />
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
