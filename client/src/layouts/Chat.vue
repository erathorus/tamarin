<template lang="pug">
    q-layout(view="hHh LpR lFf")
        q-layout-header
            q-toolbar
                q-btn.q-mx-sm(round dense @click="showConversations = !showConversations" icon="chat")
                q-toolbar-title Chat
                q-btn.q-mx-sm(round @click="$router.push({name: 'home'})" icon="people")
                q-btn.q-mx-sm(round :style="avatarStyle" @click="$router.replace({name: 'logout'})")

        q-layout-drawer(side="left" v-model="showConversations" content-class="bg-grey-2")
            q-list(no-border link inset-seperator highlight)
                q-list-header Conversations
                q-item(
                v-for="conversation in $store.state.conversations.conversations"
                :key="conversation.id"
                @click.native="$router.push({name: 'chat', params: {conversation_id: conversation.id}})"
                )
                    q-item-side(
                    v-if="getConversationAvatar(conversation) !== ''"
                    :avatar="getConversationAvatar(conversation)"
                    )
                    q-item-side(v-else inverted icon="account_circle")
                    q-item-main(
                    v-if="conversation.messages.length > 0"
                    :label="getConversationName(conversation)"
                    :sublabel="conversation.messages[conversation.messages.length - 1].content"
                    )
                    q-item-main(
                    v-else
                    :label="getConversationName(conversation)"
                    )
        q-page-container
            router-view
</template>

<script lang="ts">
import {Vue, Component} from 'vue-property-decorator';

import User from '@/models/user';
import Conversation from '@/models/conversation';

import store from '@/store';

@Component
export default class ChatLayout extends Vue {
    public showConversations = true;

    public avatarStyle = {
        backgroundImage: `url(${store.state.profile.profile.profilePicture})`,
        backgroundSize: 'contain',
    };

    private getUserName(user: User): string {
        return user.givenName + ' ' + user.familyName;
    }

    private getConversationAvatar(conversation: Conversation): string {
        for (const userId of conversation.userIds) {
            if (userId !== store.state.profile.profile.id) {
                const pos = store.state.users.position.get(userId)!;
                return store.state.users.users[pos].profilePicture;
            }
        }
        return '';
    }

    private getConversationName(conversation: Conversation): string {
        for (const userId of conversation.userIds) {
            if (userId !== store.state.profile.profile.id) {
                const pos = store.state.users.position.get(userId)!;
                return this.getUserName(store.state.users.users[pos]);
            }
        }
        return '';
    }
}
</script>
