import Vue, {VNode} from 'vue';

declare global {
    namespace JSX {
        // tslint:disable no-empty-interface
        interface Element extends VNode {
        }

        // tslint:disable no-empty-interface
        interface ElementClass extends Vue {
        }

        interface IntrinsicElements {
            [elem: string]: any;
        }
    }
}

declare module 'vue-router' {
    type Next<T extends Vue = Vue> = (to?: (vm: T) => any) => void;
}

declare module 'vue/types/vue' {
    interface Vue {
        $webSocket: WebSocket;
        $q: any;
    }
}
