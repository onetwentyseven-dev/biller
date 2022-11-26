import type { IProvider, ICreateUpdateProvider } from './types/provider';
import type { IBill, ICreateUpdateBill } from './types/bill';

interface IAPIRequest {
  path: string;
  method: 'GET' | 'POST' | 'DELETE' | 'PATCH';
  body?: string;
  // token: Token;
}

async function APIRequest<ReturnType>(opts: IAPIRequest): Promise<ReturnType> {
  const uri = new URL('https://fi478t61sj.execute-api.us-east-1.amazonaws.com');
  uri.pathname = opts.path;

  return await fetch(uri.href, {
    method: opts.method,
    body: opts.body,
    headers: {
      // ...(opts.token ? { Authorization: `Bearer ${opts.token}` } : {})
    },
  }).then((r) => {
    if (r.status >= 400) {
      throw new Error(
        `expected status code ${r.status} received on ${opts.method.toUpperCase()} ${opts.path}`
      );
    }

    if (r.status === 204) {
      return undefined as ReturnType;
    }

    return r.json() as ReturnType;
  });
}

export const ProviderRequest = {
  List: async (): Promise<IProvider[]> =>
    APIRequest<IProvider[]>({
      method: 'GET',
      path: '/providers',
    }),
  Get: async (id: string): Promise<IProvider> =>
    APIRequest<IProvider>({
      method: 'GET',
      path: `/providers/${id}`,
    }),
  GetProviderBills: async (id: string): Promise<IBill[]> =>
    APIRequest<IBill[]>({
      method: 'GET',
      path: `/providers/${id}/bills`,
    }),
  Create: async (provider: ICreateUpdateProvider): Promise<IProvider> =>
    APIRequest<IProvider>({
      method: 'POST',
      path: '/providers',
      body: JSON.stringify(provider),
    }),
  Update: async (id: string, provider: IProvider): Promise<IProvider> =>
    APIRequest<IProvider>({
      method: 'PATCH',
      path: `/providers/${id}`,
      body: JSON.stringify(provider),
    }),
  Delete: async (id: string): Promise<void> =>
    APIRequest<void>({
      method: 'DELETE',
      path: `/providers/${id}`,
    }),
};

export const BillRequest = {
  List: async (): Promise<IBill[]> =>
    APIRequest<IBill[]>({
      method: 'GET',
      path: '/bills',
    }),
  Get: async (id: string): Promise<IBill> =>
    APIRequest<IBill>({
      method: 'GET',
      path: `/bills/${id}`,
    }),
  Create: async (bill: ICreateUpdateBill): Promise<IBill> =>
    APIRequest<IBill>({
      method: 'POST',
      path: '/bills',
      body: JSON.stringify(bill),
    }),
  Update: async (id: string, bill: IBill): Promise<IBill> =>
    APIRequest<IBill>({
      method: 'PATCH',
      path: `/bills/${id}`,
      body: JSON.stringify(bill),
    }),
  Delete: async (id: string): Promise<IBill> =>
    APIRequest<IBill>({
      method: 'DELETE',
      path: `/bills/${id}`,
    }),
};
