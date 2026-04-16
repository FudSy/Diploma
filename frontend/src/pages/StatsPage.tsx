import { useEffect, useState } from "react";
import { getAdminBookings, getStats } from "../api";
import type { AdminBooking, StatsOverview } from "../types";

interface Props {
  token: string;
}

const TYPE_LABELS: Record<string, string> = {
  MEETING_ROOM: "Переговорные",
  CAR: "Автомобили",
  DEVICE: "Устройства",
};

const TYPE_ICONS: Record<string, string> = {
  MEETING_ROOM: "🏢",
  CAR: "🚗",
  DEVICE: "💻",
};

const STATUS_LABELS: Record<string, string> = {
  CONFIRMED: "Подтверждено",
  CANCELLED: "Отменено",
  PENDING: "Ожидает",
};

function statusClass(s: string) {
  if (s.toUpperCase() === "CONFIRMED") return "status-confirmed";
  if (s.toUpperCase() === "CANCELLED") return "status-cancelled";
  return "status-pending";
}

function formatDateTime(iso: string) {
  return new Date(iso).toLocaleString("ru-RU", {
    day: "2-digit", month: "short", year: "numeric",
    hour: "2-digit", minute: "2-digit",
  });
}

function Bar({ value, max, color }: { value: number; max: number; color: string }) {
  const pct = max > 0 ? Math.round((value / max) * 100) : 0;
  return (
    <div className="stat-bar-track">
      <div className="stat-bar-fill" style={{ width: `${pct}%`, background: color }} />
    </div>
  );
}

export function StatsPage({ token }: Props) {
  const [stats, setStats] = useState<StatsOverview | null>(null);
  const [bookings, setBookings] = useState<AdminBooking[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [bookingFilter, setBookingFilter] = useState("");

  useEffect(() => {
    async function load() {
      try {
        const [s, b] = await Promise.all([getStats(token), getAdminBookings(token)]);
        setStats(s);
        setBookings(b);
      } catch (err) {
        setError((err as Error).message);
      } finally {
        setLoading(false);
      }
    }
    void load();
  }, [token]);

  const maxTypeCount = stats ? Math.max(...stats.bookings_by_type.map((t) => t.count), 1) : 1;
  const maxResourceCount = stats ? Math.max(...stats.top_resources.map((r) => r.count), 1) : 1;
  const maxDayCount = stats ? Math.max(...stats.bookings_last_30_days.map((d) => d.count), 1) : 1;
  const maxHourCount = stats ? Math.max(...stats.peak_hours.map((h) => h.count), 1) : 1;

  const filteredBookings = bookings.filter((b) => {
    if (!bookingFilter) return true;
    const q = bookingFilter.toLowerCase();
    return (
      b.user_name.toLowerCase().includes(q) ||
      b.resource_name.toLowerCase().includes(q) ||
      b.status.toLowerCase().includes(q)
    );
  });

  if (loading) return <div className="page"><p className="text-muted">Загрузка статистики...</p></div>;

  return (
    <section className="page">
      <div className="page-header">
        <h2>Аналитика</h2>
      </div>

      {error && <p className="error">{error}</p>}

      {stats && (
        <>
          {/* KPI cards */}
          <div className="kpi-grid">
            <div className="kpi-card">
              <div className="kpi-value">{stats.total_bookings}</div>
              <div className="kpi-label">Всего бронирований</div>
            </div>
            <div className="kpi-card kpi-card--success">
              <div className="kpi-value">{stats.active_bookings}</div>
              <div className="kpi-label">Активных</div>
            </div>
            <div className="kpi-card kpi-card--danger">
              <div className="kpi-value">{stats.cancelled_bookings}</div>
              <div className="kpi-label">Отменено</div>
            </div>
            <div className="kpi-card">
              <div className="kpi-value">{stats.active_resources} / {stats.total_resources}</div>
              <div className="kpi-label">Активных ресурсов</div>
            </div>
          </div>

          <div className="stats-grid">
            {/* Bookings by type */}
            <div className="panel">
              <h3>Бронирования по типу ресурса</h3>
              {stats.bookings_by_type.length === 0 ? (
                <p className="text-muted">Нет данных</p>
              ) : (
                <ul className="stat-list">
                  {stats.bookings_by_type.map((t) => (
                    <li key={t.type} className="stat-list-item">
                      <span className="stat-label">
                        {TYPE_ICONS[t.type] ?? "📦"} {TYPE_LABELS[t.type] ?? t.type}
                      </span>
                      <Bar value={t.count} max={maxTypeCount} color="var(--accent)" />
                      <span className="stat-count">{t.count}</span>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {/* Top resources */}
            <div className="panel">
              <h3>Топ ресурсов по популярности</h3>
              {stats.top_resources.length === 0 ? (
                <p className="text-muted">Нет данных</p>
              ) : (
                <ul className="stat-list">
                  {stats.top_resources.map((r) => (
                    <li key={r.resource_id.toString()} className="stat-list-item">
                      <span className="stat-label">
                        {TYPE_ICONS[r.type] ?? "📦"} {r.resource_name}
                      </span>
                      <Bar value={r.count} max={maxResourceCount} color="#10b981" />
                      <span className="stat-count">{r.count}</span>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {/* Peak hours */}
            <div className="panel">
              <h3>Пиковые часы бронирований</h3>
              {stats.peak_hours.length === 0 ? (
                <p className="text-muted">Нет данных</p>
              ) : (
                <ul className="stat-list stat-list--compact">
                  {stats.peak_hours.map((h) => (
                    <li key={h.hour} className="stat-list-item">
                      <span className="stat-label">{String(h.hour).padStart(2, "0")}:00</span>
                      <Bar value={h.count} max={maxHourCount} color="#f59e0b" />
                      <span className="stat-count">{h.count}</span>
                    </li>
                  ))}
                </ul>
              )}
            </div>

            {/* Last 30 days */}
            <div className="panel">
              <h3>Бронирования за последние 30 дней</h3>
              {stats.bookings_last_30_days.length === 0 ? (
                <p className="text-muted">Нет данных за период</p>
              ) : (
                <ul className="stat-list stat-list--compact">
                  {stats.bookings_last_30_days.map((d) => (
                    <li key={d.date} className="stat-list-item">
                      <span className="stat-label">{d.date}</span>
                      <Bar value={d.count} max={maxDayCount} color="var(--accent)" />
                      <span className="stat-count">{d.count}</span>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </>
      )}

      {/* All bookings table */}
      <div className="panel" style={{ marginTop: "1.5rem" }}>
        <div className="panel-header-row">
          <h3>Все бронирования</h3>
          <input
            className="filter-input"
            placeholder="Поиск по пользователю, ресурсу, статусу..."
            value={bookingFilter}
            onChange={(e) => setBookingFilter(e.target.value)}
          />
        </div>
        {filteredBookings.length === 0 ? (
          <p className="text-muted">Нет бронирований</p>
        ) : (
          <div className="table-wrap">
            <table className="bookings-table">
              <thead>
                <tr>
                  <th>Пользователь</th>
                  <th>Ресурс</th>
                  <th>Начало</th>
                  <th>Конец</th>
                  <th>Статус</th>
                </tr>
              </thead>
              <tbody>
                {filteredBookings.map((b) => (
                  <tr key={b.id}>
                    <td>{b.user_name}</td>
                    <td>
                      <span>{TYPE_ICONS[b.resource_type] ?? "📦"} {b.resource_name}</span>
                    </td>
                    <td>{formatDateTime(b.start_time)}</td>
                    <td>{formatDateTime(b.end_time)}</td>
                    <td>
                      <span className={`status-badge ${statusClass(b.status)}`}>
                        {STATUS_LABELS[b.status.toUpperCase()] ?? b.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </section>
  );
}
