import Vue from 'vue';
import Vuex from 'vuex';

import ProfileModule from '@/store/modules/profile';
import ConversationsModule from '@/store/modules/conversations';
import UsersModule from '@/store/modules/users';

Vue.use(Vuex);

export default new Vuex.Store({
    modules: {
        profile: ProfileModule,
        conversations: ConversationsModule,
        users: UsersModule,
    },
});
