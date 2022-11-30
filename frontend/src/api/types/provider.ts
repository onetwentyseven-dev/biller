export interface IProvider {
  id: string;
  name: string;
  web_address?: string;
  ts_created: Date;
  ts_updated: Date;
}
export interface ICreateUpdateProvider {
  name: string;
  web_address: string;
}
