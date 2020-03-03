import axios from 'axios';
import config from '@/config';
import Conversation from '@/models/conversation';
import User from '@/models/user';
import Message from '@/models/message';

// API Endpoints
const ENDPOINT_AUTHORIZE = '/authorize';
const ENDPOINT_LOGOUT = '/logout';
const ENDPOINT_GET_PROFILE = '/profile';
const ENDPOINT_GET_FRIENDS = '/friends';
const ENDPOINT_GET_CONVERSATIONS = '/conversations';

// TODO: Better url join.
class ApiService {
    private readonly uri: string;

    constructor(uri: string) {
        this.uri = uri;
    }

    public async authorize(accessToken: string) {
        await this.request('post', ENDPOINT_AUTHORIZE, JSON.stringify({
            accessToken,
        }));
    }

    public async logout() {
        await this.request('post', ENDPOINT_LOGOUT);
    }

    public async getConversations() {
        const r = await this.request('get', ENDPOINT_GET_CONVERSATIONS);
        return r.data as Conversation[];
    }

    public async searchFriends(name: string) {
        const params = {name};
        const r = await this.request('get', `${ENDPOINT_GET_FRIENDS}/search`, '', params);
        return r.data as User[];
    }

    public async getConversationMessages(id: number, before?: number) {
        let params = null;
        if (before !== undefined && before !== null) {
            params = {before};
        }
        const r = await this.request('get', `${ENDPOINT_GET_CONVERSATIONS}/${id}/messages`, '', params);
        return r.data as Message[];
    }

    public async getProfile() {
        const r = await this.request('get', ENDPOINT_GET_PROFILE);
        return r.data as User;
    }

    public async getFriends() {
        const r = await this.request('get', ENDPOINT_GET_FRIENDS);
        return r.data as User[];
    }

    public async addFriend(userId: number) {
        const r = await this.request('post', `${ENDPOINT_GET_FRIENDS}/add/${userId}`);
        return r.data as Conversation;
    }

    private async request(method: string, endpoint: string, data = '', params?: any) {
        return await axios(this.uri + endpoint, {
            withCredentials: true,
            method,
            data,
            params,
        });
    }
}

const service = new ApiService(config.api.uri);

export default service;
