package attributes

import (
	"log"
	"os"
	"testing"

	"github.com/yuin/goldmark"
)

func TestAttributes(t *testing.T) {
	source := []byte(`
Paragraph with attributes.
{.myPar1}

> Paragraph with attributes inside a block quote.
> {.myPar2}

> Blockquote with attributes.
{.myBlockquote1}

- list with
- attributes
{.myList1}

and now:

- a loose list

- with attributes

{.myList2}

another:

- loose list

- with attributes
{.myList3}

and now:

- a loose list where

- the last paragraph has attributes.
	Note that the indentation of the attribute block is significant.
	{.myPar3}

and finally:

- > a list where each
	> item is a blockquote
	> {.myPar4}

- > to see that everything is possible
	> {.myPar5}
	{.myBlockquote2}
{.myList4}
`)

	var md = goldmark.New(Enable)
	err := md.Convert(source, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}

}
