<template lang="pug">
    q-page(padding)
        div.q-ma-lg
            q-search(
            v-model="searchUser"
            autofocus
            placeholder="Search for a new friend!"
            clearable
            style
            )
        q-list(link inset-seperator highlight no-border)
            q-item(
            v-for="user in newFriends"
            :key="user.id"
            @click.native="confirmAddFriend(user)"
            )
                q-item-side(:avatar="user.profilePicture")
                q-item-main(:label="`${user.givenName} ${user.familyName}`")
        q-list(link inset-seperator highlight no-border)
            q-list-header Your Friends
            q-item(
            v-for="user in $store.state.users.users"
            :key="user.id"
            @click.native="$router.push({name: 'chat', params: {conversation_id: getConversationId(user)}})"
            )
                q-item-side(:avatar="user.profilePicture")
                q-item-main(:label="`${user.givenName} ${user.familyName}`")

</template>

<script lang="ts">
import {Vue, Watch, Component} from 'vue-property-decorator';
import User from '@/models/user';
import store from '@/store';
import apiService from '@/services/api';

@Component
export default class HomeView extends Vue {
    public searchUser = '';
    public newFriends: User[] = [];


    private getConversationId(user: User): number {
        return store.state.conversations.userConversations.get(user.id)!;
    }

    @Watch('searchUser')
    private async onSearchUserChanged(val: string, oldVal: string) {
        if (val.length < 3) {
            return;
        }
        this.newFriends = await apiService.searchFriends(val);
    }

    private confirmAddFriend(user: User) {
        this.$q.dialog({
            title: 'Confirm',
            message: `Do you want to add ${user.givenName} ${user.familyName} as your new friend?`,
            ok: true,
            cancel: true,
        }).then(async () => {
            const conversation = await apiService.addFriend(user.id);
            store.commit('users/SET_USER', user);
            store.commit('conversations/SET_CONVERSATION', conversation);
            this.$router.push({name: 'chat', params: {conversation_id: String(conversation.id)}});
        }).catch(() => {
        });
    }
}
</script>
