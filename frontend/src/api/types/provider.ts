export interface IProvider {
  id: string;
  name: string;
  web_address?: string;
  ts_created: string;
  ts_updated: string;
}
export interface ICreateUpdateProvider {
  name: string;
  web_address: string;
}
