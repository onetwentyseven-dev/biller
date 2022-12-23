import { createApp } from 'vue';
// import { createPinia } from 'pinia';

import { library } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import {
  faTrash,
  faEye,
  faPlus,
  faPencil,
  faRightToBracket,
} from '@fortawesome/free-solid-svg-icons';

library.add(faTrash, faPlus, faPencil, faRightToBracket, faEye);

import App from './App.vue';
import router from './router';
import { auth0 } from './auth0';

const app = createApp(App);

// app.use(createPinia());
app.use(auth0);
app.use(router);
app.component('font-awesome-icon', FontAwesomeIcon);

app.mount('#app');
