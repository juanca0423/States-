---
name: gorm-error-chaining
description: Avoid the GORM `.Error` chaining trap that makes error checks always true and breaks `go vet` (func value, not called).
source: auto-skill
extracted_at: '2026-07-12T15:44:29.613Z'
---

## Problem
In this project (`ef/db`, GORM v2), code like:

```go
if mErr := DB.AutoMigrate(&models.Transaccion{}).Error; mErr != nil {
    fmt.Printf("❌ ...: %v\n", mErr)
}
```

compiles but is **buggy**. `go vet` reports:

```
db/connection.go:85:70: fmt.Printf format %v arg mErr is a func value, not called
db/connection.go:85:9:  comparison of function Error != nil is always true
```

Consequences:
* `mErr` is a **method value** (`func() string`), not an error, so `mErr != nil` is **always true** → the error branch is taken unconditionally and the `else` (real error recovery) never runs.
* `fmt.Printf("%v", mErr)` prints a func pointer, not a message.

## Why it happens
`(*gorm.DB).AutoMigrate(dst ...interface{}) error` returns `error` **directly**. Chaining `.Error` off that returned `error` interface selects the interface's `Error() string` method (a method value), which is a function. By contrast, `(*gorm.DB).Exec(...)` returns `*gorm.DB`, whose `.Error` is a real `error` **field** — so `DB.Exec(...).Error` is fine.

Rule of thumb: **if the GORM call returns `error` directly, do NOT chain `.Error`; if it returns `*gorm.DB`, use `.Error`.**

| Call | Returns | Use |
|------|---------|-----|
| `db.AutoMigrate(...)` | `error` | `if err := db.AutoMigrate(...); err != nil` |
| `db.Create(...)` | `*gorm.DB` | `if err := db.Create(...).Error; err != nil` |
| `db.Exec(...)` | `*gorm.DB` | `if err := db.Exec(...).Error; err != nil` |
| `db.First(...)` | `*gorm.DB` | `if err := db.First(...).Error; err != nil` |

## Fix procedure
1. Run `go vet ./db/...` (or `go doc gorm.io/gorm DB.<Method>` to confirm the return type).
2. If the method returns `error` directly, capture it immediately:
   ```go
   if err := DB.AutoMigrate(&models.Transaccion{}); err != nil {
       fmt.Printf("❌ Error en AutoMigrate tras eliminar políticas: %v\n", err)
   } else {
       // recovery path
   }
   ```
3. Never read `.Error` off an `error` value — only off a `*gorm.DB` value.

## How to recognize this class of bug
`go vet` "arg X is a func value, not called" on a `fmt.Printf`/`Sprintf` argument almost always means the argument is a **method value** (a function) instead of a value. Trace the receiver's return type: if it returns an interface that itself has a method of the same name (e.g. `error.Error()`), the selector yields the method, not a field. Call it (`()`) or, better, restructure to capture the returned error directly.

## Why this matters here
* CI runs `go test ./... -v`, which runs `go vet` on every package; **one vet error in any package fails the whole suite** (observed: `ef/db` broke `go test ./...` despite the rest passing).
* The bug is silent at runtime (no compile error) but silently disables error recovery, so it must be caught via `go vet`.

## Verification
* `go vet ./db/...` → clean
* `go build ./...` → exit 0
* `go test ./...` → all packages `ok`
