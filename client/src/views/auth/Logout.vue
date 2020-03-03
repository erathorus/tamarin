<template lang="pug">
    q-page
        .fixed-center
            q-spinner(:size="200" color="primary")
</template>

<script lang="ts">
import Vue from 'vue';
import {logout} from '@/services/auth';

export default Vue.extend({
    async created() {
        await logout();
    },
});
</script>
