# LFE Go SDK API Contract Matrix

Status: RC Contract

This document defines the accepted public Go SDK usage surface for the current
release-candidate line. It is a usage contract, not a description of Engine
internals. Historical decision records do not override this contract.

The central identity used by the SDK is the source-owned sequence coordinate
(`seq`). Applications own source records and payloads; LFE associates logical
state with those source coordinates and returns candidate sets that preserve the
same identity.

## Contract reading guide

Each API is described by five questions:

- **Purpose** — why an application calls it.
- **Input** — what the caller supplies.
- **Output** — what the caller receives.
- **Semantic guarantee** — behavior callers may rely on.
- **Typical use** — where the API normally belongs in an application flow.

APIs returning `*Engine`, `*SeqSet`, or `*SeqSetEx` transfer ownership of that
returned handle to the caller. The caller MUST call `Close()` when the object is
no longer needed. The current Go implementation also installs finalizers as a
safety net; finalizers are not a substitute for deterministic `Close()` calls.

---

## 1. Engine lifecycle

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `New()` | Open an Engine using normal SDK license discovery. | None. | `(*Engine, error)` | Returns a usable Engine only when discovery and license verification succeed. | Normal application startup. |
| `NewWithLicensePath(path)` | Open an Engine with an explicit license artifact. | Filesystem path to a license artifact. | `(*Engine, error)` | The explicit path is authoritative for this call; failure is returned rather than silently treating the Engine as licensed. | Tests, controlled deployments, explicit runtime configuration. |
| `(*Engine).Close()` | Release Engine resources. | Engine handle. | None. | Safe to call when the Engine is no longer needed; the handle must not be used afterward. | `defer engine.Close()` immediately after successful construction. |

### Lifecycle example

```go
engine, err := lfe.New()
if err != nil {
    return err
}
defer engine.Close()
```

---

## 2. Logical definition

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `DefineUInt(id, name, bits)` | Declare an unsigned integer logical projection. | Projection ID, name, declared bit width. | `error` | Subsequent unsigned ingest/query calls may address the projection by the same ID. | Scores, counters, amounts, categories encoded as unsigned values. |
| `DefineInt(id, name, bits)` | Declare a signed integer logical projection. | Projection ID, name, declared bit width. | `error` | Subsequent signed ingest/query calls may address the projection by the same ID. | Signed deltas, temperatures, balances. |
| `DefineDateTime(id, name)` | Declare a date/time logical projection. | Projection ID and name. | `error` | Date/time values are queried through typed date/time operands. | Event time, update time, business time dimensions. |
| `DefineFlagSet(id, name, count)` | Declare a logical flag collection. | Projection ID, name, logical flag count. | `error` | Individual flags may later be queried by index and expected boolean state. | Status capabilities, feature flags, eligibility markers. |

Projection IDs are application-level identifiers within an Engine. Applications
SHOULD keep their definition mapping stable across persistence/restore and across
components that are expected to address the same logical projection.

---

## 3. Single-record mutation

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `IngestUInt(seq, id, value)` | Add unsigned logical state for a source coordinate. | Source-owned `seq`, projection ID, unsigned value. | `error` | Associates the supplied logical value with the supplied source coordinate. | Initial ingest or full synchronization. |
| `IngestInt(seq, id, value)` | Add signed logical state. | `seq`, projection ID, signed value. | `error` | Signed value semantics are preserved by the typed route. | Initial ingest of signed data. |
| `IngestDateTime(seq, id, value)` | Add date/time logical state. | `seq`, projection ID, `DateTime`. | `error` | Date/time components remain queryable through typed operands. | Event or source timestamps. |
| `IngestFlagSet(seq, id, value)` | Add flag logical state for a source coordinate. | `seq`, projection ID, `FlagSetValue`. | `error` | The supplied packed logical value becomes the flag state for that coordinate. Applications supply stable flag positions and raw `0/1` values through `FlagSet`; the SDK owns packing. | Eligibility/status ingestion. |
| `UpdateUInt(seq, id, value)` | Change existing unsigned logical state. | `seq`, projection ID, new value. | `error` | Subsequent resolves observe the accepted update without rebuilding the Engine. | Live event processing and mutable source state. |
| `UpdateInt(seq, id, value)` | Change existing signed logical state. | `seq`, projection ID, new signed value. | `error` | Same typed signed semantics as ingest. | Live signed-value mutation. |
| `UpdateDateTime(seq, id, value)` | Change existing date/time logical state. | `seq`, projection ID, new `DateTime`. | `error` | Subsequent typed date/time resolves observe the accepted update. | Moving timestamps/status times. |
| `UpdateFlagSet(seq, id, value)` | Replace/update flag state. | `seq`, projection ID, `FlagSetValue`. | `error` | Subsequent flag resolves observe the accepted packed logical state. Applications supply stable flag positions and raw `0/1` values through `FlagSet`; the SDK owns packing. | Live eligibility/status changes. |
| `Delete(seq)` | Remove logical state for a source coordinate from the Engine. | Source-owned `seq`. | `error` | The deleted source coordinate no longer participates as live Engine state after a successful delete. | Source deletion, tombstone processing, retention flows. |

### Live mutation contract

An Engine may be kept resident while an application alternates mutation and
resolution:

```text
update batch -> resolve -> update batch -> resolve -> ...
```

This contract does **not** by itself claim concurrent thread safety. Concurrency
requirements must be defined separately from interleaved single-process usage.

---

## 4. Batch ingest

### Types

```go
type AddRecord struct {
    Seq          uint64
    ProjectionID uint32
    Value        Value
}

type BatchStats struct {
    Records     uint64
    Segments    uint64
    WorkersUsed uint64
}
```

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `AddBatch(records)` | Send heterogeneous logical values through one SDK batch call. | `[]AddRecord` containing `UIntValue`, `IntValue`, `DateTime`, or `FlagSetValue` values. | `(BatchStats, error)` | Numeric, DateTime, and FlagSet logical values may be mixed in one batch while preserving the same source-coordinate and projection semantics as single-record ingest. | Initial load, batch synchronization, high-throughput import. |

`AddBatch` is the heterogeneous batch-ingest contract. Numeric, DateTime, and
FlagSet logical values may be mixed in one call. Numeric batch values are
constructed with `UIntValue` or `IntValue`; `DateTime` values may be created
directly or through `DateTimeUTC` / `DateTimeUTCString`; flag values are created
with `FlagSet` and may be composed with `Flags`. Applications own flag meaning
and stable flag positions; the SDK owns shifting and packing.

`BatchStats` is operational information about the accepted call. Applications
SHOULD NOT build logical correctness rules from `WorkersUsed` or `Segments`.

---

## 5. Query construction

### Operators

```go
Eq
Neq
Gt
Gte
Lt
Lte
```

### Typed operands

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `UInt(value)` | Build an unsigned numeric operand. | `uint64`. | `Operand` | Selects unsigned numeric query semantics. | Queries against `DefineUInt`. |
| `Int(value)` | Build a signed numeric operand. | `int64`. | `Operand` | Selects signed numeric query semantics. | Queries against `DefineInt`. |
| `YearOperand(value)` | Build a year component operand. | `uint16`. | `Operand` | Selects the year component of a date/time definition. | `Eq`, range operators on year. |
| `MonthOperand(value)` | Build a month component operand. | `uint8`. | `Operand` | Selects month. | Date/time projection queries. |
| `DayOperand(value)` | Build a day component operand. | `uint8`. | `Operand` | Selects day. | Date/time projection queries. |
| `HourOperand(value)` | Build an hour component operand. | `uint8`. | `Operand` | Selects hour. | Hour-of-day logical decisions. |
| `MinuteOperand(value)` | Build a minute component operand. | `uint8`. | `Operand` | Selects minute. | Date/time projection queries. |
| `SecondOperand(value)` | Build a second component operand. | `uint8`. | `Operand` | Selects second. | Date/time projection queries. |
| `Flag(index, expected)` | Build a typed flag operand. | Flag index and expected boolean state. | `Operand` | Tests the selected flag against the expected state. | Eligibility/status rules. |

### Query API

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `NewQuery(id, op, operand)` | Construct one logical predicate. | Projection ID, operator, typed operand. | `Query` | The operand route determines the value semantics; callers should use the operand matching the definition type. | All resolve calls. |
| `Query.SeqFrom(seq)` | Add a lower source-coordinate bound. | `uint64`. | New `Query` value. | The query is bounded from the supplied source coordinate. | Windowed processing. |
| `Query.SeqTo(seq)` | Add an upper source-coordinate bound. | `uint64`. | New `Query` value. | The query is bounded to the supplied source coordinate. | Windowed processing. |
| `Query.Limit(limit)` | Bound result cardinality. | `uint64`. | New `Query` value. | The returned query carries the requested result limit. In candidate-domain resolve, the limit applies after candidate intersection. | Bounded consumption, top-N-by-coordinate style workflows. |

`Query` builder methods return a value and do not mutate the previously held
`Query` variable unless the caller assigns the returned value.

---

## 6. Resolve surfaces: choosing the result type

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `Resolve(query)` | Convenience resolve when the caller explicitly wants all matching source coordinates as a Go slice. | `Query`. | `([]uint64, error)` | Materializes the complete resolved `SeqSet` into Go memory. | Small/known-bounded results, compatibility code, simple application loops. |
| `ResolveSet(query)` | Resolve without eagerly materializing the complete result. | `Query`. | `(*SeqSet, error)` | Returns a native-backed result that can be inspected, bounded-materialized, or streamed. | Generic application consumption. |
| `ResolveEx(query)` | Resolve to the semantic interoperability/pipeline result. | `Query`. | `(*SeqSetEx, error)` | Returns semantic membership without requiring `seq[]` materialization or binary transport. | Pipelines, cross-engine composition, set algebra, downstream adapters. |

### Selection rule

```text
Need []uint64 now?                  -> Resolve
Need bounded iteration/streaming?   -> ResolveSet
Need to keep computing on the set?  -> ResolveEx
```

---

## 7. Candidate-domain pipeline composition

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `ResolveFromSet(input, query)` | Apply a query only within an existing semantic candidate domain and return a stream/materialization-oriented result. | Non-nil `*SeqSetEx`, `Query`. | `(*SeqSet, error)` | Every result member was already present in `input` and also matches the query in the receiving Engine. Query bounds/limit remain applicable; limit is applied after candidate intersection. | End a semantic pipeline in a `SeqSet` consumer. |
| `ResolveFromSetEx(input, query)` | Apply a query within an existing candidate domain and keep the result semantic. | Non-nil `*SeqSetEx`, `Query`. | `(*SeqSetEx, error)` | `result` is a subset of `input`. No binary transport or `seq[]` materialization is required between stages. | Multi-stage and cross-engine pipelines. |

A `SeqSetEx` may originate from another Engine when both Engines interpret
`seq` using the same source-owned coordinate space. Engines do not need to
contain the same definitions. A coordinate present in the input but lacking the
required logical state in the receiving Engine does not become a match. This is
the accepted Go composition model.

### Cross-engine example

```go
amountCandidates, err := pricing.ResolveEx(
    lfe.NewQuery(Amount, lfe.Gte, lfe.UInt(900_000)),
)
if err != nil {
    return err
}
defer amountCandidates.Close()

final, err := risk.ResolveFromSetEx(
    amountCandidates,
    lfe.NewQuery(Risk, lfe.Gte, lfe.UInt(90)),
)
if err != nil {
    return err
}
defer final.Close()
```

---

## 8. `SeqSet` result contract

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `Close()` | Release result resources. | `*SeqSet`. | None. | Result must not be used afterward. | Always `defer set.Close()` after success. |
| `Len()` | Read cardinality. | None. | `uint64`. | Returns zero for nil/closed handles in the current Go surface. | Fast result sizing. |
| `IsEmpty()` | Test emptiness. | None. | `bool`. | Reports whether the result has no members; nil/closed is treated as empty by the current Go wrapper. | Guard downstream work. |
| `Contains(seq)` | Test membership. | Source `seq`. | `bool`. | Does not materialize the whole result. | Membership checks. |
| `Min()` | Read minimum source coordinate. | None. | `(uint64, bool)`. | `bool=false` when no minimum is available. | Range inspection. |
| `Max()` | Read maximum source coordinate. | None. | `(uint64, bool)`. | `bool=false` when no maximum is available. | Range inspection. |
| `Materialize(max)` | Copy at most the requested materialization bound to a Go slice. | Maximum number of coordinates requested. | `([]uint64, error)` | Produces explicit source coordinates; callers pay materialization cost here. | Bounded handoff to APIs requiring slices. |
| `FetchNext(lastSeq, limit)` | Fetch the next ascending bounded page. | Optional previous cursor, limit `1..256`. | Chunk, continuation cursor, `hasNext`, error. | Returned chunks are bounded; the returned cursor is supplied to the next call when continuation exists. | Pull-based iteration. |
| `StreamChunks(size, fn)` | Consume the result through bounded ascending chunks. | Chunk size `1..256`, callback. | `error` | Does not require complete result materialization; callback failure stops the stream and is returned. | Database/API batches, streaming consumers. |

The current bounded chunk capacity is 256 source coordinates.

---

## 9. `SeqSetEx` semantic result contract

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `Close()` | Release semantic result resources. | `*SeqSetEx`. | None. | Result must not be used afterward. | Always close owned results. |
| `Len()` | Read semantic cardinality. | None. | `uint64`. | Does not require `seq[]` materialization. | Pipeline sizing/metrics. |
| `IsEmpty()` | Test semantic emptiness. | None. | `bool`. | Does not require transport encoding. | Guard composition/downstream work. |
| `Contains(seq)` | Test source-coordinate membership. | Source `seq`. | `bool`. | Membership is evaluated against the semantic set. | Correctness checks, application decisions. |
| `Min()` | Read minimum member. | None. | `(uint64, bool)`. | `bool=false` when unavailable/empty. | Range inspection. |
| `Max()` | Read maximum member. | None. | `(uint64, bool)`. | `bool=false` when unavailable/empty. | Range inspection. |
| `Binary()` | Request canonical adapter transport bytes. | None. | `([]byte, error)` | Produces `LFESEQ01` only when explicitly requested; binary transport is not required for in-process semantic composition. | Process/system/adapter boundary. |

`SeqSetEx` should be treated as the working logical set when additional LFE
operations are expected. `Binary()` is a boundary operation, not a prerequisite
for further computation.

---

## 10. `SeqSetEx` set algebra

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `Merge(a, b)` | Combine independent candidate regions. | Two non-nil `*SeqSetEx`. | `(*SeqSetEx, error)` | Semantic result is the union `A ∪ B`; duplicates are represented once. | Independent rule regions, fan-in before downstream work. |
| `Intersect(a, b)` | Keep membership shared by both sets. | Two non-nil `*SeqSetEx`. | `(*SeqSetEx, error)` | Semantic result is `A ∩ B`. | Shared eligibility, multiple independent constraints. |
| `Difference(a, b)` | Exclude members of B from A. | Two non-nil `*SeqSetEx`. | `(*SeqSetEx, error)` | Directional semantic result is `A - B`. | Exclusion lists, already-processed/known-safe removal. |

All three operations return a first-class `SeqSetEx`. The result may be fed
into `ResolveFromSetEx` or encoded later with `Binary()`. The accepted test
surface verifies both exact set algebra and reuse of a merged result as the
candidate input to another Engine.

### Example

```go
combined, err := lfe.Merge(regionA, regionB)
if err != nil {
    return err
}
defer combined.Close()

eligible, err := lfe.Intersect(combined, activeCustomers)
if err != nil {
    return err
}
defer eligible.Close()

final, err := lfe.Difference(eligible, knownExclusions)
if err != nil {
    return err
}
defer final.Close()
```

---

## 11. Persistence

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `Persist(root)` | Persist Engine state to the configured filesystem root. | Root path. | `error` | Successful return indicates the persistence operation completed for this Engine snapshot. | Controlled shutdown, warm restart, durable projection state. |
| `Restore(root)` | Restore previously persisted Engine state. | Root path. | `error` | A successful restore enters the correctness lifecycle required by the SDK before normal mutation resumes. | Application restart/recovery. |

Persistence does not transfer ownership of source payloads to LFE. The source
system remains authoritative for source records.

---

## 12. Correctness lifecycle (advanced)

### States

```go
CorrectnessReady
CorrectnessPending
CorrectnessDrifted
CorrectnessFullSync
```

| API | Purpose | Input | Output | Semantic guarantee | Typical use |
|---|---|---|---|---|---|
| `CorrectnessGateState()` | Inspect current restore/synchronization state. | None. | `(CorrectnessGateState, error)` | Returns the SDK lifecycle state used to gate mutation after restore or detected mismatch. | Recovery orchestration and health checks. |
| `ConfirmCorrectnessGate(correct)` | Confirm whether restored logical state matches the authoritative source. | Boolean confirmation. | `error` | Accepted confirmation moves the lifecycle according to the correctness state machine; invalid transitions return an error. | Post-restore source validation. |
| `ForceFullSync()` | Enter full synchronization after drift/mismatch. | None. | `error` | Only valid from the lifecycle state that requires a full synchronization. | Rebuild logical state from authoritative source. |
| `CompleteFullSync()` | Mark full synchronization complete. | None. | `error` | Returns the Engine to normal ready operation only after the valid full-sync lifecycle. | Finish recovery. |

The correctness lifecycle is an advanced operational surface. Applications that
use persistence MUST document how they validate restored logical state against
the authoritative source before reopening normal mutation.

---

## 13. Public data types

### `DateTime`

```go
type DateTime struct {
    Year   uint16
    Month  uint8
    Day    uint8
    Hour   uint8
    Minute uint8
    Second uint8
}
```

Applications SHOULD normally create UTC values through the SDK helpers rather
than manually decomposing a `time.Time`:

```go
value := lfe.DateTimeUTC(time.Now())

parsed, err := lfe.DateTimeUTCString("2026-08-30 12:23:45")
if err != nil {
    return err
}
```

`DateTimeUTC` normalizes the supplied `time.Time` to UTC and projects year,
month, day, hour, minute, and second. `DateTimeUTCString` accepts the fixed
`YYYY-MM-DD HH:MM:SS` format and interprets it as UTC.

### `FlagSetValue`

Applications own flag meaning and stable zero-based positions. Raw values passed
to `FlagSet(position, raw)` MUST be `0` or `1`; the SDK owns shifting and
packing. Multiple values may be composed with `Flags`:

```go
const (
    Active uint8 = iota
    Verified
    Fraud
    Premium
)

status := lfe.Flags(
    lfe.FlagSet(Active, active),
    lfe.FlagSet(Verified, verified),
    lfe.FlagSet(Fraud, fraud),
    lfe.FlagSet(Premium, premium),
)
```

Resolve addresses the stable flag position rather than the application's
business name:

```go
q := lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true))
```

### `Operator`

```go
type Operator uint8
```

Callers SHOULD use the exported operator constants rather than relying on their
numeric values.

### `Query`

`Query` is publicly inspectable in Go, but callers SHOULD prefer `NewQuery` plus
`SeqFrom`, `SeqTo`, and `Limit` rather than constructing internal operand state
manually. `Operand` values SHOULD be created through the exported typed operand
constructors.

### `Error`

SDK operations return normal Go `error` values. Callers MAY inspect `*lfe.Error`
when SDK-specific operation/code information is needed, but application logic
SHOULD prefer the documented semantic condition over depending on undocumented
numeric error codes.

---

## 14. Recommended application patterns

### A. Heterogeneous source-row ingest

One source coordinate may carry multiple logical projections. The canonical
batch shape keeps the source `seq` unchanged and varies `ProjectionID`:

```go
rows := []lfe.AddRecord{
    {Seq: seq, ProjectionID: Amount, Value: lfe.UIntValue(amount)},
    {Seq: seq, ProjectionID: RiskScore, Value: lfe.IntValue(risk)},
    {Seq: seq, ProjectionID: EventTime, Value: lfe.DateTimeUTC(eventTime)},
    {Seq: seq, ProjectionID: Status, Value: statusFlags},
}

_, err := engine.AddBatch(rows)
```

A repeated `seq` in this form is not an identity collision: each record targets
a different logical projection for the same source coordinate.

### B. Source-backed logical projection

```text
source payload ----------------------> database / source system
      |
      +-- seq + selected values -----> LFE Engine
                                         |
                                         +-> ResolveEx / SeqSetEx
```

LFE does not require ownership of the source payload.

### C. Narrowing pipeline

```text
ResolveEx(Q1)
  -> SeqSetEx
  -> ResolveFromSetEx(Q2)
  -> SeqSetEx
  -> ResolveFromSetEx(Q3)
  -> SeqSetEx
```

Use this when each rule should evaluate only the candidate domain produced by
the previous rule.

A canonical application-shaped flow is:

```text
source transaction
 |-- Amount     UInt
 |-- RiskScore  Int
 |-- EventTime  DateTime
 `-- Status     FlagSet
        |
        v
heterogeneous AddBatch
        |
        v
ResolveEx(amount)
        |
        v
ResolveFromSetEx(risk)
        |
        v
ResolveFromSetEx(flag)
        |
        v
SeqSetEx
```

The repository's `examples/transaction-risk` program executes this pattern and
checks exact membership against a deterministic Go truth oracle.

### D. Independent regions and composition

```text
Engine A -> SeqSetEx A --+
                         +-> Merge -> SeqSetEx -> downstream
Engine B -> SeqSetEx B --+
```

Use `Intersect` for shared membership and `Difference` for exclusions.

### E. Downstream boundary

```text
SeqSetEx
   |
   +-- continue computing in-process -> keep SeqSetEx
   |
   +-- leave SDK/process boundary ----> Binary() -> LFESEQ01
```

### F. Live resident Engine

```text
receive changes
-> Ingest/Update/Delete
-> ResolveEx / ResolveFromSetEx
-> act on candidates
-> repeat
```

The Engine does not need to be rebuilt between accepted mutations and resolves.

---

## 15. Anti-patterns

### Premature materialization

Avoid:

```text
ResolveEx -> []seq -> sort/deduplicate -> set operation
```

Prefer:

```text
ResolveEx -> SeqSetEx -> Merge/Intersect/Difference
```

### Binary transport between in-process stages

Avoid encoding `LFESEQ01` solely to pass a candidate set from one in-process LFE
operation to another. Keep `SeqSetEx` semantic until a real transport/adapter
boundary exists.

### Adapter-local identity replacement

Do not silently replace source-owned `seq` coordinates with local row numbers or
adapter-generated offsets when results are expected to compose across Engines or
systems.

### Assuming set composition requires identical Engine definitions

Cross-engine composition depends on compatible source-owned coordinates, not on
identical projection definitions.

---

## 16. RC review decisions

The following decisions define the accepted RC surface:

1. `Engine` is the single canonical public entry point.
2. `Resolve`, `ResolveSet`, and `ResolveEx` intentionally serve three different
   consumption modes rather than competing aliases.
3. `ResolveFromSet` / `ResolveFromSetEx` are the canonical Go names for
   candidate-domain resolve because Go cannot overload the pre-existing
   `ResolveSet(Query)` method.
4. `SeqSetEx` is the canonical semantic pipeline/interoperability result.
5. `Merge`, `Intersect`, and `Difference` are public set-composition operations
   over `SeqSetEx`.
6. `Binary()` is an explicit downstream transport boundary.
7. `AddBatch` is the heterogeneous batch-ingest contract for this RC.
   Numeric, DateTime, and FlagSet logical values may be mixed in one call
   through the unified `Value` surface.
8. Persistence correctness lifecycle APIs remain public but advanced.
9. Returned Engine/result handles require deterministic `Close()` ownership.
10. No additional public API should be added before RC unless it closes a
    demonstrated usage gap that cannot be expressed by the current surface.
