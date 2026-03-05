import { FormEvent, useEffect, useState } from "react";
import { createResource, deleteResource, getResources } from "../api";
import type { Resource } from "../types";

interface Props {
  token: string;
  isAdmin: boolean;
}

const initialForm: Omit<Resource, "id"> = {
  name: "",
  description: "",
  type: "MEETING_ROOM",
  capacity: 1,
  is_active: true
};

export function ResourcesPage({ token, isAdmin }: Props) {
  const [resources, setResources] = useState<Resource[]>([]);
  const [form, setForm] = useState(initialForm);
  const [error, setError] = useState<string | null>(null);

  async function load() {
    try {
      setError(null);
      setResources(await getResources(token));
    } catch (err) {
      setError((err as Error).message);
    }
  }

  useEffect(() => {
    void load();
  }, []);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    try {
      await createResource(token, form);
      setForm(initialForm);
      await load();
    } catch (err) {
      setError((err as Error).message);
    }
  }

  async function handleDelete(id: string) {
    try {
      await deleteResource(token, id);
      await load();
    } catch (err) {
      setError((err as Error).message);
    }
  }

  return (
    <section className="page">
      <h2>Ресурсы</h2>
      {error && <p className="error">{error}</p>}
      <div className="cards-grid">
        {resources.map((r) => (
          <article key={r.id} className="card">
            <h3>{r.name}</h3>
            <p>{r.description || "Без описания"}</p>
            <p>Тип: {r.type}</p>
            <p>Вместимость: {r.capacity}</p>
            <p>Статус: {r.is_active ? "активен" : "отключен"}</p>
            {isAdmin && <button onClick={() => handleDelete(r.id)}>Удалить</button>}
          </article>
        ))}
      </div>

      {isAdmin && (
        <form className="panel" onSubmit={handleCreate}>
          <h3>Добавить ресурс</h3>
          <label>
            Название
            <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          </label>
          <label>
            Описание
            <input value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
          </label>
          <label>
            Тип
            <select value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value as Resource["type"] })}>
              <option value="MEETING_ROOM">MEETING_ROOM</option>
              <option value="CAR">CAR</option>
              <option value="DEVICE">DEVICE</option>
            </select>
          </label>
          <label>
            Вместимость
            <input type="number" min={1} value={form.capacity} onChange={(e) => setForm({ ...form, capacity: Number(e.target.value) })} required />
          </label>
          <label className="inline">
            <input type="checkbox" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} />
            Активен
          </label>
          <button type="submit">Создать</button>
        </form>
      )}
    </section>
  );
}
