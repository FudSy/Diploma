import { createContext, useCallback, useContext, useEffect, useMemo, useRef, useState } from "react";
import type { ReactNode } from "react";

interface ConfirmOptions {
  title: string;
  message?: string;
  confirmText?: string;
  cancelText?: string;
  variant?: "danger" | "default";
}

interface ConfirmContextValue {
  confirm: (opts: ConfirmOptions) => Promise<boolean>;
}

const ConfirmContext = createContext<ConfirmContextValue | null>(null);

interface PendingState extends ConfirmOptions {
  resolve: (ok: boolean) => void;
}

export function ConfirmProvider({ children }: { children: ReactNode }) {
  const [pending, setPending] = useState<PendingState | null>(null);
  const confirmBtnRef = useRef<HTMLButtonElement>(null);

  const confirm = useCallback((opts: ConfirmOptions): Promise<boolean> => {
    return new Promise<boolean>((resolve) => {
      setPending({ ...opts, resolve });
    });
  }, []);

  const resolveWith = useCallback(
    (ok: boolean) => {
      if (pending) pending.resolve(ok);
      setPending(null);
    },
    [pending]
  );

  useEffect(() => {
    if (!pending) return;
    confirmBtnRef.current?.focus();
    function onKey(e: KeyboardEvent) {
      if (e.key === "Escape") resolveWith(false);
      if (e.key === "Enter") resolveWith(true);
    }
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [pending, resolveWith]);

  const value = useMemo<ConfirmContextValue>(() => ({ confirm }), [confirm]);

  return (
    <ConfirmContext.Provider value={value}>
      {children}
      {pending && (
        <div className="modal-backdrop" onClick={() => resolveWith(false)}>
          <div
            className="modal-card"
            role="dialog"
            aria-modal="true"
            aria-labelledby="confirm-title"
            onClick={(e) => e.stopPropagation()}
          >
            <h3 id="confirm-title" style={{ marginBottom: pending.message ? "0.5rem" : 0 }}>
              {pending.title}
            </h3>
            {pending.message && <p className="text-muted" style={{ margin: 0 }}>{pending.message}</p>}
            <div className="modal-actions">
              <button type="button" className="btn btn-ghost btn-sm" onClick={() => resolveWith(false)}>
                {pending.cancelText ?? "Отмена"}
              </button>
              <button
                ref={confirmBtnRef}
                type="button"
                className={`btn btn-sm ${pending.variant === "danger" ? "btn-danger" : "btn-primary"}`}
                onClick={() => resolveWith(true)}
              >
                {pending.confirmText ?? "Подтвердить"}
              </button>
            </div>
          </div>
        </div>
      )}
    </ConfirmContext.Provider>
  );
}

export function useConfirm(): (opts: ConfirmOptions) => Promise<boolean> {
  const ctx = useContext(ConfirmContext);
  if (!ctx) throw new Error("useConfirm must be used inside <ConfirmProvider>");
  return ctx.confirm;
}
