package datasize

import (
	"fmt"
	"strconv"
	"strings"
)

type ByteSize uint64

const (
	B  ByteSize = 1
	KB          = B << 10
	MB          = KB << 10
	GB          = MB << 10
	TB          = GB << 10
	PB          = TB << 10
	EB          = PB << 10

	fnUnmarshalText string = "UnmarshalText"
	maxUint64       uint64 = (1 << 64) - 1
	cutoff          uint64 = maxUint64 / 10
)

func (b ByteSize) Bytes() uint64 {
	return uint64(b)
}

func (b ByteSize) KBytes() float64 {
	return float64(b) / float64(KB)
}

func (b ByteSize) MBytes() float64 {
	return float64(b) / float64(MB)
}

func (b ByteSize) GBytes() float64 {
	return float64(b) / float64(GB)
}

func (b ByteSize) TBytes() float64 {
	return float64(b) / float64(TB)
}

func (b ByteSize) PBytes() float64 {
	return float64(b) / float64(PB)
}

func (b ByteSize) EBytes() float64 {
	return float64(b) / float64(EB)
}

func (b ByteSize) String() string {
	switch {
	case b == 0:
		return fmt.Sprint("0B")
	case b%EB == 0:
		return fmt.Sprintf("%dEB", b/EB)
	case b%PB == 0:
		return fmt.Sprintf("%dPB", b/PB)
	case b%TB == 0:
		return fmt.Sprintf("%dTB", b/TB)
	case b%GB == 0:
		return fmt.Sprintf("%dGB", b/GB)
	case b%MB == 0:
		return fmt.Sprintf("%dMB", b/MB)
	case b%KB == 0:
		return fmt.Sprintf("%dKB", b/KB)
	default:
		return fmt.Sprintf("%dB", b)
	}
}

func (b ByteSize) HR() string {
	return b.HumanReadable()
}

func (b ByteSize) HumanReadable() string {
	switch {
	case b > EB:
		return fmt.Sprintf("%f.1 EB", b.EBytes())
	case b > PB:
		return fmt.Sprintf("%f.1 PB", b.PBytes())
	case b > TB:
		return fmt.Sprintf("%f.1 TB", b.TBytes())
	case b > GB:
		return fmt.Sprintf("%f.1 GB", b.GBytes())
	case b > MB:
		return fmt.Sprintf("%f.1 MB", b.MBytes())
	case b > KB:
		return fmt.Sprintf("%f.1 KB", b.KBytes())
	default:
		return fmt.Sprintf("%d B", b)
	}
}

func (b ByteSize) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

func (b *ByteSize) UnmarshalText(t []byte) error {
	var val uint64
	var unit string

	// copy for error message
	t0 := t

	var c byte
	var i int

ParseLoop:
	for i < len(t) {
		c = t[i]
		switch {
		case '0' <= c && c <= '9':
			if val > cutoff {
				goto Overflow
			}

			c = c - '0'
			val *= 10

			if val > val+uint64(c) {
				// val+v overflows
				goto Overflow
			}
			val += uint64(c)
			i++

		default:
			if i == 0 {
				goto SyntaxError
			}
			break ParseLoop
		}
	}

	unit = strings.TrimSpace(string(t[i:]))
	switch unit {
	case "", "B", "b":
		// do nothing - already in bytes

	case "K", "KB", "k", "kb", "kB":
		if val > maxUint64/uint64(KB) {
			goto Overflow
		}
		val *= uint64(KB)

	case "M", "MB", "m", "mb", "mB":
		if val > maxUint64/uint64(MB) {
			goto Overflow
		}
		val *= uint64(MB)

	case "G", "GB", "g", "gb", "gB":
		if val > maxUint64/uint64(GB) {
			goto Overflow
		}
		val *= uint64(GB)

	case "T", "TB", "t", "tb", "tB":
		if val > maxUint64/uint64(TB) {
			goto Overflow
		}
		val *= uint64(TB)

	case "P", "PB", "p", "pb", "pB":
		if val > maxUint64/uint64(PB) {
			goto Overflow
		}
		val *= uint64(PB)

	case "E", "EB", "e", "eb", "eB":
		if val > maxUint64/uint64(EB) {
			goto Overflow
		}
		val *= uint64(EB)

	default:
		goto SyntaxError
	}

	*b = ByteSize(val)
	return nil

Overflow:
	*b = ByteSize(maxUint64)
	return &strconv.NumError{fnUnmarshalText, string(t0), strconv.ErrRange}

SyntaxError:
	*b = 0
	return &strconv.NumError{fnUnmarshalText, string(t0), strconv.ErrSyntax}
}
