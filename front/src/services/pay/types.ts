import type { Paise } from "./money.js";

export type PaymentType = "funding" | "purchase";
export type PaymentMethod = "card" | "wallet" | "cash_on_delivery" | "cod";
export type DateLike = string | number | Date;

export interface PaymentConfig {
  allowedEntities: string[];
  methods: PaymentMethod[];
}

export type PaymentRules = Record<PaymentType, PaymentConfig>;

export interface ValidationResult {
  valid: boolean;
  error?: string;
}

export interface PaymentIntentRequest {
  paymentType?: string;
  entityType: string;
  entityId: string | number;
  method?: string;
  amount?: number;
}

export interface PaymentIntentResponse {
  clientSecret?: string;
  [key: string]: unknown;
}

export interface PaymentSuccessPayload extends PaymentIntentRequest {
  paymentIntentId: string;
}

export interface WalletAccount {
  id?: string | number;
  userid?: string | number;
  userId?: string | number;
  currency?: string;
  status?: string;
  cached_balance?: number;
  cachedBalance?: number;
  version?: number;
  created_at?: DateLike;
  updated_at?: DateLike;
  [key: string]: unknown;
}

export interface WalletBalanceResponse {
  exists?: boolean;
  accountExists?: boolean;
  balance?: number;
  cached_balance?: number;
  currency?: string;
  userid?: string | number;
  account?: WalletAccount;
  [key: string]: unknown;
}

export interface WalletCreateResponse {
  success: boolean;
  message?: string;
  [key: string]: unknown;
}

export interface WalletTopupResponse {
  success?: boolean;
  message?: string;
  transaction_id?: string | number;
  transactionId?: string | number;
  [key: string]: unknown;
}

export interface WalletPayRequest extends PaymentIntentRequest {
  method: "wallet" | "cod" | string;
  amount?: number;
}

export interface WalletPayResponse {
  success?: boolean;
  message?: string;
  transaction_id?: string | number;
  transactionId?: string | number;
  id?: string | number;
  [key: string]: unknown;
}

export interface TransactionItem {
  id: string | number;
  _id?: string | number;
  userid?: string | number;
  userId?: string | number;
  parent_txn?: string | number;
  type?: string;
  method?: string;
  entity_type?: string;
  entityType?: string;
  entity_id?: string | number;
  entityId?: string | number;
  from_account?: string | number;
  fromAccount?: string | number;
  to_account?: string | number;
  toAccount?: string | number;
  amount: Paise | number;
  currency?: string;
  status?: string;
  external_ref?: string;
  meta?: Record<string, unknown>;
  created_at: DateLike;
  createdAt?: DateLike;
  updated_at?: DateLike;
  updatedAt?: DateLike;
  [key: string]: unknown;
}

export interface TransactionResponse {
  transactions?: TransactionItem[];
  [key: string]: unknown;
}

export interface RefundResponse {
  success?: boolean;
  message?: string;
  [key: string]: unknown;
}

export interface WalletTransferResponse {
  success?: boolean;
  message?: string;
  transaction_id?: string | number;
  transactionId?: string | number;
  [key: string]: unknown;
}

export interface CouponValidationResult {
  valid: boolean;
  discount: number;
  reason?: string;
}

export interface CouponApiResponse {
  data?: CouponValidationResult;
  [key: string]: unknown;
}

export interface TopupResponse {
  transactionId?: string | number;
  transaction_id?: string | number;
  status?: string;
  balance?: number;
  [key: string]: unknown;
}

export interface StripePaymentParams {
  paymentType?: PaymentType;
  entityType: string;
  entityId: string | number;
}

export interface ShowPaymentModalParams extends StripePaymentParams {
  entityName: string;
}

export interface PaymentResult {
  success: boolean;
  paymentIntentId?: string;
  method?: PaymentMethod;
  error?: string;
  message?: string;
  redirectingToStripe?: boolean;
}
