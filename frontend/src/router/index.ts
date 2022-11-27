import { authGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';
import BillsList from '@/views/Bills/List.vue';
import BillView from '@/views/Bills/View.vue';
import BillNew from '@/views/Bills/Create.vue';
import ProvidersList from '@/views/Providers/List.vue';
import ProviderView from '@/views/Providers/View.vue';
import ProviderNew from '@/views/Providers/Create.vue';
import SheetsList from '@/views/Sheets/List.vue';
import CreateBillSheet from '@/views/Sheets/Create.vue';
import ViewBillSheet from '@/views/Sheets/View.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'bills',
      component: BillsList,
    },
    {
      path: '/bills/create',
      name: 'create-bill',
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
    {
      path: '/sheets',
      name: 'sheets',
      component: SheetsList,
    },
    {
      path: '/sheets/create',
      name: 'create-sheet',
      component: CreateBillSheet,
    },
    {
      path: '/sheets/:sheetID',
      name: 'sheet',
      component: ViewBillSheet,
    },
  ],
});

router.beforeEach(authGuard);

export default router;
