package color_test

import (
	"encoding/json"
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

func ExampleRGB_MarshalJSON() {
	var c, d color.RGB
	for _, s := range []string{
		"#000",
		"#123",
		"#fff",
		"#012345",
		"#543210",
	} {
		c.Set(s)
		b, err := json.Marshal(c)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", b)
		err = json.Unmarshal(b, &d)
		if err != nil {
			panic(err)
		}
		if c != d {
			panic("roundtrip failure!")
		}
	}

	// Output:
	// [0,0,0]
	// [17,34,51]
	// [255,255,255]
	// [1,35,69]
	// [84,50,16]
}
