import { authGuard } from '@auth0/auth0-vue';
import { createRouter, createWebHistory } from 'vue-router';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'bills',
      component: import('@/views/Bills/List.vue'),
    },
    {
      path: '/bills/create',
      name: 'create-bill',
      component: () => import('@/views/Bills/Create.vue'),
    },
    {
      path: '/bills/:billID',
      name: 'bill',
      component: () => import('@/views/Bills/View.vue'),
    },
    {
      path: '/providers',
      name: 'providers',
      component: () => import('@/views/Providers/List.vue'),
    },
    {
      path: '/providers/create',
      name: 'create-provider',
      component: () => import('@/views/Providers/Create.vue'),
    },
    {
      path: '/providers/:providerID',
      name: 'provider',
      component: () => import('@/views/Providers/View.vue'),
    },
    {
      path: '/sheets',
      name: 'sheets',
      component: () => import('@/views/Sheets/List.vue'),
    },
    {
      path: '/sheets/create',
      name: 'create-sheet',
      component: () => import('@/views/Sheets/Create.vue'),
    },
    {
      path: '/sheets/:sheetID',
      name: 'sheet',
      component: () => import('@/views/Sheets/View.vue'),
    },
    {
      path: '/receipts',
      name: 'receipts',
      component: () => import('@/views/Receipts/List.vue'),
    },
    {
      path: '/receipts/create',
      name: 'create-receipt',
      component: () => import('@/views/Receipts/Create.vue'),
    },
    {
      path: '/receipts/:receiptID',
      name: 'receipt',
      component: () => import('@/views/Receipts/View.vue'),
    },
  ],
});

router.beforeEach(authGuard);

export default router;
