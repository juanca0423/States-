---
name: subscription-gating
description: How to enforce trial/subscription expiry in the States Go+Fiber app by wiring the CheckSubscription middleware and fixing the date logic
source: auto-skill
extracted_at: '2026-07-12T12:22:46.281Z'
---

# Enforcing trial/subscription expiry in States

## Context
The app creates users with a 30-day trial (`FechaFinPrueba = time.Now().AddDate(0,0,30)`,
`SuscripcionActiva = false`) in `ctrl/registro_ctrl.go`. A `CheckSubscription` middleware
exists in `middleware/ahut.go` that should block access once the trial expires, but:

1. **It was never wired into the routes** — `rutas/index_router.go` only applied
   `middleware.AuthRequired` to protected routes, so expired users kept full access.
2. **The date comparison had an off-by-one bug** — `int(time.Until(finPrueba).Hours()/24)`
   truncates toward zero, so for the first ~24h after expiry `diasRestantes == 0`, granting
   access. `ahora.Before(finPrueba)` alone was the correct guard but was OR'd with the buggy term.

## Procedure
1. In `rutas/index_router.go`, add `middleware.CheckSubscription` right after
   `middleware.AuthRequired` on the routes that must enforce the trial:
   `/eeff`, `/estados`, `/costos`, `/costosform`.
   - Do NOT add it to `/perfil` (shows plan status — the user explicitly opted out of gating it)
     nor `/planes` (it is the redirect target — adding it there causes a redirect loop).
   - Admin role bypasses inside `CheckSubscription` itself (`if role == "admin"`), so
     `/api/admin/*` is covered without extra wiring.
   - Grep to confirm wiring: a missing `CheckSubscription` on a protected route is the #1
     cause of "expired users still have access" reports.
2. In `middleware/ahut.go`, replace the buggy block with a pure instant comparison:
   ```go
   if suscrito {
       return c.Next()
   }
   if !time.Now().After(finPrueba) {
       return c.Next()
   }
   return c.Redirect("/planes?reason=expired")
   ```
   This grants access while now <= finPrueba and redirects once it has passed.
3. Display-side day math (NOT security) also truncates toward zero: in `AuthStatus`
   (`dias := int(time.Until(...).Hours()/24)`) and controllers use
   `max(0, int(time.Until(...).Hours()/24))` to avoid showing negative/“0 días” after expiry.
   A generic `max[T constraints.Ordered]` helper lives in `utils/helpers.go` (added via
   `golang.org/x/exp/constraints`); reuse it instead of importing `math`.

## Notes
- `CheckSubscription` reads `Suscrito`/`FinPrueba` from `c.Locals` (set by the global
  `AuthStatus` middleware) and falls back to a DB lookup if missing.
- The trial is **30 days** (`AddDate(0,0,30)` in `registro_ctrl.go`); admins get 1 year
  (`admin_ctrl.go`), seed users 99 years (`db/seed.go`). Keep these in sync if you change terms.
- Price display strings live in views and must stay consistent when changing currency/price:
  `views/planes.hbs` ("Q 300.00 por año"), `views/pago/procesar.hbs` ("Pagar Q 300.00/año"),
  `views/perfil.hbs` ("Q 300.00/año"). The annual plan is Q 300.00 (Guatemalan quetzales).
- `perfil.hbs` progress bar divides `DiasRestantes` by `TotalDias` (30 for trial, 365 for paid),
  passed from `perfil_ctrl.go` via `TextoPlan` — do not hardcode 365 there.
- Run `go build ./...` (and `go mod tidy` if you add the `x/exp` import) to confirm compilation.
