import {Action, Module, Mutation, VuexModule} from 'vuex-module-decorators';

import User from '@/models/user';

import apiService from '@/services/api';

// Mutation types
const SET_PROFILE = 'SET_PROFILE';
const SET_DATA_FETCHED = 'SET_DATA_FETCHED';

@Module({namespaced: true})
export default class ProfileModule extends VuexModule {
    public profile: User = {
        id: 0,
        givenName: '',
        familyName: '',
        profilePicture: '',
    };
    public dataFetched = false;

    @Mutation
    public [SET_PROFILE](profile: User) {
        this.profile = profile;
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
        const profile = await apiService.getProfile();
        this.context.commit(SET_PROFILE, profile);
    }
}
