package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	lfe "github.com/planeslogic/lfe"
)

const (
	Amount    uint32 = 1
	RiskScore uint32 = 2
	EventTime uint32 = 3
	Status    uint32 = 4
)

const (
	Active uint8 = iota
	Verified
	Fraud
	Premium
)

type transaction struct {
	Seq       uint64
	Amount    uint64
	RiskScore int64
	EventTime time.Time
	Active    bool
	Verified  bool
	Fraud     bool
	Premium   bool
}

func main() {
	recordCount := flag.Int(
		"records",
		5,
		"number of source transactions",
	)
	batchSize := flag.Int(
		"batch",
		100_000,
		"source transactions per AddBatch",
	)

	flag.Parse()

	if *recordCount <= 0 {
		log.Fatal("-records must be greater than zero")
	}

	if *batchSize <= 0 {
		log.Fatal("-batch must be greater than zero")
	}

	if *recordCount == 5 {
		runAcceptance()
		return
	}

	runScale(*recordCount, *batchSize)
}
func runScale(recordCount, batchSize int) {
	licensePath := os.Getenv("LFE_BE_LICENSE_PATH")
	if licensePath == "" {
		log.Fatal("LFE_BE_LICENSE_PATH is required")
	}

	engine, err := lfe.NewWithLicensePath(licensePath)
	must(err)
	defer engine.Close()

	must(engine.DefineUInt(Amount, "amount", 64))
	must(engine.DefineInt(RiskScore, "risk_score", 32))
	must(engine.DefineDateTime(EventTime, "event_time"))
	must(engine.DefineFlagSet(Status, "status", 4))

	fmt.Println("LFE Go SDK · Transaction Risk Pipeline")
	fmt.Println()
	fmt.Printf("source transactions : %d\n", recordCount)
	fmt.Printf("logical projections : 4\n")
	fmt.Printf("batch size          : %d\n", batchSize)
	fmt.Println()

	ingestStarted := time.Now()

	var (
		totalRecords  uint64
		totalSegments uint64
		maxWorkers    uint64
	)

	for start := 0; start < recordCount; start += batchSize {
		end := start + batchSize
		if end > recordCount {
			end = recordCount
		}

		rows := make([]lfe.AddRecord, 0, (end-start)*4)

		for i := start; i < end; i++ {
			tx := generatedTransaction(uint64(i + 1))
			rows = append(rows, records(tx)...)
		}

		stats, err := engine.AddBatch(rows)
		must(err)

		totalRecords += stats.Records
		totalSegments += stats.Segments

		if stats.WorkersUsed > maxWorkers {
			maxWorkers = stats.WorkersUsed
		}
	}

	ingestElapsed := time.Since(ingestStarted)

	wantRecords := uint64(recordCount) * 4
	if totalRecords != wantRecords {
		log.Fatalf(
			"ingest record mismatch: got=%d want=%d",
			totalRecords,
			wantRecords,
		)
	}

	fmt.Println("ingest")
	fmt.Printf("records             : %d\n", totalRecords)
	fmt.Printf("batch segment visits: %d\n", totalSegments)
	fmt.Printf("max workers         : %d\n", maxWorkers)
	fmt.Printf("elapsed             : %s\n", ingestElapsed)
	fmt.Printf(
		"values/sec          : %.0f\n",
		float64(totalRecords)/ingestElapsed.Seconds(),
	)
	fmt.Println()

	amountStarted := time.Now()
	amount, err := engine.ResolveEx(
		lfe.NewQuery(Amount, lfe.Gte, lfe.UInt(900_000)),
	)
	must(err)
	amountElapsed := time.Since(amountStarted)
	defer amount.Close()

	riskStarted := time.Now()
	risk, err := engine.ResolveFromSetEx(
		amount,
		lfe.NewQuery(RiskScore, lfe.Gte, lfe.Int(90)),
	)
	must(err)
	riskElapsed := time.Since(riskStarted)
	defer risk.Close()

	fraudStarted := time.Now()
	fraud, err := engine.ResolveFromSetEx(
		risk,
		lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
	)
	must(err)
	fraudElapsed := time.Since(fraudStarted)
	defer fraud.Close()

	pipelineElapsed :=
		amountElapsed +
			riskElapsed +
			fraudElapsed

	fmt.Println("pipeline")
	fmt.Printf(
		"amount >= 900000    : %d candidates, %s\n",
		amount.Len(),
		amountElapsed,
	)
	fmt.Printf(
		"risk >= 90          : %d candidates, %s\n",
		risk.Len(),
		riskElapsed,
	)
	fmt.Printf(
		"fraud == true       : %d candidates, %s\n",
		fraud.Len(),
		fraudElapsed,
	)
	fmt.Printf("pipeline total      : %s\n", pipelineElapsed)
	fmt.Println()

	truthStarted := time.Now()
	expected := generatedTruth(recordCount)
	truthElapsed := time.Since(truthStarted)

	assertSet("scale-pipeline", fraud, expected)

	fmt.Println("truth oracle")
	fmt.Printf("expected            : %d\n", len(expected))
	fmt.Printf("actual              : %d\n", fraud.Len())
	fmt.Printf("oracle elapsed      : %s\n", truthElapsed)
	fmt.Println("membership parity   : PASS")
	fmt.Println()

	highRiskStarted := time.Now()
	highRisk, err := engine.ResolveEx(
		lfe.NewQuery(RiskScore, lfe.Gte, lfe.Int(90)),
	)
	must(err)
	highRiskElapsed := time.Since(highRiskStarted)
	defer highRisk.Close()

	fraudAllStarted := time.Now()
	fraudAll, err := engine.ResolveEx(
		lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
	)
	must(err)
	fraudAllElapsed := time.Since(fraudAllStarted)
	defer fraudAll.Close()

	mergeStarted := time.Now()
	merged, err := lfe.Merge(highRisk, fraudAll)
	must(err)
	mergeElapsed := time.Since(mergeStarted)
	defer merged.Close()

	intersectStarted := time.Now()
	common, err := lfe.Intersect(highRisk, fraudAll)
	must(err)
	intersectElapsed := time.Since(intersectStarted)
	defer common.Close()

	differenceStarted := time.Now()
	difference, err := lfe.Difference(highRisk, fraudAll)
	must(err)
	differenceElapsed := time.Since(differenceStarted)
	defer difference.Close()

	fmt.Println("set algebra")
	fmt.Printf(
		"high-risk           : %d, resolve=%s\n",
		highRisk.Len(),
		highRiskElapsed,
	)
	fmt.Printf(
		"fraud               : %d, resolve=%s\n",
		fraudAll.Len(),
		fraudAllElapsed,
	)
	fmt.Printf(
		"merge               : %d, %s\n",
		merged.Len(),
		mergeElapsed,
	)
	fmt.Printf(
		"intersect           : %d, %s\n",
		common.Len(),
		intersectElapsed,
	)
	fmt.Printf(
		"difference          : %d, %s\n",
		difference.Len(),
		differenceElapsed,
	)
	fmt.Println()

	binaryStarted := time.Now()
	binary, err := fraud.Binary()
	must(err)
	binaryElapsed := time.Since(binaryStarted)

	if len(binary) < 8 {
		log.Fatalf("LFESEQ01 payload too short: %d", len(binary))
	}

	if !bytes.Equal(binary[:8], []byte("LFESEQ01")) {
		log.Fatalf(
			"unexpected binary magic: %q",
			binary[:8],
		)
	}

	fmt.Println("LFESEQ01")
	fmt.Printf("bytes               : %d\n", len(binary))
	fmt.Printf("encode              : %s\n", binaryElapsed)
	fmt.Printf("magic               : %s PASS\n", binary[:8])
	fmt.Println()

	fmt.Println("scale transaction-risk demo: PASS")
}
func runAcceptance() {
	licensePath := os.Getenv("LFE_BE_LICENSE_PATH")
	if licensePath == "" {
		log.Fatal("LFE_BE_LICENSE_PATH is required")
	}

	engine, err := lfe.NewWithLicensePath(licensePath)
	must(err)
	defer engine.Close()

	must(engine.DefineUInt(Amount, "amount", 64))
	must(engine.DefineInt(RiskScore, "risk_score", 32))
	must(engine.DefineDateTime(EventTime, "event_time"))
	must(engine.DefineFlagSet(Status, "status", 4))

	transactions := []transaction{
		{
			Seq:       7,
			Amount:    1_250_000,
			RiskScore: 96,
			EventTime: utc(2026, 8, 30, 10, 15, 0),
			Active:    true,
			Verified:  true,
			Fraud:     true,
		},
		{
			Seq:       999_999,
			Amount:    920_000,
			RiskScore: 94,
			EventTime: utc(2026, 8, 30, 11, 20, 0),
			Active:    true,
			Verified:  true,
			Fraud:     false,
		},
		{
			Seq:       1_000_000,
			Amount:    980_000,
			RiskScore: 91,
			EventTime: utc(2026, 8, 30, 12, 25, 0),
			Active:    true,
			Verified:  false,
			Fraud:     true,
			Premium:   true,
		},
		{
			Seq:       1_000_001,
			Amount:    450_000,
			RiskScore: 99,
			EventTime: utc(2026, 8, 30, 13, 30, 0),
			Active:    true,
			Verified:  true,
			Fraud:     true,
		},
		{
			Seq:       2_000_003,
			Amount:    2_100_000,
			RiskScore: 72,
			EventTime: utc(2026, 8, 30, 14, 35, 0),
			Active:    true,
			Verified:  true,
			Fraud:     true,
			Premium:   true,
		},
	}

	rows := make([]lfe.AddRecord, 0, len(transactions)*4)
	for _, tx := range transactions {
		rows = append(rows, records(tx)...)
	}

	stats, err := engine.AddBatch(rows)
	must(err)

	fmt.Printf(
		"ingest: records=%d segments=%d workers=%d\n",
		stats.Records,
		stats.Segments,
		stats.WorkersUsed,
	)

	if stats.Records != uint64(len(transactions)*4) {
		log.Fatalf("unexpected record count: got=%d want=%d",
			stats.Records,
			len(transactions)*4,
		)
	}

	if stats.Segments != 3 {
		log.Fatalf("unexpected segment count: got=%d want=3", stats.Segments)
	}

	// ------------------------------------------------------------
	// Pipeline
	//
	// amount >= 900,000
	//       ↓
	// risk >= 90
	//       ↓
	// Fraud = true
	// ------------------------------------------------------------

	amount, err := engine.ResolveEx(
		lfe.NewQuery(Amount, lfe.Gte, lfe.UInt(900_000)),
	)
	must(err)
	defer amount.Close()

	risk, err := engine.ResolveFromSetEx(
		amount,
		lfe.NewQuery(RiskScore, lfe.Gte, lfe.Int(90)),
	)
	must(err)
	defer risk.Close()

	fraud, err := engine.ResolveFromSetEx(
		risk,
		lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
	)
	must(err)
	defer fraud.Close()

	fmt.Printf(
		"pipeline: amount=%d risk=%d fraud=%d\n",
		amount.Len(),
		risk.Len(),
		fraud.Len(),
	)

	expected := truth(transactions)
	assertSet("pipeline", fraud, expected)

	fmt.Printf("truth:    %v PASS\n", expected)

	// ------------------------------------------------------------
	// SeqSetEx algebra
	// ------------------------------------------------------------

	highRisk, err := engine.ResolveEx(
		lfe.NewQuery(RiskScore, lfe.Gte, lfe.Int(90)),
	)
	must(err)
	defer highRisk.Close()

	fraudAll, err := engine.ResolveEx(
		lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
	)
	must(err)
	defer fraudAll.Close()

	merged, err := lfe.Merge(highRisk, fraudAll)
	must(err)
	defer merged.Close()

	common, err := lfe.Intersect(highRisk, fraudAll)
	must(err)
	defer common.Close()

	difference, err := lfe.Difference(highRisk, fraudAll)
	must(err)
	defer difference.Close()

	fmt.Printf(
		"algebra:  high-risk=%d fraud=%d merge=%d intersect=%d difference=%d\n",
		highRisk.Len(),
		fraudAll.Len(),
		merged.Len(),
		common.Len(),
		difference.Len(),
	)

	assertSet(
		"intersect",
		common,
		[]uint64{7, 1_000_000, 1_000_001},
	)

	assertSet(
		"difference",
		difference,
		[]uint64{999_999},
	)

	// ------------------------------------------------------------
	// LFESEQ01
	// ------------------------------------------------------------

	binary, err := fraud.Binary()
	must(err)

	if len(binary) < 8 {
		log.Fatalf("LFESEQ01 payload too short: %d", len(binary))
	}

	if !bytes.Equal(binary[:8], []byte("LFESEQ01")) {
		log.Fatalf("unexpected binary magic: %q", binary[:8])
	}

	fmt.Printf(
		"LFESEQ01: bytes=%d magic=%s PASS\n",
		len(binary),
		string(binary[:8]),
	)

	// ------------------------------------------------------------
	// Live mutation on the same resident Engine.
	//
	// seq 999,999 was high amount + high risk but Fraud=false.
	// Flip Fraud=true; it must enter the final candidate set without
	// rebuilding the Engine.
	// ------------------------------------------------------------

	updated := transactions[1]
	updated.Fraud = true

	must(engine.UpdateFlagSet(
		updated.Seq,
		Status,
		flags(updated),
	))

	transactions[1] = updated

	amountAfter, err := engine.ResolveEx(
		lfe.NewQuery(Amount, lfe.Gte, lfe.UInt(900_000)),
	)
	must(err)
	defer amountAfter.Close()

	riskAfter, err := engine.ResolveFromSetEx(
		amountAfter,
		lfe.NewQuery(RiskScore, lfe.Gte, lfe.Int(90)),
	)
	must(err)
	defer riskAfter.Close()

	fraudAfter, err := engine.ResolveFromSetEx(
		riskAfter,
		lfe.NewQuery(Status, lfe.Eq, lfe.Flag(Fraud, true)),
	)
	must(err)
	defer fraudAfter.Close()

	expectedAfter := truth(transactions)
	assertSet("pipeline-after-update", fraudAfter, expectedAfter)

	fmt.Printf(
		"update:   seq=%d fraud=false->true\n",
		updated.Seq,
	)
	fmt.Printf(
		"rerun:    amount=%d risk=%d fraud=%d candidates=%v PASS\n",
		amountAfter.Len(),
		riskAfter.Len(),
		fraudAfter.Len(),
		expectedAfter,
	)

	fmt.Println("transaction-risk demo: PASS")
}

func records(tx transaction) []lfe.AddRecord {
	return []lfe.AddRecord{
		{
			Seq:          tx.Seq,
			ProjectionID: Amount,
			Value:        lfe.UIntValue(tx.Amount),
		},
		{
			Seq:          tx.Seq,
			ProjectionID: RiskScore,
			Value:        lfe.IntValue(tx.RiskScore),
		},
		{
			Seq:          tx.Seq,
			ProjectionID: EventTime,
			Value:        lfe.DateTimeUTC(tx.EventTime),
		},
		{
			Seq:          tx.Seq,
			ProjectionID: Status,
			Value:        flags(tx),
		},
	}
}

func flags(tx transaction) lfe.FlagSetValue {
	return lfe.Flags(
		lfe.FlagSet(Active, raw(tx.Active)),
		lfe.FlagSet(Verified, raw(tx.Verified)),
		lfe.FlagSet(Fraud, raw(tx.Fraud)),
		lfe.FlagSet(Premium, raw(tx.Premium)),
	)
}

func truth(transactions []transaction) []uint64 {
	var out []uint64

	for _, tx := range transactions {
		if tx.Amount >= 900_000 &&
			tx.RiskScore >= 90 &&
			tx.Fraud {
			out = append(out, tx.Seq)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})

	return out
}

type seqSet interface {
	Len() uint64
	Contains(uint64) bool
}

func assertSet(name string, set seqSet, expected []uint64) {
	if set.Len() != uint64(len(expected)) {
		log.Fatalf(
			"%s len mismatch: got=%d want=%d",
			name,
			set.Len(),
			len(expected),
		)
	}

	for _, seq := range expected {
		if !set.Contains(seq) {
			log.Fatalf("%s missing seq=%d", name, seq)
		}
	}
}

func raw(value bool) uint8 {
	if value {
		return 1
	}
	return 0
}

func utc(
	year int,
	month time.Month,
	day int,
	hour int,
	minute int,
	second int,
) time.Time {
	return time.Date(
		year,
		month,
		day,
		hour,
		minute,
		second,
		0,
		time.UTC,
	)
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func generatedTransaction(seq uint64) transaction {
	return transaction{
		Seq:       seq,
		Amount:    300_000 + ((seq * 7919) % 2_000_000),
		RiskScore: int64((seq*37)%141) - 20,
		EventTime: time.Date(
			2026,
			8,
			30,
			int(seq%24),
			int((seq*7)%60),
			int((seq*13)%60),
			0,
			time.UTC,
		),
		Active:   seq%2 == 0,
		Verified: seq%3 != 0,
		Fraud:    seq%11 == 0,
		Premium:  seq%7 == 0,
	}
}

func generatedTruth(recordCount int) []uint64 {
	out := make([]uint64, 0)

	for i := 0; i < recordCount; i++ {
		tx := generatedTransaction(uint64(i + 1))

		if tx.Amount >= 900_000 &&
			tx.RiskScore >= 90 &&
			tx.Fraud {
			out = append(out, tx.Seq)
		}
	}

	return out
}
