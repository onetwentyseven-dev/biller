<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import { BillRequest } from '@/api/index';
import { ref } from 'vue';
import type { Ref } from 'vue';
import type { IBill } from '@/api/types/bill';

export default {
  components: {
    Loading,
    DashboardLayout,
  },
  setup() {
    const loading: Ref<Boolean> = ref(true);
    const bills: Ref<IBill[]> = ref([]);

    BillRequest.List().then((r) => {
      bills.value = r;
      loading.value = false;
    });

    return {
      loading,
      bills,
    };
  },
};
</script>

<template>
  <DashboardLayout>
    <Loading message="Loading Bills" v-if="loading" />
    <div v-else>
      <div class="d-flex justify-content-between align-items-center">
        <span class="h3">Bills</span>
        <RouterLink :to="{ name: 'new-bill' }" class="btn btn-outline-primary">
          <font-awesome-icon icon="fa-solid fa-plus" />
        </RouterLink>
      </div>
      <hr />
      <div class="container">
        <div class="row">
          <div class="col">
            <div class="list-group">
              <RouterLink
                :to="{ name: 'bill', params: { billID: bill.id } }"
                v-for="bill in bills"
                :key="bill.id"
                class="list-group-item list-group-item-action"
              >
                {{ bill.name }}
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </DashboardLayout>
</template>
