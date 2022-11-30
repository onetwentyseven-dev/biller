export interface IReceipt {
  id: string;
  provider_id?: string;
  label: string;
  date_paid: string;
  amount_paid: number;
  ts_created: Date;
  ts_updated: Date;
}

export interface ICreateUpdateReceipt {
  provider_id?: string;
  label: string;
  date_paid: Date;
  amount_paid: number;
}
