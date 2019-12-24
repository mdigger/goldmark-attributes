package attributes_test

import (
	"log"
	"os"

	attributes "github.com/mdigger/goldmark-attributes"
	"github.com/yuin/goldmark"
)

func Example() {
	var md = goldmark.New(attributes.Enable())
	var source = []byte("{#id .class1}\ntext")
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	// Output: <p id="id" class="class1">text</p>
}
