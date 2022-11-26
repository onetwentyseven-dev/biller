// import { Token } from "./types/token"
import type { IProvider } from "./types/provider";
// import { IBill } from "./types/bill";
import ts from "typescript";
import type { IBill } from "./types/bill";

interface IAPIRequest {
    path: string;
    method: "GET" | "POST" | "DELETE" | "PATCH";
    body?: string;
    // token: Token;
}

async function APIRequest<ReturnType>(opts: IAPIRequest): Promise<ReturnType> {
    const uri = new URL("https://fi478t61sj.execute-api.us-east-1.amazonaws.com")
    uri.pathname = opts.path

    return await fetch(uri.href, {
        method: opts.method,
        body: opts.body,
        headers: {
            // ...(opts.token ? { Authorization: `Bearer ${opts.token}` } : {})
        }
    }).then((r) => {
        if (r.status >= 400) {
            throw new Error(`expected status code ${r.status} received on ${opts.method.toUpperCase()} ${opts.path}`)
        }

        return r.json() as ReturnType
    })
}

export const ProviderRequest = {
    List: async (): Promise<IProvider[]> => APIRequest<IProvider[]>({
        method: "GET",
        path: "/providers",
    }),
    Get: async(id: string): Promise<IProvider> => APIRequest<IProvider>({
        method: "GET",
        path: `/providers/${id}`,
    }),
    GetProviderBills:  async(id: string): Promise<IBill[]> => APIRequest<IBill[]>({
        method: "GET",
        path: `/providers/${id}/bills`,
    })
}

export const BillRequest = {
    List: async (): Promise<IBill[]> => APIRequest<IBill[]>({
        method: "GET",
        path: "/bills"
    })
}