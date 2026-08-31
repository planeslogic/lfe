package main

import (
	"fmt"
	"log"

	lfe "github.com/planeslogic/lfe-go"
)

type record struct {
	Seq   uint64
	Score uint64
}

func main() {
	// ------------------------------------------------------------
	// Create an LFE Engine.
	//
	// The demo uses the standard LFE license discovery mechanism.
	// If no license is installed, a Developer License can be
	// obtained from the LFE Portal.
	// ------------------------------------------------------------

	engine := openDemoEngine()
	if engine == nil {
		return
	}
	defer engine.Close()

	// ------------------------------------------------------------
	// Source data
	//
	// These records belong to the application.
	// LFE does not own the source records or their payload.
	//
	// Seq is the source-owned coordinate used to identify each
	// record throughout the LFE pipeline.
	// ------------------------------------------------------------

	source := []record{
		{Seq: 1, Score: 40},
		{Seq: 2, Score: 75},
		{Seq: 3, Score: 90},
		{Seq: 4, Score: 60},
		{Seq: 5, Score: 100},
	}

	fmt.Println("source data")
	for _, row := range source {
		fmt.Printf("seq=%d score=%d\n", row.Seq, row.Score)
	}
	fmt.Println()

	// ------------------------------------------------------------
	// Define a Logical Projection.
	//
	// The application assigns Projection ID 1 to "score".
	// The source records remain owned by the application; only the
	// logical value needed for evaluation is projected into LFE.
	// ------------------------------------------------------------

	const Score uint32 = 1

	if err := engine.DefineUInt(Score, "score", 16); err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------
	// Ingest the score projection.
	//
	// Each logical value is associated with its source-owned Seq.
	//
	//   Seq 1 -> score 40
	//   Seq 2 -> score 75
	//   Seq 3 -> score 90
	//   Seq 4 -> score 60
	//   Seq 5 -> score 100
	//
	// LFE does not need to own the complete source record.
	// ------------------------------------------------------------

	for _, row := range source {
		if err := engine.IngestUInt(row.Seq, Score, row.Score); err != nil {
			log.Fatal(err)
		}
	}

	// ------------------------------------------------------------
	// Resolve logical conditions.
	//
	// Resolve returns the source Seq coordinates that satisfy
	// each condition. The application can then use those Seq
	// coordinates to access its source data.
	// ------------------------------------------------------------

	fmt.Println("queries")

	highScore, err := engine.Resolve(
		lfe.NewQuery(Score, lfe.Gte, lfe.UInt(70)),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("score >= 70 :", highScore)

	perfectScore, err := engine.Resolve(
		lfe.NewQuery(Score, lfe.Eq, lfe.UInt(100)),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("score == 100:", perfectScore)

	lowScore, err := engine.Resolve(
		lfe.NewQuery(Score, lfe.Lt, lfe.UInt(60)),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("score < 60  :", lowScore)
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
