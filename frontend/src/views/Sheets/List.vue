<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import { SheetRequest } from '@/api/index';
import { ref } from 'vue';
import type { Ref } from 'vue';
import type { IBillSheet } from '@/api/types/bill';

export default {
  name: 'SheetList',
  components: {
    Loading,
    DashboardLayout,
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const sheets: Ref<IBillSheet[]> = ref([]);

    SheetRequest.List().then((r) => {
      sheets.value = r;
      loading.value = false;
    });

    return {
      loading,
      sheets,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading message="Loading Sheets" v-if="loading" />
    <div v-else>
      <div class="d-flex justify-content-between align-items-center">
        <span class="h3">My Bill Sheets</span>
        <RouterLink :to="{ name: 'create-sheet' }" class="btn btn-outline-primary">
          <font-awesome-icon icon="fa-solid fa-plus" />
        </RouterLink>
      </div>
      <hr />
      <div class="container">
        <div class="row">
          <div class="col">
            <div class="list-group">
              <RouterLink
                :to="{ name: 'sheet', params: { sheetID: sheet.id } }"
                v-for="sheet in sheets"
                :key="sheet.id"
                class="list-group-item list-group-item-action"
              >
                {{ sheet.name }}
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
