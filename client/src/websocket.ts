import Vue from 'vue';
import config from '@/config';
import {WebSocketResponse} from '@/models/websocket';
import Message from '@/models/message';
import store from '@/store';
import EventEmitter from 'eventemitter3';
import User from '@/models/user';
import Conversation from '@/models/conversation';

let newMessageNotifier: EventEmitter;

function attemptConnect(webSocket: WebSocket) {
    return new Promise((resolve) => {
        setTimeout(() => {
            resolve();
        }, 10);
    });
}

function handleMessage(ev: MessageEvent) {
    const res = JSON.parse(ev.data) as WebSocketResponse;
    switch (res.method) {
        case 'new_message':
            const message = res.data as Message;
            store.dispatch('conversations/addMessage', message);
            newMessageNotifier.emit('new_message', message);
            break;
        case 'new_friend':
            const friend = res.data as User;
            store.commit('users/SET_USER', friend);
            break;
        case 'new_conversation':
            const conversation = res.data as Conversation;
            conversation.firstLoad = true;
            store.commit('conversations/SET_CONVERSATION', conversation);
            store.commit('conversations/SET_USER_CONVERSATIONS', conversation);
            break;
    }
}

async function openWebSocket(uri: string) {
    const webSocket = new WebSocket(uri);
    while (true) {
        if (webSocket.readyState === 1 || webSocket.readyState === 3) {
            break;
        }
        await attemptConnect(webSocket);
    }
    webSocket.onmessage = handleMessage;
    webSocket.onclose = async (event: CloseEvent) => {
        console.log('WebSocket is disconnected, attempt to reconnect.');
        Vue.prototype.$webSocket = await openWebSocket(config.api.webSocketUri);
    };
    return webSocket;
}

async function registerWebSocket() {
    if (Vue.prototype.$webSocket === undefined || Vue.prototype.$webSocket === null) {
        newMessageNotifier = new EventEmitter();
        Vue.prototype.$webSocket = await openWebSocket(config.api.webSocketUri);
    }
}

export {registerWebSocket, newMessageNotifier};


