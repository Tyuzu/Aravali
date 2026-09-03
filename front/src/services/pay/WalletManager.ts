import { createElement } from "../../components/createElement.js";
import { Button } from "../../components/base/Button.js";
import { formatCurrency, Paise, toPaise } from "./money.js";
import { v4 as uuidv4 } from "uuid";
import Notify from "../../components/ui/Notify.js";
import { getWalletBalance, topupWallet } from "./api.js";
import type { WalletBalanceResponse, WalletTopupResponse } from "./types.js";

/* ───────────────────────────────────────── */
/* Types & Interfaces */
/* ───────────────────────────────────────── */

export interface WalletManagerInstance {
  element: HTMLElement;
  loadBalance: () => Promise<void>;
  destroy: () => void;
}

/* ───────────────────────────────────────── */
/* Utility Functions */
/* ───────────────────────────────────────── */

function parseAmountToPaise(value: string): Paise {
  const amount = Number(value);
  if (Number.isNaN(amount) || amount <= 0) return 0 as Paise;
  return toPaise(amount);
}

/* ───────────────────────────────────────── */
/* Wallet Manager Component */
/* ───────────────────────────────────────── */

export function WalletManager(): WalletManagerInstance {
  let currentIdempotencyKey: string = uuidv4();

  const balanceEl = createElement("div", { id: "wallet-balance", class: "balance-display" }, [
    "Loading balance..."
  ]);

  const handleBalanceChange = () => {
    void loadBalance();
  };

  const globalWindow = window as Window & {
    __walletBalanceListenerRegistered?: boolean;
  };

  if (!globalWindow.__walletBalanceListenerRegistered) {
    globalWindow.addEventListener("wallet:balance-changed", handleBalanceChange);
    globalWindow.__walletBalanceListenerRegistered = true;
  }

  const amountInput = createElement("input", {
    type: "number",
    id: "topup-amount",
    placeholder: "Enter amount in INR",
    min: "1",
    step: "0.01"
  }) as HTMLInputElement;

  const methodSelect = createElement("select", { id: "topup-method" }, [
    createElement("option", { value: "card" }, ["Credit/Debit Card"]),
    createElement("option", { value: "upi" }, ["UPI Ecosystem"])
  ]) as HTMLSelectElement;

  const topupBtn = Button({
    title: "Top Up Account",
    classes: "topup-btn btn-primary",
    events: {
      click: async () => {
        const amountPaise = parseAmountToPaise(amountInput.value);
        const method = methodSelect.value;

        if (amountPaise <= 0) {
          Notify("Please enter a valid amount", { type: "warning" });
          return;
        }

        topupBtn.disabled = true;
        try {
          const res = await topupWallet(amountPaise, method, currentIdempotencyKey);

          if (res?.success) {
            Notify(res.message || "Top-up successful", { type: "success" });
            currentIdempotencyKey = uuidv4();
            amountInput.value = "";

            // Dispatch event for any component listening to balance updates
            window.dispatchEvent(new CustomEvent("wallet:balance-changed"));
          } else {
            Notify(res?.message || "Transaction declined by gateway", { type: "error" });
          }
        } catch (err: unknown) {
          console.error("Network error:", err);
          Notify("Top-up request failed", { type: "error" });
        } finally {
          topupBtn.disabled = false;
        }
      }
    }
  }) as HTMLButtonElement;

  async function loadBalance(): Promise<void> {
    try {
      const res = await getWalletBalance();
      if (res && res.balance !== undefined) {
        balanceEl.textContent = `Wallet Balance: ${formatCurrency(res.balance)}`;
      } else {
        balanceEl.textContent = "Balance unavailable";
      }
    } catch (err: unknown) {
      console.error("Balance fetch error:", err);
      balanceEl.textContent = "Sync failed";
    }
  }

  loadBalance();

  const destroy = (): void => {
    if (globalWindow.__walletBalanceListenerRegistered) {
      globalWindow.removeEventListener("wallet:balance-changed", handleBalanceChange);
      globalWindow.__walletBalanceListenerRegistered = false;
    }
  };

  return {
    element: createElement("div", { id: "wallet-manager", class: "wallet-card" }, [
      createElement("h3", { class: "wallet-section-title" }, ["Account Balance"]),
      balanceEl,
      createElement("div", { class: "wallet-form" }, [amountInput, methodSelect, topupBtn])
    ]),
    loadBalance,
    destroy
  };
}