# goldmark-attributes
[GoldMark](https://github.com/yuin/goldmark/) block attributes extension.

```markdown
# Document {#main}

{.epigraph}
> Why, you may take the most gallant sailor, the most intrepid airman or the most audacious soldier, put them at a table together – what do you get? The sum of their fears.
```

```html
<h1 id="main">Document</h1>
<blockquote class="epigraph"><p>Why, you may take the most gallant sailor, the most intrepid airman or the most audacious soldier, put them at a table together – what do you get? The sum of their fears.</p>
</blockquote>
```

```go
var md = goldmark.New(
    goldmark.WithExtensions(attrs.BlockAttributes()),
)
var source = []byte("{#id .class1}\ntext")
err := md.Convert(source, os.Stdout)
if err != nil {
    log.Fatal(err)
}
```