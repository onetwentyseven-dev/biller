import type { IProvider, ICreateUpdateProvider } from './types/provider';
import type {
  IBill,
  IBillSheet,
  IBillSheetEntry,
  ICreateUpdateBill,
  ICreateUpdateBillSheet,
  ICreateUpdateBillSheetEntry,
} from './types/bill';
import type { Token } from './types/token';
import type { IReceipt, ICreateUpdateReceipt } from './types/receipt';

interface IAPIRequest {
  path: string;
  method: 'GET' | 'POST' | 'DELETE' | 'PATCH';
  returnContentType?: 'json' | 'blob';
  headers?: { [key: string]: string };
  body?: XMLHttpRequestBodyInit;
  token?: Token;
}

async function APIRequest<ReturnType>(opts: IAPIRequest): Promise<ReturnType> {
  const uri = new URL('https://fi478t61sj.execute-api.us-east-1.amazonaws.com');
  uri.pathname = opts.path;
  let headers: { [key: string]: string } = {};
  if (opts.headers) {
    headers = opts.headers;
  }

  if (opts.token) {
    headers['Authorization'] = `Bearer ${opts.token}`;
  }

  return await fetch(uri.href, {
    method: opts.method,
    body: opts.body,
    headers: headers,
  }).then((r) => {
    if (r.status >= 400) {
      throw new Error(
        `unexpected status code ${r.status} received on ${opts.method.toUpperCase()} ${opts.path}`
      );
    }

    if (r.status === 204) {
      return undefined as ReturnType;
    }

    switch (opts.returnContentType) {
      case 'blob':
        return r.blob() as ReturnType;
      default:
        return r.text().then((r) => {
          return JSON.parse(r, (key: string, value: any) => {
            const dateKey = ['ts_created', 'ts_updated', 'date_paid', 'date_due'];
            if (dateKey.includes(key)) {
              return new Date(value);
            }
            return value;
          });
        });
    }
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

export const SheetRequest = {
  List: async (): Promise<IBillSheet[]> =>
    APIRequest<IBillSheet[]>({
      method: 'GET',
      path: '/sheets',
    }),
  Get: async (id: string): Promise<IBillSheet> =>
    APIRequest<IBillSheet>({
      method: 'GET',
      path: `/sheets/${id}`,
    }),
  Create: async (sheet: ICreateUpdateBillSheet): Promise<IBillSheet> =>
    APIRequest<IBillSheet>({
      method: 'POST',
      path: '/sheets',
      body: JSON.stringify(sheet),
    }),
  Update: async (id: string, sheet: ICreateUpdateBillSheet): Promise<IBillSheet> =>
    APIRequest<IBillSheet>({
      method: 'PATCH',
      path: `/sheets/${id}`,
      body: JSON.stringify(sheet),
    }),
  Delete: async (id: string): Promise<IBillSheet> =>
    APIRequest<IBillSheet>({
      method: 'DELETE',
      path: `/sheets/${id}`,
    }),
};

export const EntryRequest = {
  ListBySheetID: async (id: string): Promise<IBillSheetEntry[]> =>
    APIRequest<IBillSheetEntry[]>({
      method: 'GET',
      path: `/sheets/${id}/entries`,
    }),
  CreateBySheetID: async (id: string, entry: IBillSheetEntry): Promise<IBillSheetEntry> =>
    APIRequest<IBillSheetEntry>({
      method: 'POST',
      path: `/sheets/${id}/entries`,
      body: JSON.stringify(entry),
    }),
  UpdateBySheetID: async (
    id: string,
    entryID: string,
    entry: ICreateUpdateBillSheetEntry
  ): Promise<IBillSheetEntry> =>
    APIRequest<IBillSheetEntry>({
      method: 'PATCH',
      path: `/sheets/${id}/entries/${entryID}`,
      body: JSON.stringify(entry),
    }),
  DeleteByEntryID: async (id: string, entryID: string): Promise<IBillSheetEntry> =>
    APIRequest<IBillSheetEntry>({
      method: 'DELETE',
      path: `/sheets/${id}/entries/${entryID}`,
    }),
};

export const ReceiptRequest = {
  List: async (): Promise<IReceipt[]> =>
    APIRequest<IReceipt[]>({
      method: 'GET',
      path: '/receipts',
    }),
  Create: async (receipt: ICreateUpdateReceipt): Promise<IReceipt> =>
    APIRequest<IReceipt>({
      method: 'POST',
      path: '/receipts',
      body: JSON.stringify(receipt),
    }),
  Get: async (id: string): Promise<IReceipt> =>
    APIRequest<IReceipt>({
      method: 'GET',
      path: `/receipts/${id}`,
    }),
  Update: async (id: string, receipt: ICreateUpdateBillSheetEntry): Promise<IReceipt> =>
    APIRequest<IReceipt>({
      method: 'PATCH',
      path: `/receipts/${id}`,
      body: JSON.stringify(receipt),
    }),
  Delete: async (id: string, receipt: ICreateUpdateBillSheetEntry): Promise<IReceipt> =>
    APIRequest<IReceipt>({
      method: 'GET',
      path: `/receipts/${id}`,
      body: JSON.stringify(receipt),
    }),
  GetFile: async (id: string): Promise<Blob> =>
    APIRequest<Blob>({
      method: 'GET',
      path: `/receipts/${id}/file`,
      returnContentType: 'blob',
    }),
  PostFile: async (id: string, fileType: string, file: Blob): Promise<void> =>
    APIRequest<void>({
      method: 'POST',
      path: `/receipts/${id}/file`,
      headers: {
        'Content-Type': fileType,
      },
      body: file,
    }),
  DeleteFile: async (id: string): Promise<never> =>
    APIRequest<never>({
      method: 'DELETE',
      path: `/receipts/${id}/file`,
    }),
};

export const WarmerRequest = {
  Get: async (): Promise<never> =>
    APIRequest<never>({
      method: 'GET',
      path: `/warmer`,
    }),
};
