import Vue from 'vue';
import '@/class-component-hook';
import App from '@/App.vue';
import router from '@/router';
import store from '@/store';
import {init} from '@/init';

import '@/styles/quasar.styl';
import 'quasar-framework/dist/quasar.ie.polyfills';
import 'quasar-extras/roboto-font';
import 'quasar-extras/material-icons';
import Quasar from 'quasar';

Vue.use(Quasar, {
    config: {},
});

Vue.config.productionTip = false;

async function main() {
    await init();

    new Vue({
        router,
        store,
        render: (h) => h(App),
    }).$mount('#app');
}

main();
