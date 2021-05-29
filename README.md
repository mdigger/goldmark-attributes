# goldmark-attributes
[![GoDoc](https://godoc.org/github.com/mdigger/goldmark-attributes?status.svg)](https://godoc.org/github.com/mdigger/goldmark-attributes)

[GoldMark](https://github.com/yuin/goldmark/) block attributes extension.

```markdown
# Document {#main}

> Why, you may take the most gallant sailor, the most intrepid airman or the
> most audacious soldier, put them at a table together – what do you get? The
> sum of their fears.
{.epigraph}
```

```html
<h1 id="main">Document</h1>
<blockquote class="epigraph"><p>Why, you may take the most gallant sailor, the
most intrepid airman or the most audacious soldier, put them at a table 
together – what do you get? The sum of their fears.</p>
</blockquote>
```

```go
var md = goldmark.New(attributes.Enable)
var source = []byte("{#id .class1}\ntext")
err := md.Convert(source, os.Stdout)
if err != nil {
    log.Fatal(err)
}
```