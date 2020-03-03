<template lang="pug">
    q-page(padding)
        .q-display-4.text-center.q-ma-xl
            span.text-grey-5 Latte
            span.text-orange-5 Talk
        .q-headline.text-center.q-ma-lg.text-tertiary.text-weight-light Chat me chat you!
        div
            q-btn(
            label="Login"
            style="display: block; margin: 0 auto;"
            color="primary"
            size="lg"
            @click="show")
</template>

<script lang="ts">
import Vue from 'vue';
import authService from '@/services/auth';

export default Vue.extend({
    methods: {
        show: () => {
            authService.show();
        },
    },
});
</script>
