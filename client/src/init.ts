import {authenticated} from '@/services/auth';
import {initStore} from '@/store/init';
import {registerWebSocket} from '@/websocket';

async function init() {
    if (authenticated()) {
        await initStore();
        await registerWebSocket();
    }
}

export {init};
