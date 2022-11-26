<template>
    <div class="row">
        <div class="col">
            <div class="d-flex flex-row justify-content-center mt-5">
                <div class="spinner-border" style="width: 6rem; height: 6rem" role="status" />
                <br />
            </div>
            <div class="d-flex flex-row justify-content-center mt-2">
                <h4>{{ message ? `${message}...` : 'Loading...' }}</h4>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
defineProps({
    message: String
})
</script>

<style>

</style>
