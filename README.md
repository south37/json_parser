## json_parser
Simple implementation of LL(1) parser for JSON.

## WHY
For my practice of Go-lang.

## BNF

```
json = array | object
array = [ elements ]
elements = expression | expression , elements
object = { key_value_pairs }
key_value_pairs = key_value_pair | key_value_pair , key_value_pairs
key_value_pair = key : expression
key = string | number
expression = array | object | string | number | null
string = "([^"]|\\")*"
number = 0|-?[1-9][0-9]*
null = null
```

## TODO
- Support unclosed object(array).
- Check full usage of input.
- Implement extracting of value.
