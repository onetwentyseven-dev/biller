<script lang="ts">
import { useAuth0 } from '@auth0/auth0-vue';
import { defineComponent } from 'vue';

export default defineComponent({
  setup() {
    const { user, logout } = useAuth0();

    return {
      user,
      logout,
    };
  },
});
</script>

<template>
  <nav class="navbar navbar-expand-lg bg-light">
    <div class="container">
      <RouterLink :to="{ name: 'home' }" class="navbar-brand"> Biller </RouterLink>

      <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
        <li class="nav-item dropdown">
          <a
            class="nav-link dropdown-toggle"
            role="button"
            data-bs-toggle="dropdown"
            aria-expanded="false"
          >
            Welcome {{ user.given_name }}
          </a>
          <ul class="dropdown-menu dropdown-menu-end">
            <li><a class="dropdown-item" @click="() => logout()"> Logout </a></li>
          </ul>
        </li>
      </ul>
    </div>
  </nav>
</template>
