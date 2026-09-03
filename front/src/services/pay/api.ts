import { apiFetch, stripeFetch } from "../../api/api.js";
import type { Paise } from "./money.js";
import type {
  CouponApiResponse,
  CouponValidationResult,
  PaymentIntentRequest,
  PaymentIntentResponse,
  PaymentSuccessPayload,
  TopupResponse,
  TransactionItem,
  TransactionResponse,
  WalletBalanceResponse,
  WalletCreateResponse,
  WalletPayRequest,
  WalletPayResponse,
  WalletTopupResponse,
  WalletTransferResponse,
  RefundResponse
} from "./types.js";

export type {
  CouponApiResponse,
  CouponValidationResult,
  PaymentIntentRequest,
  PaymentIntentResponse,
  PaymentSuccessPayload,
  TopupResponse,
  TransactionItem,
  TransactionResponse,
  WalletBalanceResponse,
  WalletCreateResponse,
  WalletPayRequest,
  WalletPayResponse,
  WalletTopupResponse,
  WalletTransferResponse,
  RefundResponse
} from "./types.js";

export async function createPaymentIntent(payload: PaymentIntentRequest): Promise<PaymentIntentResponse> {
  return await stripeFetch<PaymentIntentResponse>("/create-payment-intent", "POST", payload);
}

export async function recordPaymentSuccess(payload: PaymentSuccessPayload): Promise<unknown> {
  return await stripeFetch<unknown>("/payment-success", "POST", payload);
}

export async function getWalletBalance(): Promise<WalletBalanceResponse> {
  return await apiFetch<WalletBalanceResponse>("/wallet/balance");
}

export async function createWalletAccount(): Promise<WalletCreateResponse> {
  return await apiFetch<WalletCreateResponse>("/wallet/create", "POST");
}

export async function topupWallet(
  amount: number | Paise,
  method: string,
  idempotencyKey: string
): Promise<WalletTopupResponse> {
  return await apiFetch<WalletTopupResponse>(
    "/wallet/topup",
    "POST",
    { amount, method },
    { headers: { "Idempotency-Key": idempotencyKey } }
  );
}

export async function getWalletTransactions(
  skip: number,
  limit: number
): Promise<TransactionResponse | TransactionItem[]> {
  return await apiFetch<TransactionResponse | TransactionItem[]>(`/wallet/transactions?skip=${skip}&limit=${limit}`);
}

export async function refundWalletTransaction(
  transactionId: string | number,
  idempotencyKey: string
): Promise<RefundResponse> {
  return await apiFetch<RefundResponse>(
    "/wallet/refund",
    "POST",
    { transaction_id: transactionId },
    { headers: { "Idempotency-Key": idempotencyKey } }
  );
}

export async function transferWallet(
  payload: { recipient_id: string; amount: number | Paise; note?: string },
  idempotencyKey: string
): Promise<WalletTransferResponse> {
  return await apiFetch<WalletTransferResponse>("/wallet/transfer", "POST", payload, {
    headers: { "Idempotency-Key": idempotencyKey }
  });
}

export async function validateCouponCode(
  couponCode: string,
  cartTotal: number
): Promise<CouponValidationResult> {
  if (!couponCode?.trim()) return { valid: false, discount: 0 };

  try {
    const response = await apiFetch<CouponApiResponse>("/cart/validate-coupon", "POST", {
      coupon_code: couponCode,
      cart_total: cartTotal
    });
    return response?.data || { valid: false, discount: 0 };
  } catch (error: any) {
    console.error("Coupon alignment verification error:", error);
    return { valid: false, discount: 0, reason: error?.message || "Validation failed" };
  }
}

export async function payWallet(payload: WalletPayRequest): Promise<WalletPayResponse> {
  return await apiFetch<WalletPayResponse>("/wallet/pay", "POST", payload);
}

export async function payCashOnDelivery(payload: Omit<WalletPayRequest, "method">): Promise<WalletPayResponse> {
  return await apiFetch<WalletPayResponse>("/wallet/pay", "POST", {
    ...payload,
    method: "cod"
  });
}

export async function requestWalletTopup(
  amount: number,
  paymentMethod: string
): Promise<TopupResponse | undefined> {
  try {
    const res = await apiFetch<{ data?: TopupResponse }>('/wallet/topup', 'POST', {
      amount,
      payment_method: paymentMethod
    });
    return res?.data;
  } catch (err) {
    console.error("Topup submission broken:", err);
    throw err;
  }
}
