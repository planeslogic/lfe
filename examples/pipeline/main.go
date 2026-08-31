package main

import (
	"fmt"
	"log"

	lfe "github.com/planeslogic/lfe-go"
)

type transaction struct {
	Seq    uint64
	Amount uint64
	Risk   int64
	Fraud  bool
}

func main() {
	// ------------------------------------------------------------
	// Create an LFE Engine.
	// ------------------------------------------------------------

	engine := openDemoEngine()
	if engine == nil {
		return
	}
	defer engine.Close()

	// ------------------------------------------------------------
	// Source data.
	//
	// Each transaction remains owned by the application. LFE will
	// receive three Logical Projections associated with the same
	// source-owned Seq: amount, risk, and fraud status.
	// ------------------------------------------------------------

	source := []transaction{
		{Seq: 1, Amount: 1200, Risk: 95, Fraud: true},
		{Seq: 2, Amount: 1500, Risk: 70, Fraud: true},
		{Seq: 3, Amount: 1800, Risk: 98, Fraud: false},
		{Seq: 4, Amount: 2200, Risk: 99, Fraud: true},
		{Seq: 5, Amount: 500, Risk: 97, Fraud: true},
	}

	fmt.Println("source data")
	for _, row := range source {
		fmt.Printf(
			"seq=%d amount=%d risk=%d fraud=%t\n",
			row.Seq,
			row.Amount,
			row.Risk,
			row.Fraud,
		)
	}
	fmt.Println()

	// ------------------------------------------------------------
	// Define the Logical Projections.
	// ------------------------------------------------------------

	const (
		Amount uint32 = 1
		Risk   uint32 = 2
		Status uint32 = 3
	)

	const Fraud uint8 = 0

	if err := engine.DefineUInt(Amount, "amount", 32); err != nil {
		log.Fatal(err)
	}
	if err := engine.DefineInt(Risk, "risk", 16); err != nil {
		log.Fatal(err)
	}
	if err := engine.DefineFlagSet(Status, "status", 1); err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------
	// Ingest multiple projections for each source Seq.
	//
	// AddBatch accepts heterogeneous logical values while the
	// application keeps ownership of the transaction payload.
	// ------------------------------------------------------------

	rows := make([]lfe.AddRecord, 0, len(source)*3)

	for _, row := range source {
		fraud := uint8(0)
		if row.Fraud {
			fraud = 1
		}

		rows = append(rows,
			lfe.AddRecord{
				Seq:          row.Seq,
				ProjectionID: Amount,
				Value:        lfe.UIntValue(row.Amount),
			},
			lfe.AddRecord{
				Seq:          row.Seq,
				ProjectionID: Risk,
				Value:        lfe.IntValue(row.Risk),
			},
			lfe.AddRecord{
				Seq:          row.Seq,
				ProjectionID: Status,
				Value:        lfe.FlagSet(Fraud, fraud),
			},
		)
	}

	if _, err := engine.AddBatch(rows); err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------
	// Build a logical pipeline.
	//
	//  amount >= 1000
	//         ↓
	//  risk >= 90
	//         ↓
	//  fraud == true
	//
	// ResolveEx returns SeqSetEx. ResolveFromSetEx evaluates the
	// next condition only from the candidate set produced by the
	// previous stage.
	// ------------------------------------------------------------

	amount, err := engine.ResolveEx(
		lfe.NewQuery(Amount, lfe.Gte, lfe.UInt(1000)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer amount.Close()

	fmt.Println("amount >= 1000:", amount.Len())
	// Continue from the previous logical state.
	// Only Seq coordinates in "amount" are evaluated by this query.
	risk, err := engine.ResolveFromSetEx(
		amount,
		lfe.NewQuery(Risk, lfe.Gte, lfe.Int(90)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer risk.Close()

	fmt.Println("risk >= 90    :", risk.Len())
	// Continue from the "risk" state produced by the previous stage.
	fraud, err := engine.ResolveFromSetEx(
		risk,
		lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer fraud.Close()

	fmt.Println("fraud == true :", fraud.Len())

	// SeqSetEx contains source Seq coordinates. The application can
	// use them to hydrate the original source records when needed.

	fmt.Print("candidate seqs :")
	for _, row := range source {
		if fraud.Contains(row.Seq) {
			fmt.Print(" ", row.Seq)
		}
	}
	fmt.Println()
}

// openDemoEngine opens an LFE Engine using the standard license
// discovery mechanism.
//
// LFE demos require a Developer License. If a license has not
// been installed yet, the demo points the developer directly
// to the LFE Portal instead of exposing a low-level startup error.
func openDemoEngine() *lfe.Engine {
	engine, err := lfe.New()
	if err == nil {
		return engine
	}

	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println("  LFE Developer License Required")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("  This demo requires a free LFE Developer License.")
	fmt.Println()
	fmt.Println("  Get your Developer License:")
	fmt.Println("  https://lfe.planeslogic.com/portal.html")
	fmt.Println()
	fmt.Println("  Install the license and run this demo again.")
	fmt.Println()
	fmt.Println("============================================================")
	fmt.Println()

	return nil
}
