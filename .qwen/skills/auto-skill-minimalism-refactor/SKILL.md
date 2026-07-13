---
name: minimalism-refactor
description: Identify and remove redundant code and simplify error‑handling flows while preserving security in Go web projects.
source: auto-skill
extracted_at: '2026-07-12T10:27:53.411Z'
---

## Purpose
This skill captures a reusable pattern for applying the **minimalism** design principle to a Go/Fiber codebase:

1. **Detect duplicated logic** (e.g., repeated DB updates, repeated error handling).
2. **Extract shared functionality** into a small helper function that returns an `error`.
3. **Replace the duplicated blocks** with a single call to the helper.
4. **Remove unnecessary GORM session options** (`PrepareStmt: false`) when they do not affect the query.
5. **Add explicit error checks** for DB look‑ups that were previously ignored, redirecting to a safe endpoint (e.g., login) on failure.
6. **Keep security guarantees** – do not change authentication, JWT handling, or data validation.

## Step‑by‑step procedure
1. **Search for repeated code** using a simple text search (e.g., `Updates(map[string]any{` or `Session(&gorm.Session{PrepareStmt: false})`).
2. **Create a helper** in the same package:
   ```go
   func responderConsulta(id, respuesta string) error {
       return db.DB.Model(&models.Mensaje{}).Where("id = ?", id).
           Updates(map[string]any{"respuesta": respuesta, "estado": "Resuelto"}).Error
   }
   ```
3. **Replace each duplicated block** with a call to the helper and handle the returned error appropriately.
4. **Remove unnecessary GORM session wrappers** from simple queries:
   ```go
   // Before
   db.DB.Session(&gorm.Session{PrepareStmt: false}).Where("email = ?", email).First(&u)
   // After
   db.DB.Where("email = ?", email).First(&u)
   ```
5. **Add explicit error handling** for look‑ups that previously ignored errors:
   ```go
   if err := db.DB.First(&u, uid).Error; err != nil {
       return c.Redirect("/loguin")
   }
   ```
6. **Run `go vet` / `go test`** to ensure no behavioural changes were introduced.

## Removing dead code (Go)
* A handler function is dead if it is **not referenced by any route** in `rutas/index_router.go`
  (and not used by tests). Example removed this session: `GetDashboardInfo` in `ctrl/user_ctrl.go`
  was unreachable (route `/eeff` uses `ctrl.HojaTrabajo` in `hoja_ctrl.go`), so the whole
  `user_ctrl.go` file was deleted once emptied.
* After deleting a function, prune its now-unused imports (or the whole file if empty).
* Verify with `grep_search` for the symbol name + `go build ./...`.

## Cleaning unused static assets (`estaticos/`)
* List `estaticos/js/*.js` (and `estaticos/css/*`) and grep each filename across the repo —
  **not only `views/`**. The static mount is `app.Static("/static", "./estaticos")`, so valid
  references are `/static/js/...`. Watch for the **broken `/estaticos/js/...` path** (note the
  missing `t`): any tag using `/estaticos/` 404s and is dead weight.
* Libs loaded via CDN (e.g. jsPDF, exceljs, chart.js from cdnjs) mean the local copies
  (`jspdf*.js`, etc.) are redundant — delete the local copies.
* Keep files the user explicitly names as exceptions even if they look unused.
* When deleting a JS file, also remove any `<script src=...>` tags that referenced it
  (e.g. the broken bootstrap/popper tags in `views/layouts/main.hbs`). `bootstrap.bundle.min.js`
  already supplies Popper, so the standalone `popper.min.js` and `bootstrap.min.js` are redundant.
* Deletion is destructive but recoverable via git; confirm with the user before `del` on Windows.

## Cleaning duplicated Handlebars views (`views/`)
* When two or more `.hbs` templates are nearly identical (differ only by a form `action`, a
  heading, or a passed variable), extract the shared markup into a `partials/` file and include
  it with `{{> 'partials/nombre'}}`.
* Real example (this session): `views/eeffform.hbs` and `views/costosform.hbs` were ~99% identical
  — the only real difference was the form `action` (`/estados` vs `/costos`). A floating calculator
  widget (~140 lines of HTML + inline CSS) and a shortcuts bar (~8 lines) were copy‑pasted into both.
  - Calculator widget + its `<script src="/static/js/calculadora.js"></script>` → `views/partials/calculadora.hbs`.
  - Shortcuts bar → `views/partials/shortcuts.hbs`.
  - Both forms now `{{> 'partials/calculadora'}}` and `{{> 'partials/shortcuts'}}`; each dropped
    from ~200 to ~40 lines, single source of truth.
* Keep template‑specific bits (form `action`, `{{#if Error}}` block, `window.CATALOGO_CUENTAS = {{{Cuentas}}}`,
  the React bundle `<script src="/static/js/bundle-formulario.js">`) inside the template; only pull
  out the truly identical blocks.
* After extracting, `grep_search` for a unique string from the old block (e.g. `cuadroCuadre`) to
  confirm it now appears **only** in the partial, not in the form templates.
* Templates are parsed at runtime, so `go build ./...` will pass even with a broken include — verify
  the partial file (`partials/<name>.hbs`) exists to avoid a 500 at render time.
* Do NOT touch shared JS (`bundle-formulario.js`, `calculadora.js`) — those are already single shared
  files; the duplication was in the HTML, not the JS.
* Project rule (QWEN.md): avoid unnecessary inline CSS. While extracting, the widget's heavy
  `style="..."` blocks could later be moved to `style.css`, but that is a separate cleanup step.

## When to apply
Use this skill whenever you encounter:
* Repeated DB update blocks across handlers.
* GORM queries wrapped in `Session{PrepareStmt: false}` without a performance reason.
* Silent DB look‑ups where the error is ignored.
* Handler functions not wired to any route (dead code).
* Static `.js`/`.css` files with no matching `<script>`/`<link>` reference (or only a broken
  `/estaticos/` path, or a CDN-loaded duplicate).
* Any code that can be expressed more concisely without losing clarity or security.

## Rationale
* **Reduced code size** aligns with the project’s “minimalismo” guideline.
* **Single source of truth** for DB update logic makes future changes easier and less error‑prone.
* **Explicit error handling** prevents hidden failures and improves reliability.
* **No functional impact** – the same HTTP responses, security checks, and data validations remain.

---

*Created automatically by the Qwen auto‑skill extraction system on 2026‑07‑12.*