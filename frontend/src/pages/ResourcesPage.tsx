import { FormEvent, useEffect, useRef, useState } from "react";
import {
  createResource,
  createResourceType,
  deleteResource,
  deleteResourceType,
  addResourceTypeOption,
  deleteResourceTypeOption,
  getResourceTypes,
  getResources,
  uploadResourcePhoto,
} from "../api";
import type { Resource, ResourceType } from "../types";

interface Props {
  token: string;
  isAdmin: boolean;
}

const initialForm: Omit<Resource, "id"> = {
  name: "",
  description: "",
  type: "",
  capacity: 1,
  is_active: true,
  location: "",
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

const OPTION_TYPE_LABELS: Record<string, string> = {
  text: "Текст",
  number: "Число",
  boolean: "Да/Нет",
};

interface NewOption {
  name: string;
  option_type: "text" | "number" | "boolean";
  is_required: boolean;
}

const emptyOption: NewOption = { name: "", option_type: "text", is_required: false };

function typeIcon(name: string) {
  return TYPE_ICONS[name] ?? "📦";
}

function typeLabel(name: string) {
  return TYPE_LABELS[name] ?? name;
}

function resolvePhotoUrl(photoUrl: string): string {
  if (!photoUrl) return "";
  if (photoUrl.startsWith("http")) return photoUrl;
  return photoUrl;
}

export function ResourcesPage({ token, isAdmin }: Props) {
  const [resources, setResources] = useState<Resource[]>([]);
  const [resourceTypes, setResourceTypes] = useState<ResourceType[]>([]);
  const [form, setForm] = useState(initialForm);
  const [newTypeName, setNewTypeName] = useState("");
  const [newTypeOptions, setNewTypeOptions] = useState<NewOption[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [typeError, setTypeError] = useState<string | null>(null);
  const [uploadingId, setUploadingId] = useState<string | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);
  const [expandedTypeId, setExpandedTypeId] = useState<string | null>(null);
  const [addOptionForm, setAddOptionForm] = useState<NewOption>({ ...emptyOption });
  const [addOptionError, setAddOptionError] = useState<string | null>(null);
  const fileInputRefs = useRef<Record<string, HTMLInputElement | null>>({});

  async function loadResources() {
    try {
      setError(null);
      setResources(await getResources(token));
    } catch (err) {
      setError((err as Error).message);
    }
  }

  async function loadTypes() {
    try {
      const types = await getResourceTypes(token);
      setResourceTypes(types);
      if (!form.type && types.length) {
        setForm((prev) => ({ ...prev, type: types[0].name }));
      }
    } catch {
      // non-critical
    }
  }

  useEffect(() => {
    void loadResources();
    void loadTypes();
  }, []);

  async function handleCreate(e: FormEvent) {
    e.preventDefault();
    try {
      await createResource(token, form);
      setForm({ ...initialForm, type: resourceTypes[0]?.name ?? "" });
      await loadResources();
    } catch (err) {
      setError((err as Error).message);
    }
  }

  async function handleDelete(id: string) {
    try {
      await deleteResource(token, id);
      await loadResources();
    } catch (err) {
      setError((err as Error).message);
    }
  }

  async function handleCreateType(e: FormEvent) {
    e.preventDefault();
    const name = newTypeName.trim().toUpperCase().replace(/\s+/g, "_");
    if (!name) return;
    try {
      setTypeError(null);
      const options = newTypeOptions.filter((o) => o.name.trim() !== "");
      await createResourceType(token, name, options.length > 0 ? options : undefined);
      setNewTypeName("");
      setNewTypeOptions([]);
      await loadTypes();
    } catch (err) {
      setTypeError((err as Error).message);
    }
  }

  async function handleDeleteType(id: string) {
    try {
      await deleteResourceType(token, id);
      if (expandedTypeId === id) setExpandedTypeId(null);
      await loadTypes();
    } catch (err) {
      setTypeError((err as Error).message);
    }
  }

  function handleAddNewTypeOption() {
    setNewTypeOptions([...newTypeOptions, { ...emptyOption }]);
  }

  function handleRemoveNewTypeOption(index: number) {
    setNewTypeOptions(newTypeOptions.filter((_, i) => i !== index));
  }

  function handleChangeNewTypeOption(index: number, field: keyof NewOption, value: string | boolean) {
    setNewTypeOptions(
      newTypeOptions.map((o, i) => (i === index ? { ...o, [field]: value } : o))
    );
  }

  async function handleAddOptionToType(e: FormEvent, resourceTypeId: string) {
    e.preventDefault();
    if (!addOptionForm.name.trim()) return;
    try {
      setAddOptionError(null);
      await addResourceTypeOption(token, resourceTypeId, addOptionForm);
      setAddOptionForm({ ...emptyOption });
      await loadTypes();
    } catch (err) {
      setAddOptionError((err as Error).message);
    }
  }

  async function handleDeleteOption(resourceTypeId: string, optionId: string) {
    try {
      await deleteResourceTypeOption(token, resourceTypeId, optionId);
      await loadTypes();
    } catch (err) {
      setTypeError((err as Error).message);
    }
  }

  async function handlePhotoUpload(resourceId: string) {
    const input = fileInputRefs.current[resourceId];
    if (!input?.files?.length) return;
    const file = input.files[0];
    try {
      setUploadingId(resourceId);
      setUploadError(null);
      await uploadResourcePhoto(token, resourceId, file);
      input.value = "";
      await loadResources();
    } catch (err) {
      setUploadError((err as Error).message);
    } finally {
      setUploadingId(null);
    }
  }

  return (
    <section className="page">
      <div className="page-header">
        <h2>Ресурсы</h2>
        <span className="badge badge-count">{resources.length}</span>
      </div>

      {error && <p className="error">{error}</p>}
      {uploadError && <p className="error">{uploadError}</p>}

      <div className="cards-grid">
        {resources.map((r) => (
          <article key={r.id} className={`card ${!r.is_active ? "card--inactive" : ""}`}>
            {r.photo_url && (
              <div className="card-photo">
                <img src={resolvePhotoUrl(r.photo_url)} alt={r.name} />
              </div>
            )}
            <div className="card-type-badge">
              <span>{typeIcon(r.type)}</span>
              <span className="card-type-label">{typeLabel(r.type)}</span>
            </div>
            <h3>{r.name}</h3>
            <p>{r.description || "Без описания"}</p>
            {r.location && (
              <p className="card-location">📍 {r.location}</p>
            )}
            <div className="card-meta">
              <span className="meta-item">👥 {r.capacity}</span>
              <span className={`status-badge ${r.is_active ? "status-active" : "status-inactive"}`}>
                {r.is_active ? "активен" : "отключен"}
              </span>
            </div>
            {isAdmin && (
              <div className="card-admin-actions">
                <div className="photo-upload-row">
                  <input
                    type="file"
                    accept=".jpg,.jpeg,.png,.webp"
                    ref={(el) => { fileInputRefs.current[r.id] = el; }}
                    className="file-input"
                  />
                  <button
                    className="btn btn-secondary btn-sm"
                    onClick={() => handlePhotoUpload(r.id)}
                    disabled={uploadingId === r.id}
                  >
                    {uploadingId === r.id ? "Загрузка..." : "Загрузить фото"}
                  </button>
                </div>
                <button className="btn btn-danger btn-sm" onClick={() => handleDelete(r.id)}>
                  Удалить
                </button>
              </div>
            )}
          </article>
        ))}
      </div>

      {isAdmin && (
        <div className="admin-panels">
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
              Расположение
              <input
                placeholder="Например: Этаж 3, Комната 301"
                value={form.location ?? ""}
                onChange={(e) => setForm({ ...form, location: e.target.value })}
              />
            </label>
            <label>
              Тип
              <select value={form.type} onChange={(e) => setForm({ ...form, type: e.target.value })}>
                {resourceTypes.map((rt) => (
                  <option key={rt.id} value={rt.name}>{typeIcon(rt.name)} {typeLabel(rt.name)}</option>
                ))}
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
            <button type="submit" className="btn btn-primary">Создать ресурс</button>
          </form>

          <div className="panel">
            <h3>Типы ресурсов</h3>
            {typeError && <p className="error">{typeError}</p>}
            <ul className="type-list">
              {resourceTypes.map((rt) => (
                <li key={rt.id} className="type-list-item-wrap">
                  <div className="type-list-item">
                    <span
                      className="type-name-clickable"
                      onClick={() => setExpandedTypeId(expandedTypeId === rt.id ? null : rt.id)}
                      title="Нажмите, чтобы увидеть опции"
                    >
                      {typeIcon(rt.name)} {typeLabel(rt.name)}
                      {rt.options && rt.options.length > 0 && (
                        <span className="options-count"> ({rt.options.length})</span>
                      )}
                    </span>
                    <button className="btn btn-danger btn-xs" onClick={() => handleDeleteType(rt.id)}>✕</button>
                  </div>

                  {expandedTypeId === rt.id && (
                    <div className="type-options-panel">
                      {rt.options && rt.options.length > 0 ? (
                        <ul className="options-list">
                          {rt.options.map((opt) => (
                            <li key={opt.id} className="option-item">
                              <span className="option-name">{opt.name}</span>
                              <span className="option-type-badge">{OPTION_TYPE_LABELS[opt.option_type] ?? opt.option_type}</span>
                              {opt.is_required && <span className="option-required-badge">обязательное</span>}
                              <button className="btn btn-danger btn-xs" onClick={() => handleDeleteOption(rt.id, opt.id)}>✕</button>
                            </li>
                          ))}
                        </ul>
                      ) : (
                        <p className="no-options-text">Нет опций</p>
                      )}

                      <form className="add-option-form" onSubmit={(e) => handleAddOptionToType(e, rt.id)}>
                        {addOptionError && <p className="error">{addOptionError}</p>}
                        <input
                          placeholder="Название опции"
                          value={addOptionForm.name}
                          onChange={(e) => setAddOptionForm({ ...addOptionForm, name: e.target.value })}
                          required
                        />
                        <select
                          value={addOptionForm.option_type}
                          onChange={(e) => setAddOptionForm({ ...addOptionForm, option_type: e.target.value as "text" | "number" | "boolean" })}
                        >
                          <option value="text">Текст</option>
                          <option value="number">Число</option>
                          <option value="boolean">Да/Нет</option>
                        </select>
                        <label className="inline">
                          <input
                            type="checkbox"
                            checked={addOptionForm.is_required}
                            onChange={(e) => setAddOptionForm({ ...addOptionForm, is_required: e.target.checked })}
                          />
                          Обязательное
                        </label>
                        <button type="submit" className="btn btn-primary btn-sm">Добавить опцию</button>
                      </form>
                    </div>
                  )}
                </li>
              ))}
            </ul>

            <form className="type-add-form" onSubmit={handleCreateType}>
              <h4>Новый тип ресурса</h4>
              <input
                placeholder="Название типа (напр. PARKING)"
                value={newTypeName}
                onChange={(e) => setNewTypeName(e.target.value)}
                required
              />

              {newTypeOptions.length > 0 && (
                <div className="new-type-options">
                  <h5>Опции нового типа</h5>
                  {newTypeOptions.map((opt, idx) => (
                    <div key={idx} className="new-option-row">
                      <input
                        placeholder="Название опции"
                        value={opt.name}
                        onChange={(e) => handleChangeNewTypeOption(idx, "name", e.target.value)}
                        required
                      />
                      <select
                        value={opt.option_type}
                        onChange={(e) => handleChangeNewTypeOption(idx, "option_type", e.target.value)}
                      >
                        <option value="text">Текст</option>
                        <option value="number">Число</option>
                        <option value="boolean">Да/Нет</option>
                      </select>
                      <label className="inline">
                        <input
                          type="checkbox"
                          checked={opt.is_required}
                          onChange={(e) => handleChangeNewTypeOption(idx, "is_required", e.target.checked)}
                        />
                        Обяз.
                      </label>
                      <button type="button" className="btn btn-danger btn-xs" onClick={() => handleRemoveNewTypeOption(idx)}>✕</button>
                    </div>
                  ))}
                </div>
              )}

              <div className="type-add-actions">
                <button type="button" className="btn btn-secondary btn-sm" onClick={handleAddNewTypeOption}>
                  + Опция
                </button>
                <button type="submit" className="btn btn-primary">Создать тип</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </section>
  );
}
