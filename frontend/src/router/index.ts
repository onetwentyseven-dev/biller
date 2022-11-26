import { authGuard } from "@auth0/auth0-vue";
import { createRouter, createWebHistory } from "vue-router";
import BillsView from "../views/BillsView.vue";
import BillView from "../views/BillView.vue";
import ProvidersView from "@/views/ProvidersView.vue";
import ProviderView from "@/views/ProviderView.vue";


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "bills",
      component: BillsView,
      beforeEnter: authGuard
    },
    { 
      path: "/bills/:billID",
      name: "bill",
      component: BillView,
      beforeEnter: authGuard
    },
    {
      path: "/providers",
      name: "providers",
      component: ProvidersView,
      beforeEnter: authGuard
    },
    {
      path: "/providers/:providerID",
      name: "provider",
      component: ProviderView,
      beforeEnter: authGuard
    },

  ],
});

export default router;
