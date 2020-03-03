import {Module, VuexModule, Mutation, Action} from 'vuex-module-decorators';

import User from '@/models/user';

import apiService from '@/services/api';

// Mutation types
const SET_USER = 'SET_USER';
const SET_DATA_FETCHED = 'SET_DATA_FETCHED';

@Module({namespaced: true})
export default class UsersModule extends VuexModule {
    public users: User[] = [];
    public position: Map<number, number> = new Map();
    public dataFetched = false;

    @Mutation
    public [SET_USER](user: User) {
        this.position.set(user.id, this.users.length);
        this.users.push(user);
    }

    @Mutation
    public [SET_DATA_FETCHED]() {
        this.dataFetched = true;
    }

    @Action
    public async fetchData() {
        if (this.dataFetched) {
            return;
        }
        this.context.commit(SET_DATA_FETCHED);
        const friends = await apiService.getFriends();
        for (const user of friends) {
            this.context.commit(SET_USER, user);
        }
    }
}
