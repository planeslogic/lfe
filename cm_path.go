package lfe

/*
#include <stdint.h>
#include <stdlib.h>
typedef struct CmHandle CmHandle;
typedef struct CmSeqSetHandle CmSeqSetHandle;
typedef struct CmSeqSetExHandle CmSeqSetExHandle;
typedef struct {
    uint64_t seq;
    uint32_t definition_id;
    uint8_t kind;
    _Bool negative;
    uint16_t year;
    uint8_t month;
    uint8_t day;
    uint8_t hour;
    uint8_t minute;
    uint8_t second;
    uint64_t value_hi;
    uint64_t value_lo;
} CmBatchAddFlat;
typedef struct {
    uint64_t records;
    uint64_t segments;
    uint64_t workers_used;
} CmBatchStatsFlat;
CmHandle* lfe_go_cm_new_licensed(const uint8_t*,size_t,const uint8_t*,size_t);
void lfe_go_cm_free(CmHandle*);
int32_t lfe_go_cm_define_numeric(CmHandle*,uint32_t,uint8_t,uint16_t,const uint8_t*,size_t);
int32_t lfe_go_cm_define_datetime(CmHandle*,uint32_t,const uint8_t*,size_t);
int32_t lfe_go_cm_define_flagset(CmHandle*,uint32_t,uint16_t,const uint8_t*,size_t);
int32_t lfe_go_cm_add_numeric(CmHandle*,uint64_t,uint32_t,_Bool,uint64_t,uint64_t);
int32_t lfe_go_cm_update_numeric(CmHandle*,uint64_t,uint32_t,_Bool,uint64_t,uint64_t);
int32_t lfe_go_cm_add_datetime(CmHandle*,uint64_t,uint32_t,uint16_t,uint8_t,uint8_t,uint8_t,uint8_t,uint8_t);
int32_t lfe_go_cm_update_datetime(CmHandle*,uint64_t,uint32_t,uint16_t,uint8_t,uint8_t,uint8_t,uint8_t,uint8_t);
int32_t lfe_go_cm_add_flagset(CmHandle*,uint64_t,uint32_t,uint64_t,uint64_t);
int32_t lfe_go_cm_update_flagset(CmHandle*,uint64_t,uint32_t,uint64_t,uint64_t);
int32_t lfe_go_cm_add_batch(CmHandle*,const CmBatchAddFlat*,size_t,CmBatchStatsFlat*);
int32_t lfe_go_cm_delete(CmHandle*,uint64_t);
int32_t lfe_go_cm_resolve_set(CmHandle*,uint32_t,uint8_t,uint8_t,uint8_t,_Bool,uint64_t,uint64_t,uint64_t,uint64_t,uint64_t,_Bool,_Bool,_Bool,CmSeqSetHandle**);
int32_t lfe_go_cm_resolve_ex(CmHandle*,uint32_t,uint8_t,uint8_t,uint8_t,_Bool,uint64_t,uint64_t,uint64_t,uint64_t,uint64_t,_Bool,_Bool,_Bool,CmSeqSetExHandle**);
int32_t lfe_go_cm_resolve_from_set(CmHandle*,CmSeqSetExHandle*,uint32_t,uint8_t,uint8_t,uint8_t,_Bool,uint64_t,uint64_t,uint64_t,uint64_t,uint64_t,_Bool,_Bool,_Bool,CmSeqSetHandle**);
int32_t lfe_go_cm_resolve_from_set_ex(CmHandle*,CmSeqSetExHandle*,uint32_t,uint8_t,uint8_t,uint8_t,_Bool,uint64_t,uint64_t,uint64_t,uint64_t,uint64_t,_Bool,_Bool,_Bool,CmSeqSetExHandle**);
void lfe_go_cm_seqset_free(CmSeqSetHandle*);
uint64_t lfe_go_cm_seqset_len(CmSeqSetHandle*);
_Bool lfe_go_cm_seqset_is_empty(CmSeqSetHandle*);
_Bool lfe_go_cm_seqset_contains(CmSeqSetHandle*,uint64_t);
int32_t lfe_go_cm_seqset_min(CmSeqSetHandle*,uint64_t*);
int32_t lfe_go_cm_seqset_max(CmSeqSetHandle*,uint64_t*);
intptr_t lfe_go_cm_seqset_materialize(CmSeqSetHandle*,uint64_t,uint64_t*,size_t);
intptr_t lfe_go_cm_seqset_fetch_next(CmSeqSetHandle*,uint64_t,_Bool,uint16_t,uint64_t*,size_t,_Bool*,uint64_t*,_Bool*);
void lfe_go_cm_seqset_ex_free(CmSeqSetExHandle*);
uint64_t lfe_go_cm_seqset_ex_len(CmSeqSetExHandle*);
_Bool lfe_go_cm_seqset_ex_is_empty(CmSeqSetExHandle*);
_Bool lfe_go_cm_seqset_ex_contains(CmSeqSetExHandle*,uint64_t);
int32_t lfe_go_cm_seqset_ex_min(CmSeqSetExHandle*,uint64_t*);
int32_t lfe_go_cm_seqset_ex_max(CmSeqSetExHandle*,uint64_t*);
intptr_t lfe_go_cm_seqset_ex_binary(CmSeqSetExHandle*,uint8_t*,size_t);
int32_t lfe_go_cm_seqset_ex_merge(CmSeqSetExHandle*,CmSeqSetExHandle*,CmSeqSetExHandle**);
int32_t lfe_go_cm_seqset_ex_intersect(CmSeqSetExHandle*,CmSeqSetExHandle*,CmSeqSetExHandle**);
int32_t lfe_go_cm_seqset_ex_difference(CmSeqSetExHandle*,CmSeqSetExHandle*,CmSeqSetExHandle**);
int32_t lfe_go_cm_persist(CmHandle*,const uint8_t*,size_t);
int32_t lfe_go_cm_restore(CmHandle*,const uint8_t*,size_t);
int32_t lfe_go_cm_correctness_gate_state(CmHandle*,uint8_t*);
int32_t lfe_go_cm_confirm_correctness_gate(CmHandle*,_Bool);
int32_t lfe_go_cm_force_full_sync(CmHandle*);
int32_t lfe_go_cm_complete_full_sync(CmHandle*);
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// Engine is the canonical public Go SDK entry point.
type Engine struct{ handle *C.CmHandle }

func New() (*Engine, error) {
	discovered, err := discoverLicense("")
	if err != nil {
		return nil, err
	}
	return newWithDiscoveredLicense(discovered)
}

func NewWithLicensePath(path string) (*Engine, error) {
	discovered, err := discoverLicense(path)
	if err != nil {
		return nil, err
	}
	return newWithDiscoveredLicense(discovered)
}

func newWithDiscoveredLicense(discovered discoveredLicense) (*Engine, error) {
	hostname, err := runtimeHostname()
	if err != nil {
		return nil, err
	}
	licenseBytes := []byte(discovered.json)
	hostnameBytes := []byte(hostname)
	var licensep, hostnamep *C.uint8_t
	if len(licenseBytes) != 0 {
		licensep = (*C.uint8_t)(unsafe.Pointer(&licenseBytes[0]))
	}
	if len(hostnameBytes) != 0 {
		hostnamep = (*C.uint8_t)(unsafe.Pointer(&hostnameBytes[0]))
	}
	h := C.lfe_go_cm_new_licensed(
		licensep, C.size_t(len(licenseBytes)),
		hostnamep, C.size_t(len(hostnameBytes)),
	)
	if h == nil {
		return nil, &Error{Op: "cm new licensed", Code: -1}
	}
	e := &Engine{handle: h}
	runtime.SetFinalizer(e, (*Engine).Close)
	return e, nil
}

func (e *Engine) Close() {
	if e != nil && e.handle != nil {
		C.lfe_go_cm_free(e.handle)
		e.handle = nil
		runtime.SetFinalizer(e, nil)
	}
}

func (e *Engine) Persist(root string) error {
	if e == nil || e.handle == nil {
		return &Error{Op: "cm persist", Code: -1}
	}
	b := []byte(root)
	var p *C.uint8_t
	if len(b) != 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	return nativeError("cm persist", C.lfe_go_cm_persist(e.handle, p, C.size_t(len(b))))
}

func (e *Engine) Restore(root string) error {
	if e == nil || e.handle == nil {
		return &Error{Op: "cm restore", Code: -1}
	}
	b := []byte(root)
	var p *C.uint8_t
	if len(b) != 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	return nativeError("cm restore", C.lfe_go_cm_restore(e.handle, p, C.size_t(len(b))))
}

type CorrectnessGateState uint8

const (
	CorrectnessReady CorrectnessGateState = iota
	CorrectnessPending
	CorrectnessDrifted
	CorrectnessFullSync
)

func (e *Engine) CorrectnessGateState() (CorrectnessGateState, error) {
	if e == nil || e.handle == nil {
		return CorrectnessReady, &Error{Op: "cm correctness gate state", Code: -1}
	}
	var state C.uint8_t
	if err := nativeError("cm correctness gate state", C.lfe_go_cm_correctness_gate_state(e.handle, &state)); err != nil {
		return CorrectnessReady, err
	}
	return CorrectnessGateState(state), nil
}

func (e *Engine) ConfirmCorrectnessGate(correct bool) error {
	if e == nil || e.handle == nil {
		return &Error{Op: "cm confirm correctness gate", Code: -1}
	}
	return nativeError("cm confirm correctness gate", C.lfe_go_cm_confirm_correctness_gate(e.handle, C._Bool(correct)))
}

func (e *Engine) ForceFullSync() error {
	if e == nil || e.handle == nil {
		return &Error{Op: "cm force full sync", Code: -1}
	}
	return nativeError("cm force full sync", C.lfe_go_cm_force_full_sync(e.handle))
}

func (e *Engine) CompleteFullSync() error {
	if e == nil || e.handle == nil {
		return &Error{Op: "cm complete full sync", Code: -1}
	}
	return nativeError("cm complete full sync", C.lfe_go_cm_complete_full_sync(e.handle))
}

func (e *Engine) defineNumeric(id uint32, name string, bits uint16, kind uint8) error {
	b := []byte(name)
	var p *C.uint8_t
	if len(b) != 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	return nativeError("cm define numeric", C.lfe_go_cm_define_numeric(e.handle, C.uint32_t(id), C.uint8_t(kind), C.uint16_t(bits), p, C.size_t(len(b))))
}

func (e *Engine) DefineUInt(id uint32, name string, bits uint16) error {
	return e.defineNumeric(id, name, bits, 0)
}
func (e *Engine) DefineInt(id uint32, name string, bits uint16) error {
	return e.defineNumeric(id, name, bits, 1)
}
func (e *Engine) DefineDateTime(id uint32, name string) error {
	b := []byte(name)
	var p *C.uint8_t
	if len(b) != 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	return nativeError("cm define datetime", C.lfe_go_cm_define_datetime(e.handle, C.uint32_t(id), p, C.size_t(len(b))))
}
func (e *Engine) DefineFlagSet(id uint32, name string, count uint16) error {
	b := []byte(name)
	var p *C.uint8_t
	if len(b) != 0 {
		p = (*C.uint8_t)(unsafe.Pointer(&b[0]))
	}
	return nativeError("cm define flagset", C.lfe_go_cm_define_flagset(e.handle, C.uint32_t(id), C.uint16_t(count), p, C.size_t(len(b))))
}

func (e *Engine) IngestUInt(seq uint64, id uint32, value uint64) error {
	return nativeError("cm add uint", C.lfe_go_cm_add_numeric(e.handle, C.uint64_t(seq), C.uint32_t(id), false, 0, C.uint64_t(value)))
}
func (e *Engine) IngestInt(seq uint64, id uint32, value int64) error {
	neg, mag := signedMagnitude(value)
	return nativeError("cm add int", C.lfe_go_cm_add_numeric(e.handle, C.uint64_t(seq), C.uint32_t(id), C._Bool(neg), 0, C.uint64_t(mag)))
}
func (e *Engine) IngestDateTime(seq uint64, id uint32, value DateTime) error {
	return nativeError("cm add datetime", C.lfe_go_cm_add_datetime(e.handle, C.uint64_t(seq), C.uint32_t(id), C.uint16_t(value.Year), C.uint8_t(value.Month), C.uint8_t(value.Day), C.uint8_t(value.Hour), C.uint8_t(value.Minute), C.uint8_t(value.Second)))
}
func (e *Engine) IngestFlagSet(seq uint64, id uint32, value FlagSetValue) error {
	if err := value.validate(); err != nil {
		return err
	}
	return nativeError("cm add flagset", C.lfe_go_cm_add_flagset(e.handle, C.uint64_t(seq), C.uint32_t(id), C.uint64_t(value.hi), C.uint64_t(value.lo)))
}
func (e *Engine) UpdateUInt(seq uint64, id uint32, value uint64) error {
	return nativeError("cm update uint", C.lfe_go_cm_update_numeric(e.handle, C.uint64_t(seq), C.uint32_t(id), false, 0, C.uint64_t(value)))
}
func (e *Engine) UpdateInt(seq uint64, id uint32, value int64) error {
	neg, mag := signedMagnitude(value)
	return nativeError("cm update int", C.lfe_go_cm_update_numeric(e.handle, C.uint64_t(seq), C.uint32_t(id), C._Bool(neg), 0, C.uint64_t(mag)))
}
func (e *Engine) UpdateDateTime(seq uint64, id uint32, value DateTime) error {
	return nativeError("cm update datetime", C.lfe_go_cm_update_datetime(e.handle, C.uint64_t(seq), C.uint32_t(id), C.uint16_t(value.Year), C.uint8_t(value.Month), C.uint8_t(value.Day), C.uint8_t(value.Hour), C.uint8_t(value.Minute), C.uint8_t(value.Second)))
}
func (e *Engine) UpdateFlagSet(seq uint64, id uint32, value FlagSetValue) error {
	if err := value.validate(); err != nil {
		return err
	}
	return nativeError("cm update flagset", C.lfe_go_cm_update_flagset(e.handle, C.uint64_t(seq), C.uint32_t(id), C.uint64_t(value.hi), C.uint64_t(value.lo)))
}

type batchValue struct {
	kind     uint8
	negative bool
	year     uint16
	month    uint8
	day      uint8
	hour     uint8
	minute   uint8
	second   uint8
	hi       uint64
	lo       uint64
	invalid  error
}

// Value is a logical value accepted by AddBatch. Concrete values are created
// by UIntValue, IntValue, DateTimeUTC/DateTimeUTCString, or FlagSet/Flags.
type Value interface {
	cmBatchValue() batchValue
}

type numericValue struct {
	negative  bool
	magnitude uint64
}

func (v numericValue) cmBatchValue() batchValue {
	return batchValue{kind: 0, negative: v.negative, lo: v.magnitude}
}

// UIntValue creates an unsigned numeric AddBatch value.
func UIntValue(value uint64) Value {
	return numericValue{magnitude: value}
}

// IntValue creates a signed numeric AddBatch value.
func IntValue(value int64) Value {
	negative, magnitude := signedMagnitude(value)
	return numericValue{negative: negative, magnitude: magnitude}
}

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

// AddBatch sends one heterogeneous batch into the Core batch-ingest path.
// Numeric, DateTime, and FlagSet values may be mixed in the same batch.
func (e *Engine) AddBatch(records []AddRecord) (BatchStats, error) {
	if e == nil || e.handle == nil {
		return BatchStats{}, &Error{Op: "cm add batch", Code: -1}
	}
	if len(records) == 0 {
		return BatchStats{}, nil
	}
	native := make([]C.CmBatchAddFlat, len(records))
	for i, record := range records {
		if record.Value == nil {
			return BatchStats{}, &Error{Op: "cm add batch", Code: -2}
		}
		value := record.Value.cmBatchValue()
		if value.invalid != nil {
			return BatchStats{}, value.invalid
		}
		native[i].seq = C.uint64_t(record.Seq)
		native[i].definition_id = C.uint32_t(record.ProjectionID)
		native[i].kind = C.uint8_t(value.kind)
		native[i].negative = C._Bool(value.negative)
		native[i].year = C.uint16_t(value.year)
		native[i].month = C.uint8_t(value.month)
		native[i].day = C.uint8_t(value.day)
		native[i].hour = C.uint8_t(value.hour)
		native[i].minute = C.uint8_t(value.minute)
		native[i].second = C.uint8_t(value.second)
		native[i].value_hi = C.uint64_t(value.hi)
		native[i].value_lo = C.uint64_t(value.lo)
	}
	var stats C.CmBatchStatsFlat
	code := C.lfe_go_cm_add_batch(
		e.handle,
		(*C.CmBatchAddFlat)(unsafe.Pointer(&native[0])),
		C.size_t(len(native)),
		&stats,
	)
	if err := nativeError("cm add batch", code); err != nil {
		return BatchStats{}, err
	}
	return BatchStats{
		Records:     uint64(stats.records),
		Segments:    uint64(stats.segments),
		WorkersUsed: uint64(stats.workers_used),
	}, nil
}

func (e *Engine) Delete(seq uint64) error {
	return nativeError("cm delete", C.lfe_go_cm_delete(e.handle, C.uint64_t(seq)))
}

type operandKind uint8

const (
	cmOperandUInt operandKind = iota
	cmOperandInt
	cmOperandDateTime
	cmOperandFlag
)

type Operand struct {
	kind     operandKind
	selector uint8
	value    uint64
	negative bool
	expected bool
}

func UInt(value uint64) Operand { return Operand{kind: cmOperandUInt, value: value} }
func Int(value int64) Operand {
	negative, magnitude := signedMagnitude(value)
	return Operand{kind: cmOperandInt, value: magnitude, negative: negative}
}
func YearOperand(value uint16) Operand {
	return Operand{kind: cmOperandDateTime, selector: 0, value: uint64(value)}
}
func MonthOperand(value uint8) Operand {
	return Operand{kind: cmOperandDateTime, selector: 1, value: uint64(value)}
}
func DayOperand(value uint8) Operand {
	return Operand{kind: cmOperandDateTime, selector: 2, value: uint64(value)}
}
func HourOperand(value uint8) Operand {
	return Operand{kind: cmOperandDateTime, selector: 3, value: uint64(value)}
}
func MinuteOperand(value uint8) Operand {
	return Operand{kind: cmOperandDateTime, selector: 4, value: uint64(value)}
}
func SecondOperand(value uint8) Operand {
	return Operand{kind: cmOperandDateTime, selector: 5, value: uint64(value)}
}
func Flag(flag uint8, expected bool) Operand {
	return Operand{kind: cmOperandFlag, selector: flag, expected: expected}
}

type Query struct {
	ProjectionID uint32
	Operator     Operator
	Operand      Operand
	SeqFromValue *uint64
	SeqToValue   *uint64
	LimitValue   *uint64
}

func NewQuery(id uint32, op Operator, operand Operand) Query {
	return Query{ProjectionID: id, Operator: op, Operand: operand}
}

func (q Query) SeqFrom(seq uint64) Query { q.SeqFromValue = &seq; return q }
func (q Query) SeqTo(seq uint64) Query   { q.SeqToValue = &seq; return q }
func (q Query) Limit(limit uint64) Query { q.LimitValue = &limit; return q }

func (e *Engine) ResolveSet(query Query) (*SeqSet, error) {
	var out *C.CmSeqSetHandle
	var seqFrom, seqTo, limit C.uint64_t
	var hasFrom, hasTo, hasLimit C._Bool
	if query.SeqFromValue != nil {
		seqFrom, hasFrom = C.uint64_t(*query.SeqFromValue), true
	}
	if query.SeqToValue != nil {
		seqTo, hasTo = C.uint64_t(*query.SeqToValue), true
	}
	if query.LimitValue != nil {
		limit, hasLimit = C.uint64_t(*query.LimitValue), true
	}
	code := C.lfe_go_cm_resolve_set(
		e.handle,
		C.uint32_t(query.ProjectionID),
		C.uint8_t(query.Operator),
		C.uint8_t(query.Operand.kind),
		C.uint8_t(query.Operand.selector),
		C._Bool(query.Operand.expected || query.Operand.negative),
		0,
		C.uint64_t(query.Operand.value),
		seqFrom, seqTo, limit,
		hasFrom, hasTo, hasLimit,
		&out,
	)
	if err := nativeError("cm resolve", code); err != nil {
		return nil, err
	}
	if out == nil {
		return nil, &Error{Op: "cm resolve", Code: -10}
	}
	set := &SeqSet{handle: out}
	runtime.SetFinalizer(set, (*SeqSet).Close)
	return set, nil
}

// ResolveEx returns the canonical extended semantic result for downstream
// adapters. Unlike Resolve, this path does not materialize source seq values.
func (e *Engine) ResolveEx(query Query) (*SeqSetEx, error) {
	if e == nil || e.handle == nil {
		return nil, &Error{Op: "cm resolve ex", Code: -1}
	}
	var out *C.CmSeqSetExHandle
	var seqFrom, seqTo, limit C.uint64_t
	var hasFrom, hasTo, hasLimit C._Bool
	if query.SeqFromValue != nil {
		seqFrom, hasFrom = C.uint64_t(*query.SeqFromValue), true
	}
	if query.SeqToValue != nil {
		seqTo, hasTo = C.uint64_t(*query.SeqToValue), true
	}
	if query.LimitValue != nil {
		limit, hasLimit = C.uint64_t(*query.LimitValue), true
	}
	code := C.lfe_go_cm_resolve_ex(
		e.handle,
		C.uint32_t(query.ProjectionID),
		C.uint8_t(query.Operator),
		C.uint8_t(query.Operand.kind),
		C.uint8_t(query.Operand.selector),
		C._Bool(query.Operand.expected || query.Operand.negative),
		0,
		C.uint64_t(query.Operand.value),
		seqFrom, seqTo, limit,
		hasFrom, hasTo, hasLimit,
		&out,
	)
	if err := nativeError("cm resolve ex", code); err != nil {
		return nil, err
	}
	if out == nil {
		return nil, &Error{Op: "cm resolve ex", Code: -10}
	}
	set := &SeqSetEx{handle: out}
	runtime.SetFinalizer(set, (*SeqSetEx).Close)
	return set, nil
}

// ResolveFromSet applies a query only within an existing semantic SeqSetEx
// candidate domain. The input may originate from another Engine when both
// engines use the same source-owned seq coordinate space.
//
// Core calls this primitive ResolveSet. The Go SDK uses ResolveFromSet because
// ResolveSet(Query) is an existing public method that returns a SeqSet.
func (e *Engine) ResolveFromSet(input *SeqSetEx, query Query) (*SeqSet, error) {
	if e == nil || e.handle == nil || input == nil || input.handle == nil {
		return nil, &Error{Op: "cm resolve from set", Code: -1}
	}
	var out *C.CmSeqSetHandle
	var seqFrom, seqTo, limit C.uint64_t
	var hasFrom, hasTo, hasLimit C._Bool
	if query.SeqFromValue != nil {
		seqFrom, hasFrom = C.uint64_t(*query.SeqFromValue), true
	}
	if query.SeqToValue != nil {
		seqTo, hasTo = C.uint64_t(*query.SeqToValue), true
	}
	if query.LimitValue != nil {
		limit, hasLimit = C.uint64_t(*query.LimitValue), true
	}
	code := C.lfe_go_cm_resolve_from_set(
		e.handle, input.handle,
		C.uint32_t(query.ProjectionID),
		C.uint8_t(query.Operator),
		C.uint8_t(query.Operand.kind),
		C.uint8_t(query.Operand.selector),
		C._Bool(query.Operand.expected || query.Operand.negative),
		0, C.uint64_t(query.Operand.value),
		seqFrom, seqTo, limit, hasFrom, hasTo, hasLimit, &out,
	)
	if err := nativeError("cm resolve from set", code); err != nil {
		return nil, err
	}
	if out == nil {
		return nil, &Error{Op: "cm resolve from set", Code: -10}
	}
	set := &SeqSet{handle: out}
	runtime.SetFinalizer(set, (*SeqSet).Close)
	return set, nil
}

// ResolveFromSetEx is the extended pipeline form of ResolveFromSet. It keeps
// semantic segmented membership intact so the result can feed another LFE
// engine without seq[] materialization or LFESEQ01 serialization.
func (e *Engine) ResolveFromSetEx(input *SeqSetEx, query Query) (*SeqSetEx, error) {
	if e == nil || e.handle == nil || input == nil || input.handle == nil {
		return nil, &Error{Op: "cm resolve from set ex", Code: -1}
	}
	var out *C.CmSeqSetExHandle
	var seqFrom, seqTo, limit C.uint64_t
	var hasFrom, hasTo, hasLimit C._Bool
	if query.SeqFromValue != nil {
		seqFrom, hasFrom = C.uint64_t(*query.SeqFromValue), true
	}
	if query.SeqToValue != nil {
		seqTo, hasTo = C.uint64_t(*query.SeqToValue), true
	}
	if query.LimitValue != nil {
		limit, hasLimit = C.uint64_t(*query.LimitValue), true
	}
	code := C.lfe_go_cm_resolve_from_set_ex(
		e.handle, input.handle,
		C.uint32_t(query.ProjectionID),
		C.uint8_t(query.Operator),
		C.uint8_t(query.Operand.kind),
		C.uint8_t(query.Operand.selector),
		C._Bool(query.Operand.expected || query.Operand.negative),
		0, C.uint64_t(query.Operand.value),
		seqFrom, seqTo, limit, hasFrom, hasTo, hasLimit, &out,
	)
	if err := nativeError("cm resolve from set ex", code); err != nil {
		return nil, err
	}
	if out == nil {
		return nil, &Error{Op: "cm resolve from set ex", Code: -10}
	}
	set := &SeqSetEx{handle: out}
	runtime.SetFinalizer(set, (*SeqSetEx).Close)
	return set, nil
}

// SeqSetEx is the native-backed semantic result returned by ResolveEx.
// Binary transport is produced only when Binary is requested by a downstream
// adapter; the ResolveEx path itself does not materialize seq[].
type SeqSetEx struct{ handle *C.CmSeqSetExHandle }

func (s *SeqSetEx) Close() {
	if s != nil && s.handle != nil {
		C.lfe_go_cm_seqset_ex_free(s.handle)
		s.handle = nil
		runtime.SetFinalizer(s, nil)
	}
}

func (s *SeqSetEx) Len() uint64 {
	if s == nil || s.handle == nil {
		return 0
	}
	return uint64(C.lfe_go_cm_seqset_ex_len(s.handle))
}

func (s *SeqSetEx) IsEmpty() bool {
	return s == nil || s.handle == nil || bool(C.lfe_go_cm_seqset_ex_is_empty(s.handle))
}

func (s *SeqSetEx) Contains(seq uint64) bool {
	return s != nil && s.handle != nil && bool(C.lfe_go_cm_seqset_ex_contains(s.handle, C.uint64_t(seq)))
}

func (s *SeqSetEx) Min() (uint64, bool) {
	if s == nil || s.handle == nil {
		return 0, false
	}
	var out C.uint64_t
	code := C.lfe_go_cm_seqset_ex_min(s.handle, &out)
	return uint64(out), code == 0
}

func (s *SeqSetEx) Max() (uint64, bool) {
	if s == nil || s.handle == nil {
		return 0, false
	}
	var out C.uint64_t
	code := C.lfe_go_cm_seqset_ex_max(s.handle, &out)
	return uint64(out), code == 0
}

// Merge returns the semantic union of two SeqSetEx candidate sets.
// The operation stays in sparse segment/block/word form and does not
// materialize seq[] or pass through LFESEQ01.
func Merge(a, b *SeqSetEx) (*SeqSetEx, error) {
	return seqSetExBinaryOp("cm seqset ex merge", a, b, seqSetExMerge)
}

// Intersect returns the semantic intersection of two SeqSetEx candidate sets.
func Intersect(a, b *SeqSetEx) (*SeqSetEx, error) {
	return seqSetExBinaryOp("cm seqset ex intersect", a, b, seqSetExIntersect)
}

// Difference returns the semantic difference a - b.
func Difference(a, b *SeqSetEx) (*SeqSetEx, error) {
	return seqSetExBinaryOp("cm seqset ex difference", a, b, seqSetExDifference)
}

type seqSetExOp uint8

const (
	seqSetExMerge seqSetExOp = iota
	seqSetExIntersect
	seqSetExDifference
)

func seqSetExBinaryOp(opName string, a, b *SeqSetEx, op seqSetExOp) (*SeqSetEx, error) {
	if a == nil || a.handle == nil || b == nil || b.handle == nil {
		return nil, &Error{Op: opName, Code: -1}
	}

	var out *C.CmSeqSetExHandle
	var code C.int32_t

	switch op {
	case seqSetExMerge:
		code = C.lfe_go_cm_seqset_ex_merge(a.handle, b.handle, &out)
	case seqSetExIntersect:
		code = C.lfe_go_cm_seqset_ex_intersect(a.handle, b.handle, &out)
	case seqSetExDifference:
		code = C.lfe_go_cm_seqset_ex_difference(a.handle, b.handle, &out)
	default:
		return nil, &Error{Op: opName, Code: -11}
	}

	if err := nativeError(opName, code); err != nil {
		return nil, err
	}
	if out == nil {
		return nil, &Error{Op: opName, Code: -10}
	}

	set := &SeqSetEx{handle: out}
	runtime.SetFinalizer(set, (*SeqSetEx).Close)
	return set, nil
}

// Binary returns the canonical LFESEQ01 adapter transport for this semantic
// SeqSetEx. The encoded bytes are cached by the native handle after first use.
func (s *SeqSetEx) Binary() ([]byte, error) {
	if s == nil || s.handle == nil {
		return nil, &Error{Op: "cm seqset ex binary", Code: -1}
	}
	n := C.lfe_go_cm_seqset_ex_binary(s.handle, nil, 0)
	if n < 0 {
		return nil, nativeResolveError("cm seqset ex binary", n)
	}
	if n == 0 {
		return []byte{}, nil
	}
	out := make([]byte, int(n))
	got := C.lfe_go_cm_seqset_ex_binary(
		s.handle,
		(*C.uint8_t)(unsafe.Pointer(&out[0])),
		C.size_t(len(out)),
	)
	if got < 0 {
		return nil, nativeResolveError("cm seqset ex binary", got)
	}
	return out[:int(got)], nil
}

func (e *Engine) Resolve(query Query) ([]uint64, error) {
	set, err := e.ResolveSet(query)
	if err != nil {
		return nil, err
	}
	defer set.Close()
	return set.Materialize(set.Len())
}

type SeqSet struct{ handle *C.CmSeqSetHandle }

func (s *SeqSet) Close() {
	if s != nil && s.handle != nil {
		C.lfe_go_cm_seqset_free(s.handle)
		s.handle = nil
		runtime.SetFinalizer(s, nil)
	}
}
func (s *SeqSet) Len() uint64 {
	if s == nil || s.handle == nil {
		return 0
	}
	return uint64(C.lfe_go_cm_seqset_len(s.handle))
}
func (s *SeqSet) IsEmpty() bool {
	return s == nil || s.handle == nil || bool(C.lfe_go_cm_seqset_is_empty(s.handle))
}
func (s *SeqSet) Contains(seq uint64) bool {
	return s != nil && s.handle != nil && bool(C.lfe_go_cm_seqset_contains(s.handle, C.uint64_t(seq)))
}
func (s *SeqSet) Min() (uint64, bool) {
	if s == nil || s.handle == nil {
		return 0, false
	}
	var out C.uint64_t
	code := C.lfe_go_cm_seqset_min(s.handle, &out)
	return uint64(out), code == 0
}
func (s *SeqSet) Max() (uint64, bool) {
	if s == nil || s.handle == nil {
		return 0, false
	}
	var out C.uint64_t
	code := C.lfe_go_cm_seqset_max(s.handle, &out)
	return uint64(out), code == 0
}
func (s *SeqSet) Materialize(max uint64) ([]uint64, error) {
	if s == nil || s.handle == nil {
		return nil, &Error{Op: "cm seqset materialize", Code: -1}
	}
	n := C.lfe_go_cm_seqset_materialize(s.handle, C.uint64_t(max), nil, 0)
	if n < 0 {
		return nil, nativeResolveError("cm seqset materialize", n)
	}
	if n == 0 {
		return []uint64{}, nil
	}
	out := make([]uint64, int(n))
	got := C.lfe_go_cm_seqset_materialize(s.handle, C.uint64_t(max), (*C.uint64_t)(unsafe.Pointer(&out[0])), C.size_t(len(out)))
	if got < 0 {
		return nil, nativeResolveError("cm seqset materialize", got)
	}
	return out[:int(got)], nil
}

// FetchNext returns the next ascending bounded chunk from this SeqSet.
// lastSeq is nil for the first chunk and the returned cursor must be supplied
// to the next call when hasNext is true. limit must be 1..256.
func (s *SeqSet) FetchNext(lastSeq *uint64, limit uint16) (seqs []uint64, nextLastSeq *uint64, hasNext bool, err error) {
	if s == nil || s.handle == nil {
		return nil, nil, false, &Error{Op: "cm seqset fetch next", Code: -1}
	}
	if limit == 0 || limit > 256 {
		return nil, nil, false, &Error{Op: "cm seqset fetch next", Code: -2}
	}

	out := make([]uint64, int(limit))
	var last C.uint64_t
	var hasLast C._Bool
	if lastSeq != nil {
		last = C.uint64_t(*lastSeq)
		hasLast = true
	}
	var nativeHasNext C._Bool
	var nativeLast C.uint64_t
	var nativeHasLast C._Bool

	n := C.lfe_go_cm_seqset_fetch_next(
		s.handle,
		last,
		hasLast,
		C.uint16_t(limit),
		(*C.uint64_t)(unsafe.Pointer(&out[0])),
		C.size_t(len(out)),
		&nativeHasNext,
		&nativeLast,
		&nativeHasLast,
	)
	if n < 0 {
		return nil, nil, false, nativeResolveError("cm seqset fetch next", n)
	}
	out = out[:int(n)]
	if bool(nativeHasLast) {
		v := uint64(nativeLast)
		nextLastSeq = &v
	}
	return out, nextLastSeq, bool(nativeHasNext), nil
}

// StreamChunks consumes the SeqSet without materializing the complete result.
// size must be 1..256, matching the canonical fixed SeqSet chunk capacity.
func (s *SeqSet) StreamChunks(size uint16, fn func([]uint64) error) error {
	if fn == nil {
		return &Error{Op: "cm seqset stream chunks", Code: -1}
	}
	var cursor *uint64
	for {
		chunk, next, hasNext, err := s.FetchNext(cursor, size)
		if err != nil {
			return err
		}
		if len(chunk) != 0 {
			if err := fn(chunk); err != nil {
				return err
			}
		}
		if !hasNext {
			return nil
		}
		if next == nil {
			return &Error{Op: "cm seqset stream chunks", Code: -2}
		}
		cursor = next
	}
}
