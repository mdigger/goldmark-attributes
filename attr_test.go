package attrs_test

import (
	"log"
	"os"

	attrs "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
)

func Example() {
	var md = goldmark.New(
		goldmark.WithExtensions(attrs.BlockAttributes()),
	)
	var source = []byte("{#id .class1}\ntext")
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	// Output: <p id="id" class="class1">text</p>
}
