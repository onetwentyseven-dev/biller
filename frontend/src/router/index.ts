import { authGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';
import BillsList from '@/views/Bills/List.vue';
import BillView from '@/views/Bills/View.vue';
import BillNew from '@/views/Bills/New.vue';
import ProvidersList from '@/views/Providers/List.vue';
import ProviderView from '@/views/Providers/View.vue';
import ProviderNew from '@/views/Providers/New.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'bills',
      component: BillsList,
    },
    {
      path: '/bills/new',
      name: 'new-bill',
      component: BillNew,
    },
    {
      path: '/bills/:billID',
      name: 'bill',
      component: BillView,
    },
    {
      path: '/providers',
      name: 'providers',
      component: ProvidersList,
    },
    {
      path: '/providers/create',
      name: 'create-provider',
      component: ProviderNew,
    },
    {
      path: '/providers/:providerID',
      name: 'provider',
      component: ProviderView,
    },
  ],
});

router.beforeEach(authGuard);

export default router;
