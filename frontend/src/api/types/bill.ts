export interface IBill {
  id: string;
  provider_id: string;
  name: string;
  amount_paid: number;
  amount_due: number;
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
  amount_due?: number;
  amount_paid?: number;
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
  bill_name: string;
  date_due: Date;
  amount_due: number;
  receipt_id?: string;
  receipt_name?: string;
  date_paid?: Date;
  amount_paid?: number;
  ts_created: Date;
  ts_updated: Date;
}

export interface ICreateUpdateBillSheetEntry {
  bill_id: string;
  date_due: Date;
  amount_due: number;
  receipt_id?: string;
  date_paid?: Date;
  amount_paid?: number;
}
