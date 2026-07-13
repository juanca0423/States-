---
name: unify-reccol
description: Consolidate duplicated RecCol* functions into a single reusable core with wrappers.
source: auto-skill
extracted_at: '2026-07-12T15:30:00.000Z'
---

## Goal
Reduce code duplication in `help/resultados_help.go` where three functions (`RecCol1`, `RecCol1Tot`, `RecCol2Tot`) share the same looping logic and only differ in how the third column (`Col3`) is calculated.

## Approach
1. **Identify common pattern** – each function:
   * Iterates over a slice of `models.Cue`.
   * Skips zero‑valued entries.
   * Accumulates costs into the provided total pointers.
   * Builds a slice of `models.KR` with `Key`, `Value.Col1`, optional `Col2`, `Col3`, and border classes.
2. **Extract core** – create a private helper `recColGrupo` that receives:
   * The original parameters (`Balance`, `Recorido`, `titulo`).
   * `sumaAnterior` (previous subtotal) for the wrappers that need it.
   * A *callback* `col3 func(anterior, grupo float64) (float64, bool)` that decides whether and how to produce `Col3`.
   * The three total‑pointer arguments.
   The helper performs the loop, fills `res`, computes `SaldoGrupo`, adds the title row, and applies the callback to optionally set `Col3` and its class.
3. **Define thin wrappers** that preserve the original public signatures (so existing call‑sites do not change):
   * `RecCol1` → passes `nil` for the callback (no `Col3`).
   * `RecCol1Tot` → passes a callback that returns `anterior + grupo` and `true`.
   * `RecCol2Tot` → passes a callback that returns `anterior - grupo` and `true`.
4. **Update imports** – the new helper lives in the same file, so no extra imports are needed.
5. **Add regression test** (`help/resultados_help_test.go`):
   * Build a small `Balance` map and a slice of `models.Cue` that covers a few codes.
   * Call the original functions (now wrappers) and the new helper directly, asserting that the returned `[]models.KR` and subtotals are identical.
   * This ensures the refactor does not alter financial output.
6. **Run `go test ./...`** to verify all existing tests still pass and the new test succeeds.

## Why this works
* The core loop is **single‑source‑of‑truth**, eliminating the risk of divergent behaviour when a future change is required (e.g., a new cost‑accumulation rule).
* The callback isolates the only variable part (the `Col3` calculation), keeping the public API untouched – no other file needs to be modified.
* Adding a test gives confidence that the financial numbers stay correct, satisfying the high‑accuracy requirements of the application.

## Expected impact
* ~60 lines of duplicated code become ~30 lines of core logic + three lightweight wrappers.
* Maintenance effort drops dramatically; any bug fix in the loop is applied once.
* Future extensions (e.g., a fourth variant) only need a new wrapper with a different callback.

---

**Next steps**
1. Apply the patch shown in the conversation (the `recColGrupo` implementation and wrapper updates).
2. Add the regression test file.
3. Run the full test suite.
4. Merge the change once the CI passes.
