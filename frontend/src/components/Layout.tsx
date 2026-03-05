import { Link, Outlet } from "react-router-dom";
import type { MeResponse } from "../types";

interface LayoutProps {
  user: MeResponse;
  onLogout: () => void;
}

export function Layout({ user, onLogout }: LayoutProps) {
  const roleLabel = user.role === "ADMIN" ? "Администратор" : "Пользователь";

  return (
    <div className="app-shell">
      <header className="topbar">
        <div>
          <h1>Система бронирования</h1>
          <p>{user.name} {user.surname} ({roleLabel})</p>
        </div>
        <nav>
          <Link to="/resources">Ресурсы</Link>
          <Link to="/bookings">Бронирования</Link>
          <button onClick={onLogout}>Выйти</button>
        </nav>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  );
}
