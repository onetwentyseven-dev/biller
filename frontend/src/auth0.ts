import { createAuth0 } from '@auth0/auth0-vue';

export const auth0 = createAuth0({
  domain: 'onetwentyseven.us.auth0.com',
  client_id: 'u2hgEu2s28xKcKWw0JRgqlgcA6hRLatk',
  redirect_uri: window.location.origin,
  audience: 'bill-api-development-resource',
});
