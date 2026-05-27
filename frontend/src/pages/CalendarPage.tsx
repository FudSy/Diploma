import { useEffect, useState } from "react";
import { getMyCalendarFeed, rotateMyCalendarFeed } from "../api";
import type { CalendarFeedResponse } from "../types";
import { useToast } from "../components/Toast";
import { useConfirm } from "../components/ConfirmDialog";

interface Props {
  token: string;
}

function buildGoogleSubscribeUrl(feedUrl: string): string {
  return `https://calendar.google.com/calendar/r/settings/addbyurl?cid=${encodeURIComponent(feedUrl)}`;
}

export function CalendarPage({ token }: Props) {
  const toast = useToast();
  const confirmDialog = useConfirm();
  const [feed, setFeed] = useState<CalendarFeedResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [rotating, setRotating] = useState<boolean>(false);
  const [copied, setCopied] = useState<"feed" | "webcal" | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      try {
        const data = await getMyCalendarFeed(token);
        if (!cancelled) setFeed(data);
      } catch (err) {
        if (!cancelled) toast.error((err as Error).message);
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    void load();
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]);

  async function handleRotate() {
    const ok = await confirmDialog({
      title: "Сбросить токен подписки?",
      message: "Старая ссылка перестанет работать. Подписку нужно будет добавить заново.",
      confirmText: "Сбросить",
      variant: "danger",
    });
    if (!ok) return;
    try {
      setRotating(true);
      const data = await rotateMyCalendarFeed(token);
      setFeed(data);
      toast.success("Токен обновлён");
    } catch (err) {
      toast.error((err as Error).message);
    } finally {
      setRotating(false);
    }
  }

  async function handleCopy(value: string, kind: "feed" | "webcal") {
    try {
      await navigator.clipboard.writeText(value);
      setCopied(kind);
      toast.success("Ссылка скопирована");
      setTimeout(() => setCopied(null), 2000);
    } catch (err) {
      toast.error((err as Error).message);
    }
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2>Календарь — подписка</h2>
      </div>

      <p style={{ color: "var(--text-muted)", marginBottom: "1rem" }}>
        Подпишитесь на личный календарь — все ваши бронирования будут автоматически появляться в&nbsp;Google
        Calendar, Apple Calendar или Outlook. Изменения и&nbsp;отмены обновятся при следующей синхронизации.
      </p>

      {loading ? (
        <p>Загрузка…</p>
      ) : feed ? (
        <div className="panel" style={{ display: "flex", flexDirection: "column", gap: "1.25rem" }}>
          <div>
            <h3 style={{ marginTop: 0 }}>HTTPS-ссылка (Google / Outlook)</h3>
            <div className="calendar-url-row">
              <code className="calendar-url">{feed.feed_url}</code>
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={() => handleCopy(feed.feed_url, "feed")}
              >
                {copied === "feed" ? "✓ Скопировано" : "Копировать"}
              </button>
              <a
                className="btn btn-primary btn-sm"
                href={buildGoogleSubscribeUrl(feed.feed_url)}
                target="_blank"
                rel="noopener noreferrer"
              >
                Добавить в Google Calendar
              </a>
            </div>
          </div>

          <div>
            <h3>Webcal-ссылка (Apple Calendar, iOS)</h3>
            <div className="calendar-url-row">
              <code className="calendar-url">{feed.webcal_url}</code>
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={() => handleCopy(feed.webcal_url, "webcal")}
              >
                {copied === "webcal" ? "✓ Скопировано" : "Копировать"}
              </button>
              <a className="btn btn-primary btn-sm" href={feed.webcal_url}>
                Открыть в календаре
              </a>
            </div>
          </div>

          <details>
            <summary style={{ cursor: "pointer", fontWeight: 500 }}>Как подписаться вручную</summary>
            <ol style={{ marginTop: "0.5rem", paddingLeft: "1.25rem", lineHeight: 1.6 }}>
              <li>
                <b>Google Calendar:</b> «Другие календари» → «Подписаться» → «Из&nbsp;URL» → вставить
                HTTPS-ссылку.
              </li>
              <li>
                <b>Apple Calendar:</b> Файл → Новая подписка на&nbsp;календарь → вставить webcal-ссылку.
              </li>
              <li>
                <b>Outlook:</b> Календарь → Добавить календарь → Из&nbsp;Интернета → вставить HTTPS-ссылку.
              </li>
            </ol>
          </details>

          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between",
              gap: "1rem",
              borderTop: "1px solid var(--border)",
              paddingTop: "1rem"
            }}
          >
            <div style={{ fontSize: "0.9rem", color: "var(--text-muted)" }}>
              Если ссылка попала к&nbsp;посторонним — сбросьте токен. Старая подписка перестанет работать.
            </div>
            <button
              type="button"
              className="btn btn-danger btn-sm"
              onClick={handleRotate}
              disabled={rotating}
            >
              {rotating ? "Сброс…" : "Сбросить токен"}
            </button>
          </div>
        </div>
      ) : null}
    </section>
  );
}
