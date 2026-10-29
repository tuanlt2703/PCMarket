# Administator Web

Vue 3 + Vite admin site built on the AdminLTE 4 template in `admin-lte-v4.10.1/dist`.

```bash
npm install
npm run dev      # http://localhost:5173  (mock login: admin / admin123)
npm run build
```

Data comes from an in-browser mock API (`src/services/*Api.js`, persisted to localStorage; reset via the user menu).
Product spec forms are generated from `src/config/componentSchemas.js`, which mirrors the product-db migrations.