import { FormEvent, useState } from "react";
import { login, register } from "../api";
import type { LoginRequest, RegisterRequest } from "../types";

interface AuthPageProps {
  onToken: (token: string) => void;
}

const loginInitial: LoginRequest = { login: "", password: "" };
const registerInitial: RegisterRequest = {
  login: "",
  password: "",
  email: "",
  name: "",
  surname: ""
};

export function AuthPage({ onToken }: AuthPageProps) {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [loginForm, setLoginForm] = useState(loginInitial);
  const [registerForm, setRegisterForm] = useState(registerInitial);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function handleLogin(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      const result = await login(loginForm);
      onToken(result.token);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  async function handleRegister(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      await register(registerForm);
      const result = await login({ login: registerForm.login, password: registerForm.password });
      onToken(result.token);
    } catch (err) {
      setError((err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="auth-wrap">
      <div className="auth-card">
        <h2>{mode === "login" ? "Добро пожаловать" : "Создать аккаунт"}</h2>
        <p className="auth-subtitle">Система бронирования ресурсов</p>
        <div className="segmented">
          <button className={mode === "login" ? "active" : ""} onClick={() => setMode("login")}>Вход</button>
          <button className={mode === "register" ? "active" : ""} onClick={() => setMode("register")}>Регистрация</button>
        </div>

        {mode === "login" ? (
          <form onSubmit={handleLogin}>
            <label>
              Логин
              <input value={loginForm.login} onChange={(e) => setLoginForm({ ...loginForm, login: e.target.value })} required />
            </label>
            <label>
              Пароль
              <input type="password" value={loginForm.password} onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })} required />
            </label>
            <button type="submit" disabled={busy}>Войти</button>
          </form>
        ) : (
          <form onSubmit={handleRegister}>
            <label>
              Логин
              <input value={registerForm.login} onChange={(e) => setRegisterForm({ ...registerForm, login: e.target.value })} required />
            </label>
            <label>
              Пароль
              <input type="password" minLength={6} value={registerForm.password} onChange={(e) => setRegisterForm({ ...registerForm, password: e.target.value })} required />
            </label>
            <label>
              Email
              <input type="email" value={registerForm.email} onChange={(e) => setRegisterForm({ ...registerForm, email: e.target.value })} required />
            </label>
            <label>
              Имя
              <input value={registerForm.name} onChange={(e) => setRegisterForm({ ...registerForm, name: e.target.value })} required />
            </label>
            <label>
              Фамилия
              <input value={registerForm.surname} onChange={(e) => setRegisterForm({ ...registerForm, surname: e.target.value })} required />
            </label>
            <button type="submit" disabled={busy}>Создать аккаунт</button>
          </form>
        )}

        {error && <p className="error">{error}</p>}
      </div>
    </div>
  );
}
