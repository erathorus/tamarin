import Vue from 'vue';
import Router, {Route, RawLocation} from 'vue-router';

import {authenticated} from '@/services/auth';

import AuthLayout from '@/layouts/Auth.vue';
import LoginView from '@/views/auth/Login.vue';
import CallbackView from '@/views/auth/Callback.vue';

import DefaultLayout from '@/layouts/Default.vue';
import LogoutView from '@/views/auth/Logout.vue';

import ChatLayout from '@/layouts/Chat.vue';
import HomeView from '@/views/Home.vue';
import ChatView from '@/views/Chat.vue';

Vue.use(Router);

const ensureAuthenticated = (to: Route, from: Route, next: (to?: RawLocation) => void) => {
    if (authenticated()) {
        next();
        return;
    }
    next({name: 'login'});
};

const ensureAnonymous = (to: Route, from: Route, next: (to?: RawLocation) => void) => {
    if (!authenticated()) {
        next();
        return;
    }
    next({name: 'home'});
};

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [
        {
            path: '/auth',
            component: AuthLayout,
            children: [
                {
                    path: 'login',
                    name: 'login',
                    beforeEnter: ensureAnonymous,
                    component: LoginView,
                },
                {
                    path: 'callback',
                    name: 'callback',
                    beforeEnter: ensureAnonymous,
                    component: CallbackView,
                },
                {
                    path: 'logout',
                    name: 'logout',
                    beforeEnter: ensureAuthenticated,
                    component: LogoutView,
                },
            ],
        },
        {
            path: '/',
            beforeEnter: ensureAuthenticated,
            component: ChatLayout,
            children: [
                {
                    path: '',
                    name: 'home',
                    component: HomeView,
                },
                {
                    path: 'c/:conversation_id',
                    name: 'chat',
                    component: ChatView,
                },
            ],
        },
    ],
});
