import { createApp } from 'vue';
// import { createPinia } from 'pinia';

import { createAuth0 } from '@auth0/auth0-vue';

import { library } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { faTrash, faPlus, faPencil, faRightToBracket } from '@fortawesome/free-solid-svg-icons';

library.add(faTrash, faPlus, faPencil, faRightToBracket);

import App from './App.vue';
import router from './router';

const app = createApp(App);

// app.use(createPinia());
app.use(router);
app.component('font-awesome-icon', FontAwesomeIcon);
app.use(
  createAuth0({
    domain: 'onetwentyseven.us.auth0.com',
    client_id: 'u2hgEu2s28xKcKWw0JRgqlgcA6hRLatk',
    redirect_uri: window.location.origin,
    audience: 'bill-api-development-resource',
  })
);

app.mount('#app');
