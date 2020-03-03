<template lang="pug">
    q-page(padding v-scroll="scrolled")
        div(v-if="showStartChat")
            .q-caption.text-center All messages are loaded!
        div(v-else)
            div(v-if="showLoadMore")
                q-btn(
                label="Load more!"
                style="display: block; margin: 0 auto;"
                color="primary"
                size="md"
                @click="loadOldMessages()"
                )
        q-chat-message(
        v-for="message in conversation.messages"
        :key="`msg-${message.id}`"
        :sent="isMessageSent(message)"
        :text-color="getMessageTextColor(message)"
        :bg-color="getMessageBackgroundColor(message)"
        :name="getMessageName(message)"
        :avatar="getMessageAvatar(message)"
        :text="[message.content]"
        :stamp="getMessageStamp(message)"
        )
        q-page-sticky.bg-white.bg-grey-2(expand position='bottom')
            q-input.full-width.q-pa-md(
            v-model='message'
            :after='[{icon: "send", handler: sendMessage }]'
            placeholder='Type a message...'
            @keyup.enter='sendMessage()'
            autofocus='true'
            )
        #end-chat-marker.col(style='height: 60px;')
</template>

<script lang="ts">
import {Vue, Component} from 'vue-property-decorator';
import {Route, Next} from 'vue-router';
import {WebSocketRequest} from '@/models/websocket';
import Conversation from '@/models/conversation';
import store from '@/store';
import Message from '@/models/message';
import User from '@/models/user';
import moment from 'moment';
import {scroll, debounce} from 'quasar';
import {newMessageNotifier} from '@/websocket';

const {getScrollTarget, setScrollPosition} = scroll;

@Component
export default class ChatView extends Vue {
    public message = '';
    public conversation = {} as Conversation;
    public showLoadMore = false;
    public showStartChat = false;

    private scrolled = debounce((position: number) => {
        if (position === 0) {
            this.setLoadMore();
        }
    }, 200);

    public created() {
        newMessageNotifier.on('new_message', (message: Message) => {
            if (message.conversationId === this.conversation.id) {
                this.scrollToEndChat(500);
            }
        });
    }

    public async beforeRouteEnter(to: Route, from: Route, next: Next<ChatView>) {
        const conversationId = parseInt(to.params['conversation_id'], 10);
        await store.dispatch('conversations/loadOldMessages', {conversationId});
        next((vm) => vm.setConversation(conversationId));
    }

    public async beforeRouteUpdate(to: Route, from: Route, next: Next<ChatView>) {
        const conversationId = parseInt(to.params['conversation_id'], 10);
        await store.dispatch('conversations/loadOldMessages', {conversationId});
        this.setConversation(conversationId);
        next();
    }

    private scrollToEndChat(duration: number) {
        const el = document.getElementById('end-chat-marker')!;
        const target = getScrollTarget(el);
        setScrollPosition(target, el.offsetTop - el.scrollHeight, duration);
    }

    private setLoadMore() {
        if (this.conversation.firstMessage) {
            this.showStartChat = true;
            this.showLoadMore = false;
            return;
        }
        this.showLoadMore = true;
    }

    private setConversation(conversationId: number) {
        const position = store.state.conversations.positions.get(conversationId);
        this.conversation = store.state.conversations.conversations[position!];
        this.showStartChat = false;
        this.showLoadMore = false;
    }

    private isMessageSent(message: Message): boolean {
        return message.userId === store.state.profile.profile.id;
    }

    private loadOldMessages() {
        if (!this.conversation.firstMessage) {
            if (this.conversation.messages.length > 0) {
                store.dispatch('conversations/loadOldMessages', {
                    conversationId: this.conversation.id,
                    beforeId: this.conversation.messages[0].id,
                });
            }
        } else {
            this.showLoadMore = false;
            this.showStartChat = true;
        }
    }

    private getMessageTextColor(message: Message): string {
        if (this.isMessageSent(message)) {
            return 'white';
        }
        return 'black';
    }

    private getMessageBackgroundColor(message: Message): string {
        if (this.isMessageSent(message)) {
            return 'primary';
        }
        return 'grey-2';
    }

    private getUserById(userId: number): User {
        if (userId === store.state.profile.profile.id) {
            return store.state.profile.profile;
        } else {
            const pos = store.state.users.position.get(userId)!;
            return store.state.users.users[pos];
        }
    }

    private getMessageName(message: Message): string {
        const user = this.getUserById(message.userId);
        return `${user.givenName} ${user.familyName}`;
    }

    private getMessageAvatar(message: Message): string {
        const user = this.getUserById(message.userId);
        return user.profilePicture;
    }

    private getMessageStamp(message: Message): string {
        return moment(message.createdAt).format('MMMM Do YYYY, h:mm:ss a');
    }

    private sendMessage() {
        if (this.message === '') {
            return;
        }
        const request: WebSocketRequest = {
            method: 'new_message',
            data: {
                content: this.message,
                conversationId: this.conversation.id,
            },
        };
        this.$webSocket.send(JSON.stringify(request));
        this.message = '';
    }
}
</script>
