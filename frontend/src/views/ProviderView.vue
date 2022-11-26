<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import type { IProvider } from '../api/types/provider';
import type { IBill } from '../api/types/bill';
import { ProviderRequest } from '../api/index';
import { ref } from 'vue';
import type { Ref } from 'vue';
import { useRoute } from 'vue-router';

export default {
    components: {
        Loading,
        DashboardLayout,
    },
    setup() {
        const loading: Ref<Boolean> = ref(true);
        const provider: Ref<IProvider | undefined> = ref(undefined);
        const bills: Ref<IBill[]> = ref([]);
        const route = useRoute();
        const providerID = route.params.providerID as string;

        Promise.all([
            ProviderRequest.Get(providerID).then((r) => {
                provider.value = r;
            }),
            ProviderRequest.GetProviderBills(providerID).then((r) => {
                bills.value = r;
            }),
        ]).then(() => (loading.value = false));

        return {
            loading,
            provider,
            bills,
        };
    },
};
</script>

<template>
    <DashboardLayout>
        <Loading message="Loading Provider" v-if="loading" />
        <div v-else-if="!provider">
            <div class="alert alert-danger">
                <h4>Provider Failed to Load. Check Console</h4>
            </div>
        </div>

        <div v-else>
            <h3>{{ provider.name }}</h3>
            <hr />
            <h4>Bills Assigned This Provider</h4>
            <hr />
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
    </DashboardLayout>
</template>
