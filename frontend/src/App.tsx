import { Navigate, Route, Routes } from "react-router-dom";
import { useEffect, useState } from "react";
import { getMe, logout } from "./api";
import { Layout } from "./components/Layout";
import { AuthPage } from "./pages/AuthPage";
import { ResourcesPage } from "./pages/ResourcesPage";
import { BookingsPage } from "./pages/BookingsPage";
import { StatsPage } from "./pages/StatsPage";
import type { MeResponse } from "./types";

const TOKEN_KEY = "diploma_token";

export default function App() {
  const [token, setToken] = useState<string | null>(() => localStorage.getItem(TOKEN_KEY));
  const [user, setUser] = useState<MeResponse | null>(null);
  const [authError, setAuthError] = useState<string | null>(null);

  useEffect(() => {
    async function bootstrap() {
      if (!token) {
        setUser(null);
        return;
      }
      try {
        setAuthError(null);
        const me = await getMe(token);
        setUser(me);
      } catch (err) {
        setAuthError((err as Error).message);
        setToken(null);
        localStorage.removeItem(TOKEN_KEY);
      }
    }

    void bootstrap();
  }, [token]);

  function handleToken(nextToken: string) {
    setToken(nextToken);
    localStorage.setItem(TOKEN_KEY, nextToken);
  }

  async function handleLogout() {
    if (token) {
      try {
        await logout(token);
      } catch {
        // Ignore logout errors; local cleanup still applies.
      }
    }
    setToken(null);
    setUser(null);
    localStorage.removeItem(TOKEN_KEY);
  }

  if (!token || !user) {
    return (
      <>
        {authError && <p className="error global-error">{authError}</p>}
        <Routes>
          <Route path="*" element={<AuthPage onToken={handleToken} />} />
        </Routes>
      </>
    );
  }

  return (
    <Routes>
      <Route element={<Layout user={user} onLogout={handleLogout} />}>
        <Route path="/resources" element={<ResourcesPage token={token} isAdmin={user.role === "ADMIN"} />} />
        <Route path="/bookings" element={<BookingsPage token={token} />} />
        {user.role === "ADMIN" && (
          <Route path="/stats" element={<StatsPage token={token} />} />
        )}
        <Route path="*" element={<Navigate to="/resources" replace />} />
      </Route>
    </Routes>
  );
}
