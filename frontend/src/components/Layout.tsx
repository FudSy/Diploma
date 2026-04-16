import { Link, Outlet, useLocation } from "react-router-dom";
import type { MeResponse } from "../types";

interface LayoutProps {
  user: MeResponse;
  onLogout: () => void;
}

export function Layout({ user, onLogout }: LayoutProps) {
  const roleLabel = user.role === "ADMIN" ? "Администратор" : "Пользователь";
  const { pathname } = useLocation();

  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="topbar-brand">
          <h1>Система бронирования</h1>
          <p>{user.name} {user.surname} &middot; {roleLabel}</p>
        </div>
        <nav>
          <Link to="/resources" className={pathname.startsWith("/resources") ? "active" : ""}>Ресурсы</Link>
          <Link to="/bookings" className={pathname.startsWith("/bookings") ? "active" : ""}>Бронирования</Link>
          {user.role === "ADMIN" && (
            <Link to="/stats" className={pathname.startsWith("/stats") ? "active" : ""}>Аналитика</Link>
          )}
          <button onClick={onLogout}>Выйти</button>
        </nav>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  );
}
