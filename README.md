# CFGO

Very limited CFG generator for Golang

* if
* for
* range
* return/break/continue

## Use

```sh
go build main.go && ./main samples/if.go && dot cfg.dot -Tsvg > cfg.svg
```

Generates Dot file that can be converted to an svg image using Graphviz

## Example

![example.svg](example.svg)
