package rs485

import (
	. "github.com/volkszaehler/mbmd/meters"
)

func init() {
	Register("CGEM24_E1", NewCarloGavazziEM24_E1Producer)
}

type CarloGavazziEM24_E1Producer struct {
	Opcodes
}

func NewCarloGavazziEM24_E1Producer() Producer {
	/***
	 * Note: Carlo Gavazzi EM24_E1 (Ethernet)
	 * Doc for EM24_E1: https://www.enika.eu/data/files/produkty/energy%20m/CP/em24%20ethernet%20cp.pdf
	 */
	ops := Opcodes{
		VoltageL1: 0x00,
		VoltageL2: 0x02,
		VoltageL3: 0x04,
		CurrentL1: 0x0C,
		CurrentL2: 0x0E,
		CurrentL3: 0x10,
		PowerL1:   0x12,
		PowerL2:   0x14,
		PowerL3:   0x16,
		Power:     0x28,
		CosphiL1:  0x2E,
		CosphiL2:  0x2F,
		CosphiL3:  0x30,
		Cosphi:    0x31,
		Frequency: 0x33,
		Import:    0x34,
		ImportL1:  0x40,
		ImportL2:  0x42,
		ImportL3:  0x44,
		Export:    0x4E,
	}
	return &CarloGavazziEM24_E1Producer{Opcodes: ops}
}

// Description implements Producer interface
func (p *CarloGavazziEM24_E1Producer) Description() string {
	return "Carlo Gavazzi EM24_E1"
}

func (p *CarloGavazziEM24_E1Producer) snip16(iec Measurement, scaler ...float64) Operation {
	transform := RTUInt16ToFloat64 // default conversion
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}

	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   1,
		IEC61850:  iec,
		Transform: transform,
	}
	return operation
}

func (p *CarloGavazziEM24_E1Producer) snip32(iec Measurement, scaler ...float64) Operation {
	transform := RTUInt32ToFloat64Swapped // default conversion
	if len(scaler) > 0 {
		transform = MakeScaledTransform(transform, scaler[0])
	}

	operation := Operation{
		FuncCode:  ReadInputReg,
		OpCode:    p.Opcode(iec),
		ReadLen:   2,
		IEC61850:  iec,
		Transform: transform,
	}
	return operation
}

// Probe implements Producer interface
func (p *CarloGavazziEM24_E1Producer) Probe() Operation {
	return p.snip32(VoltageL1, 10)
}

// Produce implements Producer interface
func (p *CarloGavazziEM24_E1Producer) Produce() (res []Operation) {
	for _, op := range []Measurement{
		VoltageL1, VoltageL2, VoltageL3,
	} {
		res = append(res, p.snip32(op, 10))
	}

	for _, op := range []Measurement{
		CurrentL1, CurrentL2, CurrentL3,
	} {
		res = append(res, p.snip32(op, 1000))
	}

	for _, op := range []Measurement{
		Cosphi, CosphiL1, CosphiL2, CosphiL3,
	} {
		res = append(res, p.snip16(op, 1000))
	}

	for _, op := range []Measurement{
		Frequency,
	} {
		res = append(res, p.snip16(op, 10))
	}

	for _, op := range []Measurement{
		Power, PowerL1, PowerL2, PowerL3,
	} {
		res = append(res, p.snip32(op, 10))
	}

	for _, op := range []Measurement{
		Import, ImportL1, ImportL2, ImportL3,
		Export,
	} {
		res = append(res, p.snip32(op, 10))
	}

	return res
}
