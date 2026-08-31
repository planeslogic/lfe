package main

import (
	"fmt"
	"log"
	"os"

	lfe "github.com/planeslogic/lfe-go"
)

func main() {
	// ------------------------------------------------------------
	// Standard license discovery.
	//
	// Most applications can simply call New(). The SDK searches
	// the standard LFE license locations and opens the Engine when
	// a valid license is available.
	// ------------------------------------------------------------

	engine := openDemoEngine()
	if engine == nil {
		return
	}
	engine.Close()

	fmt.Println("automatic license discovery: PASS")

	// ------------------------------------------------------------
	// Explicit license path.
	//
	// Some applications may manage the license location directly.
	// Set LFE_BE_LICENSE_PATH to try the same Engine startup using
	// NewWithLicensePath().
	//
	// Example:
	//
	//   LFE_BE_LICENSE_PATH=/path/to/license.json \
	//     go run ./license-discovery
	// ------------------------------------------------------------

	licensePath := os.Getenv("LFE_BE_LICENSE_PATH")
	if licensePath == "" {
		fmt.Println()
		fmt.Println("explicit license path: SKIP")
		fmt.Println("set LFE_BE_LICENSE_PATH to try NewWithLicensePath()")
		return
	}

	explicitEngine, err := lfe.NewWithLicensePath(licensePath)
	if err != nil {
		log.Fatal(err)
	}
	defer explicitEngine.Close()

	fmt.Println("explicit license path: PASS")
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
