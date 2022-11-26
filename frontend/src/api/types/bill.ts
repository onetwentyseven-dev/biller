export interface IBill {
  id: string;
  provider_id: string;
  name: string;
  ts_created: Date;
  ts_updated: Date;
}

export interface ICreateUpdateBill {
  provider_id: string;
  name: string;
}

export interface IBillSheet {
  id: string;
  name: string;
  ts_created: Date;
  ts_updated: Date;
}

export interface ICreateUpdateBillSheet {
  name: string;
}

export interface IBillSheetEntry {
  entry_id: string;
  sheet_id: string;
  bill_id: string;
  date_due: Date;
  amount_due: number;
  receipt_id?: string;
  date_paid?: Date;
  amount_paid?: number;
  ts_created: Date;
  ts_updated: Date;
}

export interface ICreateUpdateBillSheet {
  sheet_id: string;
  bill_id: string;
  date_due: Date;
  amount_due: number;
  receipt_id?: string;
  date_paid?: Date;
  amount_paid?: number;
}

export interface IBillReceipt {
  id: string;
  provider_id: string;
  date_paid: Date;
  amount_paid: number;
  ts_created: Date;
  ts_updated: Date;
}
