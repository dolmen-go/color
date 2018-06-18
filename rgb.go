package color

import (
	"encoding/json"
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

func (ok success) nibble(b byte) (uint8, success) {
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

func (ok success) h(b string, target []uint8) success {
	for len(b) >= 1 {
		n, ok := ok.nibble(b[0])
		if !ok {
			return ok
		}
		target[0] = (n << 4) | n
		target = target[1:]
		b = b[1:]
	}
	return ok
}

func (ok success) hh(src string, target []uint8) success {
	for len(src) >= 2 {
		n1, ok := ok.nibble(src[0])
		if !ok {
			return ok
		}
		n2, ok := ok.nibble(src[1])
		if !ok {
			return ok
		}
		target[0] = (n1 << 4) | n2
		target = target[1:]
		src = src[2:]
	}
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
	var done success
	var b [3]uint8
	switch len(s) {
	case 4:
		if s[0] != '#' {
			return ErrInvalidRGB
		}
		s = s[1:]
		fallthrough
	case 3:
		done = ok.h(s, b[:])
	case 7:
		if s[0] != '#' {
			return ErrInvalidRGB
		}
		s = s[1:]
		fallthrough
	case 6:
		done = ok.hh(s, b[:])
	}
	if !done {
		return ErrInvalidRGB
	}
	(*c).R, (*c).G, (*c).B = b[0], b[1], b[2]
	return nil
}

func (c RGB) MarshalText() ([]byte, error) {
	if (c.R>>4) == (c.R&15) && (c.G>>4) == (c.G&15) && (c.B>>4) == (c.B&15) {
		return []byte{hex[c.R&15], hex[c.G&15], hex[c.B&15]}, nil
	}
	return []byte{hex[c.R>>4], hex[c.R&15], hex[c.G>>4], hex[c.G&15], hex[c.B>>4], hex[c.B&15]}, nil
}

func (c *RGB) UnmarshalText(s []byte) error {
	var done success
	var b [3]uint8
	switch len(s) {
	case 3:
		done = ok.h(string(s), b[:])
	case 6:
		done = ok.hh(string(s), b[:])
	}

	if !done {
		return ErrInvalidRGB
	}

	(*c).R, (*c).G, (*c).B = b[0], b[1], b[2]
	return nil
}

func (c RGB) MarshalJSON() ([]byte, error) {
	b := make([]byte, 1, 1+3+1+3+1+3+1)
	b[0] = '['
	b = strconv.AppendUint(b, uint64(c.R), 10)
	b = append(b, ',')
	b = strconv.AppendUint(b, uint64(c.G), 10)
	b = append(b, ',')
	b = strconv.AppendUint(b, uint64(c.B), 10)
	b = append(b, ']')
	return b, nil
}

func (c *RGB) UnmarshalJSON(b []byte) error {
	var a [3]uint16 // We can't directly use [3]uint8 because this is a synonym of [3]byte
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	if a[0] > 255 {
		return errors.New("invalid red value")
	}
	if a[1] > 255 {
		return errors.New("invalid green value")
	}
	if a[2] > 255 {
		return errors.New("invalid blue value")
	}
	c.R = uint8(a[0])
	c.G = uint8(a[1])
	c.B = uint8(a[2])
	return nil
}
