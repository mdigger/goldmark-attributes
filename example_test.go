package attributes_test

import (
	"log"
	"os"

	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
)

func Example() {
	var source = []byte(`
text
{#id .class}
`)
	var md = goldmark.New(attributes.Enable)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	// Output: <p id="id" class="class">text</p>
}
