import auth0 from 'auth0-js';

import config from '@/config';
import router from '@/router';
import {init} from '@/init';

import apiService from '@/services/api';

const LOCAL_STORAGE_AUTH_EXPIRES_AT = 'auth_expires_at';

function authenticated(): boolean {
    const value = localStorage.getItem(LOCAL_STORAGE_AUTH_EXPIRES_AT);
    if (value !== null) {
        const expiresAt: number = JSON.parse(value);
        return new Date().getTime() < expiresAt;
    }
    return false;
}

async function logout() {
    await apiService.logout();
    await removeSession();
    router.replace({name: 'login'});
}

async function setSession(result: auth0.Auth0DecodedHash) {
    await localStorage.setItem(
        LOCAL_STORAGE_AUTH_EXPIRES_AT,
        JSON.stringify(new Date().getTime() + result.expiresIn! * 1000),
    );
}

async function removeSession() {
    await localStorage.removeItem(LOCAL_STORAGE_AUTH_EXPIRES_AT);
}

class AuthService {
    private readonly webAuth: auth0.WebAuth;

    constructor(domain: string,
                clientId: string,
                redirectUri: string,
                audience: string,
                responseType: string,
                scope: string) {
        this.webAuth = new auth0.WebAuth({
            domain,
            clientID: clientId,
            redirectUri,
            audience,
            responseType,
            scope,
        });
    }

    public show() {
        this.webAuth.authorize();
    }

    public handleCallback() {
        this.webAuth.parseHash(async (err, result) => {
            if (result && result.accessToken) {
                await apiService.authorize(result.accessToken);
                await setSession(result);
                await init();
                router.replace({name: 'home'});
                return;
            }
            router.replace({name: 'login'});
        });
    }
}

const service = new AuthService(
    config.auth0.domain,
    config.auth0.clientId,
    config.auth0.redirectUri,
    config.auth0.audience,
    config.auth0.responseType,
    config.auth0.scope,
);

export default service;
export {authenticated, logout};
