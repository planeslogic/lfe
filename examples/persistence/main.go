package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	lfe "github.com/planeslogic/lfe-go"
)

type record struct {
	Seq   uint64
	Score uint64
}

func main() {
	// ------------------------------------------------------------
	// Create the first LFE Engine.
	//
	// This Engine will receive logical state and persist it to a
	// local store.
	// ------------------------------------------------------------

	writer := openDemoEngine()
	if writer == nil {
		return
	}

	// ------------------------------------------------------------
	// Source data.
	//
	// The application owns these records. LFE receives only the
	// source-owned Seq and the logical "score" value.
	// ------------------------------------------------------------

	source := []record{
		{Seq: 10, Score: 55},
		{Seq: 20, Score: 80},
		{Seq: 30, Score: 95},
	}

	fmt.Println("source data")
	for _, row := range source {
		fmt.Printf("seq=%d score=%d\n", row.Seq, row.Score)
	}
	fmt.Println()

	// ------------------------------------------------------------
	// Define and ingest the Logical Projection.
	// ------------------------------------------------------------

	const Score uint32 = 1

	if err := writer.DefineUInt(Score, "score", 16); err != nil {
		log.Fatal(err)
	}

	for _, row := range source {
		if err := writer.IngestUInt(row.Seq, Score, row.Score); err != nil {
			log.Fatal(err)
		}
	}

	// ------------------------------------------------------------
	// Persist the Engine state.
	//
	// Resolve the user's home directory explicitly. "~" is a shell
	// convention and should not be passed as a literal store path.
	// ------------------------------------------------------------

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	store := filepath.Join(home, ".lfe-store", "persistence-demo")

	if err := writer.Persist(store); err != nil {
		log.Fatal(err)
	}
	writer.Close()

	fmt.Println("persisted:", store)

	// ------------------------------------------------------------
	// Create a new Engine and restore the persisted state.
	//
	// The Logical Projection is defined again so this Engine has
	// the same logical schema before Restore().
	// ------------------------------------------------------------

	reader := openDemoEngine()
	if reader == nil {
		return
	}
	defer reader.Close()

	if err := reader.DefineUInt(Score, "score", 16); err != nil {
		log.Fatal(err)
	}

	if err := reader.Restore(store); err != nil {
		log.Fatal(err)
	}

	// ------------------------------------------------------------
	// Resolve against the restored Engine.
	//
	// The result contains source Seq coordinates from the state
	// restored from disk.
	// ------------------------------------------------------------

	seqs, err := reader.Resolve(
		lfe.NewQuery(Score, lfe.Gte, lfe.UInt(80)),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("restored score >= 80:", seqs)
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
