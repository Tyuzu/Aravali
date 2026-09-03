import type { Paise } from "./money.js";

export type PaymentType = "funding" | "purchase";
export type PaymentMethod = "card" | "wallet" | "cash_on_delivery";

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
}

export interface PaymentIntentResponse {
  clientSecret?: string;
  [key: string]: unknown;
}

export interface PaymentSuccessPayload extends PaymentIntentRequest {
  paymentIntentId: string;
}

export interface WalletBalanceResponse {
  exists?: boolean;
  accountExists?: boolean;
  balance?: number;
  currency?: string;
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
  id?: string | number;
  [key: string]: unknown;
}

export interface TransactionItem {
  id: string | number;
  type?: string;
  amount: Paise | number;
  method?: string;
  status?: string;
  created_at: string | number | Date;
  from_account?: string | number;
  userid?: string | number;
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
