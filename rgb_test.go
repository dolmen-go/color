package color_test

import (
	"fmt"

	"github.com/dolmen-go/color"
)

func ExampleRGB_String() {
	fmt.Println(color.RGB{0, 0, 0})
	fmt.Println(color.RGB{255, 255, 255})
	fmt.Println(color.RGB{0x01, 0x23, 0x45})

	// Output:
	// #000
	// #fff
	// #012345
}

func ExampleRGB_Set() {
	var c color.RGB
	for _, s := range []string{
		"000",
		"#000",
		"123",
		"fff",
		"012345",
		"#543210",
	} {
		c.Set(s)
		fmt.Println(c)
	}

	// Output:
	// #000
	// #000
	// #123
	// #fff
	// #012345
	// #543210
}
