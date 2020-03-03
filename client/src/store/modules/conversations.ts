import {Module, VuexModule, Mutation, Action} from 'vuex-module-decorators';

import Conversation from '@/models/conversation';
import apiService from '@/services/api';
import Message from '@/models/message';

// Mutation types
const SET_CONVERSATION = 'SET_CONVERSATION';
const SET_DATA_FETCHED = 'SET_DATA_FETCHED';
const SET_CONVERSATION_FIRST_LOAD = 'SET_CONVERSATION_FIRST_LOAD';
const SET_CONVERSATION_FIRST_MESSAGE = 'SET_CONVERSATION_FIRST_MESSAGE';
const SET_USER_CONVERSATIONS = 'SET_USER_CONVERSATIONS';
const ADD_MESSAGE = 'ADD_MESSAGE';
const ADD_CONVERSATION_MESSAGES = 'ADD_CONVERSATION_MESSAGES';

@Module({namespaced: true})
export default class ConversationsModule extends VuexModule {
    public conversations: Conversation[] = [];
    public positions: Map<number, number> = new Map();
    public userConversations: Map<number, number> = new Map();
    public dataFetched = false;

    @Mutation
    public [SET_CONVERSATION](conversation: Conversation) {
        conversation.messages = [];
        conversation.firstMessage = false;
        this.positions.set(conversation.id, this.conversations.length);
        this.conversations.push(conversation);
    }

    @Mutation
    public [SET_DATA_FETCHED]() {
        this.dataFetched = true;
    }

    @Mutation
    public [SET_CONVERSATION_FIRST_LOAD](pos: number) {
        this.conversations[pos].firstLoad = true;
    }

    @Mutation
    public [SET_CONVERSATION_FIRST_MESSAGE](pos: number) {
        this.conversations[pos].firstMessage = true;
    }

    @Mutation
    public [ADD_MESSAGE](payload: any) {
        this.conversations[payload.pos].messages.push(payload.message);
    }

    @Mutation
    public [ADD_CONVERSATION_MESSAGES](payload: any) {
        const pos = payload.pos as number;
        const messages = payload.messages as Message[];
        messages.reverse();
        this.conversations[pos].messages = messages.concat(this.conversations[pos].messages);
    }

    @Mutation
    public [SET_USER_CONVERSATIONS](conversation: Conversation) {
        for (const userId of conversation.userIds) {
            this.userConversations.set(userId, conversation.id);
        }
    }

    @Action
    public async fetchData() {
        if (this.dataFetched) {
            return;
        }
        this.context.commit(SET_DATA_FETCHED);
        const conversations = await apiService.getConversations();
        for (const conversation of conversations) {
            conversation.firstLoad = false;
            this.context.commit(SET_CONVERSATION, conversation);
            this.context.commit(SET_USER_CONVERSATIONS, conversation);
        }
    }

    @Action
    public async loadOldMessages(payload: any) {
        const conversationId = payload.conversationId as number;
        const beforeId = payload.beforeId as number;
        const pos = this.positions.get(conversationId)!;
        const conversation = this.conversations[pos];
        if (conversation.firstMessage) {
            return;
        }
        let messages: Message[];
        if (beforeId === undefined || beforeId === null) {
            if (!conversation.firstLoad) {
                this.context.commit(SET_CONVERSATION_FIRST_LOAD, pos);
                messages = await apiService.getConversationMessages(conversationId);
                if (messages.length < 30) {
                    this.context.commit(SET_CONVERSATION_FIRST_MESSAGE, pos);
                }
                this.context.commit(ADD_CONVERSATION_MESSAGES, {pos, messages});
            }
            return;
        }
        messages = await apiService.getConversationMessages(conversationId, beforeId);
        if (messages.length === 0) {
            this.context.commit(SET_CONVERSATION_FIRST_MESSAGE, pos);
        } else {
            this.context.commit(ADD_CONVERSATION_MESSAGES, {pos, messages});
        }
    }

    @Action
    public async addMessage(message: Message) {
        const pos = this.positions.get(message.conversationId);
        this.context.commit(ADD_MESSAGE, {pos, message});
    }
}
