# Ellie
Ellie - simple way of transforming data

## Spaces
Space is the cardinality of the underlying lattice. Adding a new computed column doesnt change the space.

Space can be expanded for example by doing a cartesian product.
TODO: Example here

Space can be reduced by taking unique values of a subset of columns.
TODO: Example here


## Example: 1000 founders use case
```
// Databases
-> db.Messages = wmUserId, type, author, id_id, id_remote, t, quotedStanzaID, body

```


## How to run lexer

```bash
$ antlr4-parse ellie.g4 prog -gui
$ antlr4-parse ellie.g4 prog -tree

$ antlr4 -Dlanguage=Go -o parser ellie.g4
```


# Roadmap
- [x] add support for more functions - IF, LOOKUP ...
- [x] tables and db's should be resolved in context of the current row
- [x] support for JOIN functionality using LOOKUP and other complex cases
- [ ] Split a column value by delimiter and duplicate rows (dynamic cartesian)
- [ ] Stats over rolling window (last 10 / top 10 entries)
- [ ] convert parsed tree into a postfix expression that can easily be evaluated
- [ ] profile and speed up computation
- [ ] space reduction & expansion
- [ ] Load dbs from a query / datasource
- [ ] Incremental sync for dbs
- [ ] Charts


# One column at a time
Benefits:
- makes it super easy to understand whats happening

Cons:
- repetitive if has to be done for multiple cols

Maybe its easier to start with 1 column at a time, and then expand syntax to allow importing multiple columns:

```
tbl1.col2 = FILTER( tbl2.col2, tbl2.col3 == tbl1.col1 )
```