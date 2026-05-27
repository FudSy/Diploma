import { FormEvent, useEffect, useMemo, useState } from "react";
import {
  createBooking,
  deleteBooking,
  downloadBookingICS,
  getBookingGoogleLink,
  getBookings,
  getResourceAvailability,
  getResources,
  updateBooking,
} from "../api";
import type { Booking, BusySlot, Resource } from "../types";
import { useToast } from "../components/Toast";
import { useConfirm } from "../components/ConfirmDialog";
import { SkeletonList } from "../components/Skeleton";

interface Props {
  token: string;
}

interface BookingForm {
  resource_id: string;
  start_time: string;
  end_time: string;
}

const initialForm: BookingForm = { resource_id: "", start_time: "", end_time: "" };

const STATUS_LABELS: Record<string, string> = {
  CONFIRMED: "Подтверждено",
  CANCELLED: "Отменено",
  PENDING: "Ожидает",
};

const TYPE_ICONS: Record<string, string> = {
  MEETING_ROOM: "🏢",
  CAR: "🚗",
  DEVICE: "💻",
};

const TYPE_LABELS: Record<string, string> = {
  MEETING_ROOM: "Переговорная",
  CAR: "Автомобиль",
  DEVICE: "Устройство",
};

type FilterMode = "upcoming" | "live" | "past" | "cancelled" | "all";

const FILTER_LABELS: Record<FilterMode, string> = {
  upcoming: "Предстоящие",
  live: "Идут сейчас",
  past: "Прошедшие",
  cancelled: "Отменённые",
  all: "Все",
};

function typeIcon(name: string) {
  return TYPE_ICONS[name] ?? "📦";
}

function typeLabel(name: string) {
  return TYPE_LABELS[name] ?? name;
}

function statusLabel(status: string) {
  return STATUS_LABELS[status.toUpperCase()] ?? status;
}

function statusClass(status: string) {
  const s = status.toUpperCase();
  if (s === "CONFIRMED") return "status-confirmed";
  if (s === "CANCELLED") return "status-cancelled";
  return "status-pending";
}

function resolvePhotoUrl(photoUrl: string): string {
  return photoUrl || "";
}

function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString("ru-RU", {
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function pad(n: number) {
  return String(n).padStart(2, "0");
}

// Convert ISO from <input type="datetime-local"> to a Date.
function localInputToDate(s: string): Date | null {
  if (!s) return null;
  const d = new Date(s);
  return isNaN(d.getTime()) ? null : d;
}

// Convert Date to "YYYY-MM-DDTHH:MM" suitable for <input type="datetime-local">.
function dateToLocalInput(d: Date): string {
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

function isLive(b: Booking, now: number): boolean {
  if (b.status.toUpperCase() === "CANCELLED") return false;
  const start = Date.parse(b.start_time);
  const end = Date.parse(b.end_time);
  return start <= now && now < end;
}

function isPast(b: Booking, now: number): boolean {
  return Date.parse(b.end_time) <= now;
}

function hasOverlap(startISO: string, endISO: string, slots: BusySlot[]): BusySlot | null {
  const s = Date.parse(startISO);
  const e = Date.parse(endISO);
  if (isNaN(s) || isNaN(e) || s >= e) return null;
  for (const slot of slots) {
    if (slot.status.toUpperCase() === "CANCELLED") continue;
    const ss = Date.parse(slot.start_time);
    const se = Date.parse(slot.end_time);
    if (s < se && e > ss) return slot;
  }
  return null;
}

export function BookingsPage({ token }: Props) {
  const toast = useToast();
  const confirm = useConfirm();

  const [bookings, setBookings] = useState<Booking[]>([]);
  const [resources, setResources] = useState<Resource[]>([]);
  const [form, setForm] = useState<BookingForm>(initialForm);
  const [busySlots, setBusySlots] = useState<BusySlot[]>([]);
  const [slotsDate, setSlotsDate] = useState<string>(() => new Date().toISOString().slice(0, 10));
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [filter, setFilter] = useState<FilterMode>("upcoming");
  const [now, setNow] = useState<number>(Date.now());

  // Live "now" tick — keeps live/past/upcoming labels accurate.
  useEffect(() => {
    const id = window.setInterval(() => setNow(Date.now()), 30_000);
    return () => clearInterval(id);
  }, []);

  const selectedResource = resources.find((r) => r.id === form.resource_id);

  async function loadAll() {
    try {
      setLoading(true);
      const [bookingsData, resourcesData] = await Promise.all([getBookings(token), getResources(token)]);
      setBookings(bookingsData);
      setResources(resourcesData);
      if (!form.resource_id && resourcesData.length) {
        setForm((prev) => ({ ...prev, resource_id: resourcesData[0].id }));
      }
    } catch (err) {
      toast.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    void loadAll();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (!form.resource_id) return;
    getResourceAvailability(token, form.resource_id, slotsDate)
      .then((res) => setBusySlots(res.busy_slots))
      .catch(() => setBusySlots([]));
  }, [form.resource_id, slotsDate, token]);

  // When user picks start_time and end_time is empty, auto-fill +1h.
  function handleStartChange(value: string) {
    const start = localInputToDate(value);
    if (start && !form.end_time) {
      const end = new Date(start.getTime() + 60 * 60 * 1000);
      setForm({ ...form, start_time: value, end_time: dateToLocalInput(end) });
      return;
    }
    setForm({ ...form, start_time: value });
  }

  function applyDuration(minutes: number) {
    const start = localInputToDate(form.start_time) ?? new Date();
    if (!form.start_time) {
      setForm({
        ...form,
        start_time: dateToLocalInput(start),
        end_time: dateToLocalInput(new Date(start.getTime() + minutes * 60_000)),
      });
      return;
    }
    setForm({ ...form, end_time: dateToLocalInput(new Date(start.getTime() + minutes * 60_000)) });
  }

  const conflict = useMemo(() => {
    if (!form.start_time || !form.end_time) return null;
    const startISO = new Date(form.start_time).toISOString();
    const endISO = new Date(form.end_time).toISOString();
    // Restrict to the same date as slotsDate, which is what busySlots reflects.
    const startDate = form.start_time.slice(0, 10);
    if (startDate !== slotsDate) return null;
    return hasOverlap(startISO, endISO, busySlots);
  }, [form.start_time, form.end_time, busySlots, slotsDate]);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    if (conflict) {
      toast.error("Время пересекается с другой бронью");
      return;
    }
    setSubmitting(true);
    try {
      await createBooking(token, {
        resource_id: form.resource_id,
        start_time: new Date(form.start_time).toISOString(),
        end_time: new Date(form.end_time).toISOString(),
      });
      setForm((prev) => ({ ...prev, start_time: "", end_time: "" }));
      toast.success("Бронирование создано");
      await loadAll();
    } catch (err) {
      toast.error((err as Error).message);
    } finally {
      setSubmitting(false);
    }
  }

  async function handleCancel(id: string) {
    const ok = await confirm({
      title: "Отменить бронирование?",
      message: "Бронь останется в истории со статусом «Отменено».",
      confirmText: "Отменить бронь",
      variant: "danger",
    });
    if (!ok) return;
    try {
      await updateBooking(token, id, { status: "CANCELLED" });
      toast.success("Бронирование отменено");
      await loadAll();
    } catch (err) {
      toast.error((err as Error).message);
    }
  }

  async function handleDelete(id: string) {
    const ok = await confirm({
      title: "Удалить бронирование?",
      message: "Это действие необратимо.",
      confirmText: "Удалить",
      variant: "danger",
    });
    if (!ok) return;
    try {
      await deleteBooking(token, id);
      toast.success("Бронирование удалено");
      await loadAll();
    } catch (err) {
      toast.error((err as Error).message);
    }
  }

  async function handleDownloadICS(id: string) {
    try {
      await downloadBookingICS(token, id);
      toast.success("Файл .ics скачан");
    } catch (err) {
      toast.error((err as Error).message);
    }
  }

  async function handleOpenGoogle(id: string) {
    try {
      const { url } = await getBookingGoogleLink(token, id);
      window.open(url, "_blank", "noopener");
    } catch (err) {
      toast.error((err as Error).message);
    }
  }

  const filteredBookings = useMemo(() => {
    return bookings.filter((b) => {
      const s = b.status.toUpperCase();
      switch (filter) {
        case "upcoming":
          return s !== "CANCELLED" && Date.parse(b.start_time) > now;
        case "live":
          return isLive(b, now);
        case "past":
          return s !== "CANCELLED" && isPast(b, now);
        case "cancelled":
          return s === "CANCELLED";
        case "all":
        default:
          return true;
      }
    });
  }, [bookings, filter, now]);

  const counts = useMemo(() => {
    const c: Record<FilterMode, number> = { upcoming: 0, live: 0, past: 0, cancelled: 0, all: bookings.length };
    for (const b of bookings) {
      const s = b.status.toUpperCase();
      if (s === "CANCELLED") {
        c.cancelled++;
        continue;
      }
      if (isLive(b, now)) c.live++;
      else if (isPast(b, now)) c.past++;
      else c.upcoming++;
    }
    return c;
  }, [bookings, now]);

  return (
    <section className="page">
      <div className="page-header">
        <h2>Мои бронирования</h2>
        <span className="badge badge-count">{filteredBookings.length}</span>
      </div>

      <div className="filter-bar">
        {(Object.keys(FILTER_LABELS) as FilterMode[]).map((mode) => (
          <button
            key={mode}
            type="button"
            className={`filter-chip ${filter === mode ? "filter-chip--active" : ""}`}
            onClick={() => setFilter(mode)}
          >
            {FILTER_LABELS[mode]}
            <span className="options-count"> ({counts[mode]})</span>
          </button>
        ))}
      </div>

      <div className="booking-layout">
        <form className="panel booking-form-panel" onSubmit={handleCreate}>
          <h3>Создать бронирование</h3>
          <label>
            Ресурс
            <select
              value={form.resource_id}
              onChange={(e) => setForm({ ...form, resource_id: e.target.value })}
              required
            >
              {resources.map((resource) => (
                <option key={resource.id} value={resource.id}>
                  {typeIcon(resource.type)} {resource.name}
                </option>
              ))}
            </select>
          </label>

          {selectedResource && (
            <div className="resource-preview">
              {selectedResource.photo_url ? (
                <div className="resource-preview-photo">
                  <img src={resolvePhotoUrl(selectedResource.photo_url)} alt={selectedResource.name} />
                </div>
              ) : (
                <div className="resource-preview-placeholder">
                  <span className="resource-preview-icon">{typeIcon(selectedResource.type)}</span>
                </div>
              )}
              <div className="resource-preview-info">
                <div className="resource-preview-name">{selectedResource.name}</div>
                <div className="card-type-badge" style={{ marginBottom: "0.3rem" }}>
                  <span>{typeIcon(selectedResource.type)}</span>
                  <span className="card-type-label">{typeLabel(selectedResource.type)}</span>
                </div>
                {selectedResource.description && (
                  <p className="resource-preview-desc">{selectedResource.description}</p>
                )}
                <div className="resource-preview-meta">
                  {selectedResource.location && <span>📍 {selectedResource.location}</span>}
                  <span>👥 {selectedResource.capacity}</span>
                  <span className={`status-badge ${selectedResource.is_active ? "status-active" : "status-inactive"}`}>
                    {selectedResource.is_active ? "активен" : "отключен"}
                  </span>
                </div>
              </div>
            </div>
          )}

          <div className="availability-section">
            <div className="availability-header">
              <span className="availability-title">Занятость на дату</span>
              <input
                type="date"
                value={slotsDate}
                onChange={(e) => setSlotsDate(e.target.value)}
                className="date-input-sm"
              />
            </div>
            {busySlots.length === 0 ? (
              <p className="availability-free">Свободно весь день</p>
            ) : (
              <ul className="busy-slots-list">
                {busySlots.map((s) => (
                  <li key={s.booking_id} className="busy-slot">
                    🔴 {new Date(s.start_time).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })}
                    {" — "}
                    {new Date(s.end_time).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })}
                  </li>
                ))}
              </ul>
            )}
          </div>

          <label>
            Начало
            <input
              type="datetime-local"
              value={form.start_time}
              onChange={(e) => handleStartChange(e.target.value)}
              required
            />
          </label>
          <label>
            Конец
            <input
              type="datetime-local"
              value={form.end_time}
              onChange={(e) => setForm({ ...form, end_time: e.target.value })}
              required
            />
          </label>

          <div className="duration-chips">
            <button type="button" className="duration-chip" onClick={() => applyDuration(30)}>+ 30 мин</button>
            <button type="button" className="duration-chip" onClick={() => applyDuration(60)}>+ 1 ч</button>
            <button type="button" className="duration-chip" onClick={() => applyDuration(120)}>+ 2 ч</button>
            <button type="button" className="duration-chip" onClick={() => applyDuration(240)}>+ 4 ч</button>
          </div>

          {conflict && (
            <div className="inline-warn">
              ⚠ Пересечение с бронью{" "}
              {new Date(conflict.start_time).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })}
              {" — "}
              {new Date(conflict.end_time).toLocaleTimeString("ru-RU", { hour: "2-digit", minute: "2-digit" })}
            </div>
          )}

          <button type="submit" className="btn btn-primary" disabled={submitting || !!conflict}>
            {submitting ? <><span className="spinner" /> Создаём…</> : "Забронировать"}
          </button>
        </form>

        <div className="bookings-list-section">
          {loading ? (
            <SkeletonList count={3} />
          ) : filteredBookings.length === 0 ? (
            <div className="empty-state">
              <div className="empty-state-icon">📅</div>
              <h3>Ничего не найдено</h3>
              <p>
                {filter === "upcoming" && "Нет предстоящих бронирований — создайте первое"}
                {filter === "live" && "Сейчас активных бронирований нет"}
                {filter === "past" && "История пуста"}
                {filter === "cancelled" && "Отменённых бронирований нет"}
                {filter === "all" && "У вас пока нет бронирований"}
              </p>
            </div>
          ) : (
            <div className="bookings-list">
              {filteredBookings.map((booking) => {
                const resource = resources.find((r) => r.id === booking.resource_id);
                const live = isLive(booking, now);
                return (
                  <article key={booking.id} className="booking-card">
                    {resource?.photo_url ? (
                      <div className="booking-card-photo">
                        <img src={resolvePhotoUrl(resource.photo_url)} alt={resource?.name} />
                      </div>
                    ) : (
                      <div className="booking-card-photo booking-card-photo--empty">
                        <span>{resource ? typeIcon(resource.type) : "📦"}</span>
                      </div>
                    )}
                    <div className="booking-card-body">
                      <div className="booking-card-header">
                        <h3>{resource?.name || "Ресурс"}</h3>
                        <span className={`status-badge ${live ? "status-live" : statusClass(booking.status)}`}>
                          {live ? "идёт сейчас" : statusLabel(booking.status)}
                        </span>
                      </div>
                      {resource && (
                        <div className="card-type-badge" style={{ marginBottom: "0.2rem" }}>
                          <span>{typeIcon(resource.type)}</span>
                          <span className="card-type-label">{typeLabel(resource.type)}</span>
                        </div>
                      )}
                      {resource?.location && (
                        <p className="booking-card-location">📍 {resource.location}</p>
                      )}
                      <div className="booking-time">
                        <span>🕐 {formatDateTime(booking.start_time)}</span>
                        <span>🕑 {formatDateTime(booking.end_time)}</span>
                      </div>
                      <div className="actions">
                        <button
                          type="button"
                          className="btn btn-ghost btn-sm"
                          onClick={() => handleDownloadICS(booking.id)}
                          title="Скачать .ics для Apple Calendar / Outlook"
                        >
                          📅 .ics
                        </button>
                        <button
                          type="button"
                          className="btn btn-ghost btn-sm"
                          onClick={() => handleOpenGoogle(booking.id)}
                          title="Открыть в Google Calendar"
                        >
                          🟦 Google
                        </button>
                        {booking.status.toUpperCase() !== "CANCELLED" && !isPast(booking, now) && (
                          <button className="btn btn-ghost btn-sm" onClick={() => handleCancel(booking.id)}>
                            Отменить
                          </button>
                        )}
                        <button className="btn btn-danger btn-sm" onClick={() => handleDelete(booking.id)}>
                          Удалить
                        </button>
                      </div>
                    </div>
                  </article>
                );
              })}
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
