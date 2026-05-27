import { createContext, useCallback, useContext, useMemo, useState } from "react";
import type { ReactNode } from "react";

export type ToastKind = "success" | "error" | "info";

interface Toast {
  id: number;
  kind: ToastKind;
  message: string;
}

interface ToastContextValue {
  push: (kind: ToastKind, message: string) => void;
  success: (message: string) => void;
  error: (message: string) => void;
  info: (message: string) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

let toastCounter = 0;

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const dismiss = useCallback((id: number) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const push = useCallback(
    (kind: ToastKind, message: string) => {
      const id = ++toastCounter;
      setToasts((prev) => [...prev, { id, kind, message }]);
      const ttl = kind === "error" ? 6000 : 3500;
      window.setTimeout(() => dismiss(id), ttl);
    },
    [dismiss]
  );

  const value = useMemo<ToastContextValue>(
    () => ({
      push,
      success: (m) => push("success", m),
      error: (m) => push("error", m),
      info: (m) => push("info", m),
    }),
    [push]
  );

  return (
    <ToastContext.Provider value={value}>
      {children}
      <div className="toast-stack" role="status" aria-live="polite">
        {toasts.map((t) => (
          <div key={t.id} className={`toast toast--${t.kind}`} onClick={() => dismiss(t.id)}>
            <span className="toast-icon" aria-hidden>
              {t.kind === "success" ? "✓" : t.kind === "error" ? "✕" : "ℹ"}
            </span>
            <span className="toast-message">{t.message}</span>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast(): ToastContextValue {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast must be used inside <ToastProvider>");
  return ctx;
}

// Helper hook: wraps an async fn — shows toast on error, returns success boolean.
export function useAsyncAction() {
  const toast = useToast();
  return useCallback(
    async (fn: () => Promise<unknown>, successMsg?: string): Promise<boolean> => {
      try {
        await fn();
        if (successMsg) toast.success(successMsg);
        return true;
      } catch (err) {
        toast.error((err as Error).message);
        return false;
      }
    },
    [toast]
  );
}

export type { ToastContextValue };
