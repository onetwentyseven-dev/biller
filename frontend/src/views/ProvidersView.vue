<script lang="ts">
import Loading from '@/components/Loading.vue';
import DashboardLayout from '@/views/layouts/DashboardLayout.vue';
import type { IProvider } from '../api/types/provider';
import { ProviderRequest } from '../api/index';
import { ref } from 'vue';
import type { Ref } from 'vue';

export default {
    components: {
        Loading,
        DashboardLayout,
    },
    setup() {
        const loading: Ref<Boolean> = ref(true);
        const providers: Ref<IProvider[]> = ref([]);

        ProviderRequest.List().then((r) => {
            providers.value = r;
            loading.value = false;
        });

        return {
            loading,
            providers,
        };
    },
};
</script>

<template>
    <DashboardLayout>
        <Loading message="Loading Providers" v-if="loading" />
        <div v-else>
            <h3>Providers</h3>
            <hr />
            <div class="container">
                <div class="row">
                    <div class="col">
                        <div class="list-group">
                            <RouterLink :to="{ name: 'provider', params: { providerID: provider.id } }" v-for="provider in providers" :key="provider.id" class="list-group-item list-group-item-action">
                                {{ provider.name }}
                            </RouterLink>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </DashboardLayout>
</template>
