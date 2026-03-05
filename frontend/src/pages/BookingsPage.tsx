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

export function BookingsPage({ token }: Props) {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [resources, setResources] = useState<Resource[]>([]);
  const [form, setForm] = useState<BookingForm>(initialForm);
  const [error, setError] = useState<string | null>(null);

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
      <h2>Мои бронирования</h2>
      {error && <p className="error">{error}</p>}

      <form className="panel" onSubmit={handleCreate}>
        <h3>Создать бронирование</h3>
        <label>
          Ресурс
          <select value={form.resource_id} onChange={(e) => setForm({ ...form, resource_id: e.target.value })} required>
            {resources.map((resource) => (
              <option key={resource.id} value={resource.id}>{resource.name}</option>
            ))}
          </select>
        </label>
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
        <button type="submit">Забронировать</button>
      </form>

      <div className="cards-grid">
        {bookings.map((booking) => (
          <article key={booking.id} className="card">
            <h3>{resources.find((r) => r.id === booking.resource_id)?.name || booking.resource_id}</h3>
            <p>Статус: {booking.status}</p>
            <p>Начало: {new Date(booking.start_time).toLocaleString("ru-RU")}</p>
            <p>Конец: {new Date(booking.end_time).toLocaleString("ru-RU")}</p>
            <div className="actions">
              <button onClick={() => handleCancel(booking.id)}>Отменить</button>
              <button onClick={() => handleDelete(booking.id)}>Удалить</button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
