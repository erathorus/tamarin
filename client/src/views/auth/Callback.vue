<template lang="pug">
    q-page
        .fixed-center
            q-spinner(:size="200" color="primary")
</template>

<script lang="ts">
import Vue from 'vue';
import authService from '@/services/auth';

export default Vue.extend({
    created() {
        authService.handleCallback();
    },
});
</script>
