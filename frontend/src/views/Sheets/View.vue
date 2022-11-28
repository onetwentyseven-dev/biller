<script lang="ts">
import { ref } from 'vue';
import type { Ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import CreateSheetEntryCard from '@/components/CreateSheetEntryCard.vue';
import SheetEntriesCard from '@/components/SheetEntriesCard.vue';

import { SheetRequest, EntryRequest, BillRequest } from '@/api/index';
import type {
  IBill,
  IBillSheet,
  IBillSheetEntry,
  ICreateUpdateBillSheetEntry,
} from '@/api/types/bill';

export default {
  components: {
    Loading,
    DashboardLayout,
    CreateSheetEntryCard,
    SheetEntriesCard,
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const loadingMessage: Ref<String> = ref('Loading Bill Sheet');
    const isEditing: Ref<Boolean> = ref(false);
    const isAddingEntry: Ref<boolean> = ref(false);
    const modifiableEntry: Ref<ICreateUpdateBillSheetEntry> = ref({
      bill_id: '',
      amount_due: 0,
      date_due: new Date(),
    });
    const sheet: Ref<IBillSheet | undefined> = ref(undefined);
    const entries: Ref<IBillSheetEntry[]> = ref([]);

    const bills: Ref<IBill[]> = ref([]);

    const router = useRouter();
    const route = useRoute();
    const sheetID = route.params.sheetID as string;

    const createEntry = async (entry: ICreateUpdateBillSheetEntry) => {
      loading.value = true;
      try {
        await EntryRequest.CreateBySheetID(sheetID, entry);
        await EntryRequest.ListBySheetID(sheetID).then((r) => {
          entries.value = r;
          isAddingEntry.value = false;
          loading.value = false;
        });
      } catch (e) {
        console.log(e);
      }
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
      entries,
      modifiableEntry,
      createEntry,
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
          <CreateSheetEntryCard
            v-if="isAddingEntry"
            :bills="bills"
            @cancel-add="() => (isAddingEntry = false)"
            @add-entry="createEntry"
          />
          <SheetEntriesCard v-else :entries="entries" @add-entry="() => (isAddingEntry = true)" />
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
