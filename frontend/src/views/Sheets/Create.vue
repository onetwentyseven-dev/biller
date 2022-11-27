<script lang="ts">
import { SheetRequest } from '@/api';
import type { ICreateUpdateBillSheet } from '@/api/types/bill';
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import DashboardLayout from '../layouts/DashboardLayout.vue';
import Loading from '@/components/Loading.vue';

export default {
  name: 'CreateBillSheet',
  setup() {
    const loading = ref<Boolean>(false);
    const loadingMessage = 'Creating Bill Sheet';
    const errorMessage = ref<String | undefined>();
    const sheet: ICreateUpdateBillSheet = reactive({ name: '' });
    const router = useRouter();

    const createSheet = async () => {
      if (!sheet) return;
      errorMessage.value = undefined;
      loading.value = true;

      await SheetRequest.Create(sheet).then((r) => {
        router.push({ name: 'sheet', params: { sheetID: r.id } });
      });
    };

    return {
      loading,
      loadingMessage,
      errorMessage,
      sheet,
      createSheet,
    };
  },
  components: { DashboardLayout, Loading },
};
</script>

<template>
  <DashboardLayout>
    <h3>Create A New Bill Sheet</h3>
    <hr />
    <Loading :message="loadingMessage" v-if="loading" />

    <div class="container" v-else>
      <div class="row" v-if="errorMessage">
        <div class="col">
          <div class="alert alert-danger">
            <h4>Error Encountered Creating Bill Sheet</h4>
            {{ errorMessage }}
          </div>
        </div>
      </div>
      <div class="row">
        <div class="col-lg-6 offset-3">
          <div class="card">
            <div class="card-header">Provision a New Bill Sheet</div>
            <div class="card-body">
              <div class="mb-3">
                <label for="name">Sheet Name</label>
                <input type="text" v-model="sheet.name" class="form-control" />
              </div>
            </div>
            <div class="card-footer d-grid gap-2">
              <button class="btn btn-primary" @click="createSheet">Create Bill Sheet</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
