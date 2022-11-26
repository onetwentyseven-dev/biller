<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import { BillRequest } from '../api/index';
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
            <h3>Bills</h3>
            <hr />
            <div class="container">
                <div class="row">
                    <div class="col">
                        <div class="list-group">
                            <RouterLink :to="{}"> </RouterLink>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </DashboardLayout>
</template>
