# LFE Go SDK Demos

These demos are small, self-contained examples designed to show one LFE concept at a time.

The examples intentionally keep the code explicit. LFE calls are not hidden behind application helpers or framework abstractions, so you can see how source-owned data is projected into LFE and how LFE returns source sequence coordinates.

## Start here

From the repository root:

```bash
go run ./examples/get-started
```

Or, from inside `examples/`:

```bash
go run ./get-started
```

## Developer License

The LFE Engine requires a signed license to run.

A **free LFE Developer License** is available for evaluating, learning, and building with the LFE SDK.

**No purchase is required to run these demos.**

Get your Developer License:

https://lfe.planeslogic.com/portal.html

After installing the license, run the demo again. `lfe.New()` uses the standard LFE license discovery mechanism.

## Demo map

| Demo | What it shows |
| --- | --- |
| `get-started` | Source data → Logical Projection → Resolve → source Seq coordinates |
| `license-discovery` | Standard license discovery and an explicit license path |
| `persistence` | Persist an Engine state, restore it into a new Engine, and Resolve again |
| `query` | Evaluate the same Logical Projection with different comparison operators |
| `pipeline` | Continue logical evaluation from one `SeqSetEx` into the next query |
| `set-algebra` | Merge, intersect, and subtract `SeqSetEx` results |
| `transaction-risk` | Larger end-to-end transaction-risk and scale example |

## Mental model

LFE does not replace ownership of your source records.

A typical flow is:

```text
Application source data
        ↓
source-owned Seq + logical value
        ↓
Logical Projection
        ↓
Query / Resolve
        ↓
SeqSet / SeqSetEx
        ↓
source Seq coordinates
        ↓
application / downstream source
```

Start with `get-started`, then read the other demos independently. Each demo has its own `main.go` and is intended to be understandable without reading the others.
