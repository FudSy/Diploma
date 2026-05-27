import { FormEvent, useEffect, useRef, useState } from "react";
import { login, register } from "../api";
import type { LoginRequest, RegisterRequest } from "../types";
import { useToast } from "../components/Toast";

interface AuthPageProps {
  onToken: (token: string) => void;
}

const loginInitial: LoginRequest = { login: "", password: "" };
const registerInitial: RegisterRequest = {
  login: "",
  password: "",
  email: "",
  name: "",
  surname: "",
};

export function AuthPage({ onToken }: AuthPageProps) {
  const toast = useToast();
  const [mode, setMode] = useState<"login" | "register">("login");
  const [loginForm, setLoginForm] = useState(loginInitial);
  const [registerForm, setRegisterForm] = useState(registerInitial);
  const [busy, setBusy] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const loginInputRef = useRef<HTMLInputElement>(null);
  const registerInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (mode === "login") loginInputRef.current?.focus();
    else registerInputRef.current?.focus();
  }, [mode]);

  async function handleLogin(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    try {
      const result = await login(loginForm);
      onToken(result.token);
      toast.success("Добро пожаловать");
    } catch (err) {
      toast.error((err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  async function handleRegister(e: FormEvent) {
    e.preventDefault();
    setBusy(true);
    try {
      await register(registerForm);
      const result = await login({ login: registerForm.login, password: registerForm.password });
      onToken(result.token);
      toast.success("Аккаунт создан");
    } catch (err) {
      toast.error((err as Error).message);
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
          <button className={mode === "login" ? "active" : ""} onClick={() => setMode("login")}>
            Вход
          </button>
          <button className={mode === "register" ? "active" : ""} onClick={() => setMode("register")}>
            Регистрация
          </button>
        </div>

        {mode === "login" ? (
          <form onSubmit={handleLogin}>
            <label>
              Логин
              <input
                ref={loginInputRef}
                value={loginForm.login}
                onChange={(e) => setLoginForm({ ...loginForm, login: e.target.value })}
                autoComplete="username"
                required
              />
            </label>
            <label>
              Пароль
              <div className="password-row">
                <input
                  type={showPassword ? "text" : "password"}
                  value={loginForm.password}
                  onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })}
                  autoComplete="current-password"
                  required
                />
                <button
                  type="button"
                  className="password-toggle"
                  onClick={() => setShowPassword((v) => !v)}
                  aria-label={showPassword ? "Скрыть пароль" : "Показать пароль"}
                  tabIndex={-1}
                >
                  {showPassword ? "🙈" : "👁"}
                </button>
              </div>
            </label>
            <button type="submit" disabled={busy}>
              {busy ? <><span className="spinner" /> Входим…</> : "Войти"}
            </button>
          </form>
        ) : (
          <form onSubmit={handleRegister}>
            <label>
              Логин
              <input
                ref={registerInputRef}
                value={registerForm.login}
                onChange={(e) => setRegisterForm({ ...registerForm, login: e.target.value })}
                autoComplete="username"
                required
              />
            </label>
            <label>
              Пароль
              <div className="password-row">
                <input
                  type={showPassword ? "text" : "password"}
                  minLength={6}
                  value={registerForm.password}
                  onChange={(e) => setRegisterForm({ ...registerForm, password: e.target.value })}
                  autoComplete="new-password"
                  required
                />
                <button
                  type="button"
                  className="password-toggle"
                  onClick={() => setShowPassword((v) => !v)}
                  aria-label={showPassword ? "Скрыть пароль" : "Показать пароль"}
                  tabIndex={-1}
                >
                  {showPassword ? "🙈" : "👁"}
                </button>
              </div>
            </label>
            <label>
              Email
              <input
                type="email"
                value={registerForm.email}
                onChange={(e) => setRegisterForm({ ...registerForm, email: e.target.value })}
                autoComplete="email"
                required
              />
            </label>
            <label>
              Имя
              <input
                value={registerForm.name}
                onChange={(e) => setRegisterForm({ ...registerForm, name: e.target.value })}
                autoComplete="given-name"
                required
              />
            </label>
            <label>
              Фамилия
              <input
                value={registerForm.surname}
                onChange={(e) => setRegisterForm({ ...registerForm, surname: e.target.value })}
                autoComplete="family-name"
                required
              />
            </label>
            <button type="submit" disabled={busy}>
              {busy ? <><span className="spinner" /> Создаём…</> : "Создать аккаунт"}
            </button>
          </form>
        )}
      </div>
    </div>
  );
}
