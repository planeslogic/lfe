# LFE BE SDK for Go

The LFE BE SDK for Go embeds the commercial LFE execution engine behind a small Go API. Applications keep ownership of source records and payloads; LFE associates logical projections with the application's source-owned sequence coordinate (`seq`) and returns candidate sets in that same coordinate space.

The Go package provides the Go SDK surface for the LFE Engine. Engine semantics remain consistent across supported LFE SDKs.

For the complete public API contract, see [`docs/SDK_API_CONTRACT_MATRIX.md`](docs/SDK_API_CONTRACT_MATRIX.md). For language-neutral interoperability semantics, use the separately distributed **LFE Interoperability Specification**.

## Build and test

The SDK uses a Rust native bridge linked through cgo.

```bash
make native
make test
```

`make native` builds the production native artifact with the canonical LFE BE production trust-root identity. `make test` builds an isolated test-native artifact with an ephemeral test trust root and runs the Go test suite.

## Licensing

The LFE Go SDK source code is licensed under the Apache License, Version 2.0.

The bundled LFE Engine native binary is proprietary PlanesLogic software. It is not licensed under the Apache License, Version 2.0 and is governed by separate LFE commercial license terms.

See [LICENSE](LICENSE), [NOTICE](NOTICE), and [native/COMMERCIAL-LICENSE-NOTICE](native/COMMERCIAL-LICENSE-NOTICE) for the distribution license boundary.

### Engine entitlement

Production Engine construction requires a valid signed LFE BE license.

Normal discovery:

```go
engine, err := lfe.New()
if err != nil {
    return err
}
defer engine.Close()
```

Explicit artifact path:

```go
engine, err := lfe.NewWithLicensePath("/path/to/license.json")
if err != nil {
    return err
}
defer engine.Close()
```

License verification is fail-closed. Missing, invalid, expired, incorrectly bound, or incompatible product/SKU artifacts do not open a licensed Engine.

## Source sequence and logical projections

`seq` is owned by the application or source system. LFE does not replace it with an internal application-facing identifier.

A single source coordinate may carry multiple logical projections:

```text
source seq 100
  |-- amount       UInt
  |-- risk_score   Int
  |-- event_time   DateTime
  `-- status       FlagSet
```

Define the projections once for an Engine:

```go
const (
    Amount    uint32 = 1
    RiskScore uint32 = 2
    EventTime uint32 = 3
    Status    uint32 = 4
)

if err := engine.DefineUInt(Amount, "amount", 64); err != nil {
    return err
}
if err := engine.DefineInt(RiskScore, "risk_score", 32); err != nil {
    return err
}
if err := engine.DefineDateTime(EventTime, "event_time"); err != nil {
    return err
}
if err := engine.DefineFlagSet(Status, "status", 4); err != nil {
    return err
}
```

The IDs supplied here are application-owned **Projection IDs** and should remain stable anywhere components are expected to address the same logical projection.

## Heterogeneous `AddBatch`

`AddBatch` accepts numeric, DateTime, and FlagSet logical values in one call. Multiple projections for the same source `seq` are valid.

```go
const (
    Active uint8 = iota
    Verified
    Fraud
    Premium
)

status := lfe.Flags(
    lfe.FlagSet(Active, 1),
    lfe.FlagSet(Verified, 1),
    lfe.FlagSet(Fraud, 1),
    lfe.FlagSet(Premium, 0),
)

rows := []lfe.AddRecord{
    {
        Seq:          100,
        ProjectionID: Amount,
        Value:        lfe.UIntValue(950_000),
    },
    {
        Seq:          100,
        ProjectionID: RiskScore,
        Value:        lfe.IntValue(95),
    },
    {
        Seq:          100,
        ProjectionID: EventTime,
        Value:        lfe.DateTimeUTC(time.Now()),
    },
    {
        Seq:          100,
        ProjectionID: Status,
        Value:        status,
    },
}

stats, err := engine.AddBatch(rows)
if err != nil {
    return err
}
_ = stats
```

Batch validation is performed against the declared projection type. Do not substitute an unsigned value for a DateTime or FlagSet projection.

## DateTime

Use `DateTimeUTC` when the application already has a `time.Time`:

```go
value := lfe.DateTimeUTC(time.Now())
```

The helper normalizes the supplied time to UTC and projects year, month, day, hour, minute, and second.

For a UTC string, use the fixed format `YYYY-MM-DD HH:MM:SS`:

```go
value, err := lfe.DateTimeUTCString("2026-08-30 12:23:45")
if err != nil {
    return err
}
```

DateTime queries select a component explicitly:

```go
q := lfe.NewQuery(EventTime, lfe.Eq, lfe.HourOperand(12))
```

## FlagSet

Applications own flag meaning and stable flag positions. Raw application values passed to `FlagSet` are `0` or `1`; the SDK owns shifting and packing.

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

Resolve uses the flag position, not the application's business name:

```go
q := lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true))
```

## Semantic pipeline with `SeqSetEx`

Use `ResolveEx` when the result will continue through LFE operations. `ResolveFromSetEx` applies the next predicate only to the existing candidate domain.

```go
amount, err := engine.ResolveEx(
    lfe.NewQuery(Amount, lfe.Gte, lfe.UInt(900_000)),
)
if err != nil {
    return err
}
defer amount.Close()

risk, err := engine.ResolveFromSetEx(
    amount,
    lfe.NewQuery(RiskScore, lfe.Gte, lfe.Int(90)),
)
if err != nil {
    return err
}
defer risk.Close()

fraud, err := engine.ResolveFromSetEx(
    risk,
    lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
)
if err != nil {
    return err
}
defer fraud.Close()
```

This keeps the candidate domain semantic in-process; no `[]seq` materialization or binary encoding is required between stages.

## Choosing a resolve surface

```text
Need []uint64 now?                  -> Resolve
Need bounded iteration/streaming?   -> ResolveSet
Need to keep computing on the set?  -> ResolveEx
```

`SeqSet` supports bounded materialization and chunked iteration. `SeqSetEx` is the semantic pipeline/interoperability result.

## `SeqSetEx` set algebra

Independent candidate sets can be composed without materializing source coordinates:

```go
combined, err := lfe.Merge(regionA, regionB)
if err != nil {
    return err
}
defer combined.Close()

common, err := lfe.Intersect(regionA, regionB)
if err != nil {
    return err
}
defer common.Close()

onlyA, err := lfe.Difference(regionA, regionB)
if err != nil {
    return err
}
defer onlyA.Close()
```

The operations are semantic union, intersection, and directional difference over source-sequence membership.

## LFESEQ01 boundary

Keep `SeqSetEx` while computation remains inside LFE. Encode only at a real downstream or adapter boundary:

```go
payload, err := fraud.Binary()
if err != nil {
    return err
}
```

`Binary()` returns the canonical `LFESEQ01` transport representation. The language-neutral format and adapter requirements belong to the **LFE Interoperability Specification**.

## Live mutation

An Engine may remain resident while the application mutates logical state and resolves again:

```go
if err := engine.UpdateFlagSet(seq, Status, updatedStatus); err != nil {
    return err
}

next, err := engine.ResolveEx(
    lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
)
```

A successful update is visible to subsequent resolves without rebuilding the Engine.

## Transaction-risk acceptance demo

The repository includes an executable application-shaped demo covering heterogeneous ingest, semantic narrowing, set algebra, LFESEQ01, exact truth parity, and live mutation:

```bash
LFE_BE_LICENSE_PATH=/path/to/license.json \
  go run ./demos/transaction-risk
```

Scale mode uses the same public SDK surface with deterministic generated transactions:

```bash
LFE_BE_LICENSE_PATH=/path/to/license.json \
  go run ./demos/transaction-risk \
  -records 1000000 \
  -batch 100000
```

The scale output is an executable acceptance/measurement aid. Treat results as measurements of the specific host, build, workload, and execution path rather than as universal performance guarantees.

## Persistence and correctness lifecycle

`Persist` and `Restore` are public advanced surfaces. Applications using persistence must follow the correctness lifecycle exposed by:

```text
CorrectnessGateState
ConfirmCorrectnessGate
ForceFullSync
CompleteFullSync
```

The authoritative source remains responsible for source payloads and for validating restored logical state before normal mutation resumes.

## Documentation

- [`docs/SDK_API_CONTRACT_MATRIX.md`](docs/SDK_API_CONTRACT_MATRIX.md) — complete Go public API usage contract.
- `docs/decisions/` — historical design decisions. These are not substitutes for the current public contract.
- **LFE Interoperability Specification** — source-sequence semantics, `SeqSetEx`, set algebra, `LFESEQ01`, and adapter contracts for cross-language/system integration.
