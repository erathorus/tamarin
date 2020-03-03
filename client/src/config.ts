export default {
    auth0: {
        domain: 'lattetalk.auth0.com',
        clientId: 'AQfe8ws5Ijr1JYcl0YGzwc1pUEoz5PgU',
        redirectUri: 'http://lattetalk.com/auth/callback',
        audience: 'https://auth0.lattetalk.com/api',
        responseType: 'token id_token',
        scope: 'openid profile',
    },
    api: {
        uri: 'http://lattetalk.com/api',
        webSocketUri: 'ws://lattetalk.com/api/ws',
    },
};
