import store from '@/store';

async function initStore() {
    await store.dispatch('profile/fetchData');
    await store.dispatch('users/fetchData');
    await store.dispatch('conversations/fetchData');
}

export {initStore};
