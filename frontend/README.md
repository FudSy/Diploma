# Frontend (React + Vite)

## Requirements

- Node.js 20+
- npm 10+

## Run

```bash
cd frontend
npm install
npm run dev
```

Vite dev server runs at `http://localhost:5173` and proxies API calls from `/api/*` to `http://localhost:8080`.

## Environment variables

Optional `.env` in `frontend/`:

```bash
VITE_API_BASE_URL=/api
```

For deployment (without Vite proxy), set full backend URL, e.g.:

```bash
VITE_API_BASE_URL=http://localhost:8080
```
