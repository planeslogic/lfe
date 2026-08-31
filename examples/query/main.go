package main

import (
	"fmt"
	"log"

	lfe "github.com/planeslogic/lfe-go"
)

type product struct {
	Seq   uint64
	Price uint64
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
	// The application owns the product records. Seq identifies
	// each source record; Price is the logical value we want LFE
	// to evaluate.
	// ------------------------------------------------------------

	source := []product{
		{Seq: 101, Price: 50},
		{Seq: 102, Price: 120},
		{Seq: 103, Price: 250},
		{Seq: 104, Price: 700},
	}

	fmt.Println("source data")
	for _, row := range source {
		fmt.Printf("seq=%d price=%d\n", row.Seq, row.Price)
	}
	fmt.Println()

	// ------------------------------------------------------------
	// Define the Price Logical Projection and ingest its values.
	// ------------------------------------------------------------

	const Price uint32 = 1

	if err := engine.DefineUInt(Price, "price", 32); err != nil {
		log.Fatal(err)
	}

	for _, row := range source {
		if err := engine.IngestUInt(row.Seq, Price, row.Price); err != nil {
			log.Fatal(err)
		}
	}

	// ------------------------------------------------------------
	// Query the same Logical Projection with different operators.
	//
	// Resolve returns source Seq coordinates. It does not return
	// or take ownership of the source product records.
	// ------------------------------------------------------------

	fmt.Println("queries")

	gte, err := engine.Resolve(
		lfe.NewQuery(Price, lfe.Gte, lfe.UInt(100)),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("price >= 100:", gte)

	eq, err := engine.Resolve(
		lfe.NewQuery(Price, lfe.Eq, lfe.UInt(250)),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("price == 250:", eq)

	lt, err := engine.Resolve(
		lfe.NewQuery(Price, lfe.Lt, lfe.UInt(500)),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("price < 500 :", lt)
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
