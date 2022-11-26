<template>
  <div class="container my-5">
    <div class="row">
      <div class="col-lg-4">
        <div class="list-group">
          <RouterLink :to="{ name: 'bills' }" class="list-group-item list-group-item-action">
            Bills
          </RouterLink>
          <RouterLink :to="{ name: 'providers' }" class="list-group-item list-group-item-action">
            Providers
          </RouterLink>
        </div>
      </div>
      <div class="col-lg-8">
        <slot></slot>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts"></script>
