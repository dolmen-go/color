package color

import (
	"errors"
	"strconv"
)

type RGB struct {
	R, G, B uint8
}

const hex = "0123456789abcdef"

// sucess is a simple propagator of success.
type success bool

const ok success = true

func (ok success) scan(b byte) (uint8, success) {
	if !ok {
		return 0, false
	}
	switch {
	case b <= '9':
		if b < '0' {
			return 0, false
		}
		b -= '0'
	case b >= 'a':
		if b > 'f' {
			return 0, false
		}
		b -= 'a' - 10
	case b >= 'A':
		if b > 'F' {
			return 0, false
		}
		b -= 'A' - 10
	}
	return b, true
}

func (ok success) h(b byte, target *uint8) success {
	n, ok := ok.scan(b)
	if !ok {
		return ok
	}
	*target = (n << 4) | n
	return ok
}

func (ok success) hh(b1, b2 byte, target *uint8) success {
	n1, ok := ok.scan(b1)
	if !ok {
		return ok
	}
	n2, ok := ok.scan(b2)
	if !ok {
		return ok
	}
	*target = (n1 << 4) | n2
	return ok
}

func (c RGB) String() string {
	if (c.R>>4) == (c.R&15) && (c.G>>4) == (c.G&15) && (c.B>>4) == (c.B&15) {
		return string([]byte{'#', hex[c.R&15], hex[c.G&15], hex[c.B&15]})
	}
	return string([]byte{'#', hex[c.R>>4], hex[c.R&15], hex[c.G>>4], hex[c.G&15], hex[c.B>>4], hex[c.B&15]})
}

var ErrInvalidRGB = errors.New("invalid RGB value")

func (c *RGB) Set(s string) error {
	switch len(s) {
	case 4:
		if s[0] != '#' {
			return ErrInvalidRGB
		}
		s = s[1:]
		fallthrough
	case 3:
		if ok.h(s[0], &(*c).R).h(s[1], &(*c).G).h(s[2], &(*c).B) {
			return nil
		}
	case 7:
		if s[0] != '#' {
			return ErrInvalidRGB
		}
		s = s[1:]
		fallthrough
	case 6:
		if ok.hh(s[0], s[1], &(*c).R).hh(s[2], s[3], &(*c).G).hh(s[4], s[5], &(*c).B) {
			return nil
		}
	}
	return ErrInvalidRGB
}

func (c RGB) MarshalText() ([]byte, error) {
	if (c.R>>4) == (c.R&15) && (c.G>>4) == (c.G&15) && (c.B>>4) == (c.B&15) {
		return []byte{hex[c.R&15], hex[c.G&15], hex[c.B&15]}, nil
	}
	return []byte{hex[c.R>>4], hex[c.R&15], hex[c.G>>4], hex[c.G&15], hex[c.B>>4], hex[c.B&15]}, nil
}

func (c *RGB) UnmarshalText(s []byte) error {
	switch len(s) {
	case 3:
		if ok.h(s[0], &(*c).R).h(s[1], &(*c).G).h(s[2], &(*c).B) {
			return nil
		}
	case 6:
		if ok.hh(s[0], s[1], &(*c).R).hh(s[2], s[3], &(*c).G).hh(s[4], s[5], &(*c).B) {
			return nil
		}
	}
	return ErrInvalidRGB
}

func (c RGB) MarshalJSON() ([]byte, error) {
	b := make([]byte, 1, 1+3+1+3+1+3+1)
	b[0] = '['
	b = strconv.AppendUint(b, uint64(c.R), 10)
	b = strconv.AppendUint(b, uint64(c.G), 10)
	b = strconv.AppendUint(b, uint64(c.B), 10)
	b = append(b, ']')
	return b, nil
}
