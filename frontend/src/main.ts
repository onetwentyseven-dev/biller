import { createApp } from "vue";
import { createPinia } from "pinia";
import { createAuth0 } from '@auth0/auth0-vue';

import App from "./App.vue";
import router from "./router";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(
    createAuth0({
        domain: "onetwentyseven.us.auth0.com",
        client_id: "u2hgEu2s28xKcKWw0JRgqlgcA6hRLatk",
        redirect_uri: window.location.origin
    })
);

app.mount("#app");
