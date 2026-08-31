package main

import (
	"fmt"
	"log"

	lfe "github.com/planeslogic/lfe-go"
)

type record struct {
	Seq      uint64
	Score    uint64
	Priority uint64
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
	// We will create two independent logical result sets:
	//
	//   A = score >= 80
	//   B = priority == 1
	//
	// Then compose those SeqSetEx results without materializing
	// the application records into LFE.
	// ------------------------------------------------------------

	source := []record{
		{Seq: 1, Score: 90, Priority: 1},
		{Seq: 2, Score: 95, Priority: 0},
		{Seq: 3, Score: 60, Priority: 1},
		{Seq: 4, Score: 40, Priority: 0},
	}

	fmt.Println("source data")
	for _, row := range source {
		fmt.Printf(
			"seq=%d score=%d priority=%d\n",
			row.Seq,
			row.Score,
			row.Priority,
		)
	}
	fmt.Println()

	// ------------------------------------------------------------
	// Define and ingest the Logical Projections.
	// ------------------------------------------------------------

	const (
		Score    uint32 = 1
		Priority uint32 = 2
	)

	if err := engine.DefineUInt(Score, "score", 16); err != nil {
		log.Fatal(err)
	}
	if err := engine.DefineUInt(Priority, "priority", 8); err != nil {
		log.Fatal(err)
	}

	rows := make([]lfe.AddRecord, 0, len(source)*2)

	for _, row := range source {
		rows = append(rows,
			lfe.AddRecord{
				Seq:          row.Seq,
				ProjectionID: Score,
				Value:        lfe.UIntValue(row.Score),
			},
			lfe.AddRecord{
				Seq:          row.Seq,
				ProjectionID: Priority,
				Value:        lfe.UIntValue(row.Priority),
			},
		)
	}

	if _, err := engine.AddBatch(rows); err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------
	// Resolve two independent SeqSetEx results.
	//
	// A = highScore = {1, 2}
	// B = priority  = {1, 3}
	// ------------------------------------------------------------

	highScore, err := engine.ResolveEx(
		lfe.NewQuery(Score, lfe.Gte, lfe.UInt(80)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer highScore.Close()

	priority, err := engine.ResolveEx(
		lfe.NewQuery(Priority, lfe.Eq, lfe.UInt(1)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer priority.Close()

	// ------------------------------------------------------------
	// Compose the result sets.
	//
	// Merge      A ∪ B  -> {1, 2, 3}
	// Intersect  A ∩ B  -> {1}
	// Difference A - B  -> {2}
	// ------------------------------------------------------------

	merged, err := lfe.Merge(highScore, priority)
	if err != nil {
		log.Fatal(err)
	}
	defer merged.Close()

	common, err := lfe.Intersect(highScore, priority)
	if err != nil {
		log.Fatal(err)
	}
	defer common.Close()

	remaining, err := lfe.Difference(highScore, priority)
	if err != nil {
		log.Fatal(err)
	}
	defer remaining.Close()

	fmt.Println("A · score >= 80")
	fmt.Println("  candidates :", highScore.Len())

	fmt.Println("B · priority == 1")
	fmt.Println("  candidates :", priority.Len())

	fmt.Println("A ∪ B · merge")
	fmt.Println("  candidates :", merged.Len())

	fmt.Println("A ∩ B · intersect")
	fmt.Println("  candidates :", common.Len())

	fmt.Println("A - B · difference")
	fmt.Println("  candidates :", remaining.Len())
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
