import { FormEvent, useEffect, useState } from "react";
import { createBooking, deleteBooking, getBookings, getResources, updateBooking } from "../api";
import type { Booking, Resource } from "../types";

interface Props {
  token: string;
}

interface BookingForm {
  resource_id: string;
  start_time: string;
  end_time: string;
}

const initialForm: BookingForm = {
  resource_id: "",
  start_time: "",
  end_time: ""
};

const STATUS_LABELS: Record<string, string> = {
  CONFIRMED: "Подтверждено",
  CANCELLED: "Отменено",
  PENDING: "Ожидает"
};

const TYPE_ICONS: Record<string, string> = {
  MEETING_ROOM: "🏢",
  CAR: "🚗",
  DEVICE: "💻"
};

const TYPE_LABELS: Record<string, string> = {
  MEETING_ROOM: "Переговорная",
  CAR: "Автомобиль",
  DEVICE: "Устройство"
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
  if (!photoUrl) return "";
  if (photoUrl.startsWith("http")) return photoUrl;
  return photoUrl;
}

function formatDateTime(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleString("ru-RU", {
    day: "2-digit",
    month: "long",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function BookingsPage({ token }: Props) {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [resources, setResources] = useState<Resource[]>([]);
  const [form, setForm] = useState<BookingForm>(initialForm);
  const [error, setError] = useState<string | null>(null);

  const selectedResource = resources.find((r) => r.id === form.resource_id);

  async function loadAll() {
    try {
      setError(null);
      const [bookingsData, resourcesData] = await Promise.all([getBookings(token), getResources(token)]);
      setBookings(bookingsData);
      setResources(resourcesData);
      if (!form.resource_id && resourcesData.length) {
        setForm((prev) => ({ ...prev, resource_id: resourcesData[0].id }));
      }
    } catch (err) {
      setError((err as Error).message);
    }
  }

  useEffect(() => {
    void loadAll();
  }, []);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    try {
      await createBooking(token, {
        resource_id: form.resource_id,
        start_time: new Date(form.start_time).toISOString(),
        end_time: new Date(form.end_time).toISOString()
      });
      setForm((prev) => ({ ...prev, start_time: "", end_time: "" }));
      await loadAll();
    } catch (err) {
      setError((err as Error).message);
    }
  }

  async function handleCancel(id: string) {
    try {
      await updateBooking(token, id, { status: "CANCELLED" });
      await loadAll();
    } catch (err) {
      setError((err as Error).message);
    }
  }

  async function handleDelete(id: string) {
    try {
      await deleteBooking(token, id);
      await loadAll();
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2>Мои бронирования</h2>
        <span className="badge badge-count">{bookings.length}</span>
      </div>

      {error && <p className="error">{error}</p>}

      <div className="booking-layout">
        <form className="panel booking-form-panel" onSubmit={handleCreate}>
          <h3>Создать бронирование</h3>
          <label>
            Ресурс
            <select value={form.resource_id} onChange={(e) => setForm({ ...form, resource_id: e.target.value })} required>
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

          <label>
            Начало
            <input
              type="datetime-local"
              value={form.start_time}
              onChange={(e) => setForm({ ...form, start_time: e.target.value })}
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
          <button type="submit" className="btn btn-primary">Забронировать</button>
        </form>

        <div className="bookings-list-section">
          {bookings.length === 0 ? (
            <div className="empty-state">
              <div className="empty-state-icon">📅</div>
              <h3>Нет бронирований</h3>
              <p>Выберите ресурс и создайте своё первое бронирование</p>
            </div>
          ) : (
            <div className="bookings-list">
              {bookings.map((booking) => {
                const resource = resources.find((r) => r.id === booking.resource_id);
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
                        <span className={`status-badge ${statusClass(booking.status)}`}>
                          {statusLabel(booking.status)}
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
                        {booking.status.toUpperCase() !== "CANCELLED" && (
                          <button className="btn btn-ghost btn-sm" onClick={() => handleCancel(booking.id)}>Отменить</button>
                        )}
                        <button className="btn btn-danger btn-sm" onClick={() => handleDelete(booking.id)}>Удалить</button>
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
