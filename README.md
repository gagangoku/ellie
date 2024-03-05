# Ellie
Ellie - simple way of transforming data at any scale

## How it works
Ellie works with databases and tables. A database is nothing but a simple immutable table.
A table represents a columnar data structure - named columns, and rows.

To keep things simple, you have full access to other columns in a row within an instruction.

## Spaces
Space represents the cardinality of the data. If you take the unique values in each column vector, then space represents a collection of columns that form the table. Adding a new computed column to the table doesnt change the space.

### Expansion
Space can be expanded for example by doing a cartesian product. Example:
```go
db.days = CSV(`dayId
1
2
3
`)
db.weeks = CSV(`weekId
w1
w2
`)

tbl.Master = CARTESIAN(db.days.dayId, db.weeks.weekId)
tbl.Master.col2 = 1

dayId,weekId,col2
1,w1,1
1,w2,1
2,w1,1
2,w2,1
3,w1,1
3,w2,1
```

Another common example is concatenating multiple data sources over the same space:
```go
tbl.ShortBatch = CSV(`userId,joinDate`)
tbl.ShortBatch.duration = 10

tbl.LongBatch = CSV(`userId,joinDate`)
tbl.LongBatch.duration = 100

tbl.AllBatches = CONCAT_TABLE(tbl.ShortBatch, tbl.LongBatch)
```


### Contraction
Space can be reduced by taking unique values of a subset of columns. Example:

```go
db.Master = CSV(`dayId,weekId,col2
1,w1,1
1,w2,1
2,w1,1
2,w2,1
3,w1,1
3,w2,1
`)

tbl.days = CARTESIAN(db.Master.dayId)
```


# Roadmap
- [x] tables and db's should be resolved in context of the current row
- [x] add support for more functions - IF, LOOKUP ...
- [x] support for JOIN functionality using LOOKUP and other complex cases
- [x] space reduction & expansion
- [ ] Split a column value by delimiter and duplicate rows (dynamic cartesian)
- [ ] Stats over rolling window (last 10 / top 10 entries)
- [ ] convert parsed tree into a postfix expression that can easily be evaluated
- [ ] profile and speed up computation
- [ ] Load dbs from a query / datasource
- [ ] Incremental sync for dbs
- [ ] Charts


# One column at a time
One way of modeling is to have multiple expressions for a column, for example setting a set of values where a filter matches. An example of this would be to set:

```go
tbl.UserDetails = CSV(`userId,joinDate,batchType`)
tbl.UserDetails.duration = 0
UPDATE UserDetails.duration=100 IF UserDetails.batchType == 'long'
UPDATE UserDetails.duration=10 IF UserDetails.batchType == 'short'
```

Alternately:
```go
tbl.UserDetails = CSV(`userId,joinDate,batchType`)
tbl.UserDetails.duration = IF ( tbl.UserDetails.batchType == 'long', 100, IF ( tbl.UserDetails.batchType == 'short', 10, 0 ) )
```


Ellie supports only 1 expression per column. The reason is that its a lot easier to understand data lineage this way, and forces users to know exactly how data transformations should happen.
Nested IF conditions are hard to read. So there is scope here for a better switch case like syntax.

In excel world, the equivalent is to create additional columns (and hide them) and use them to write a readable formula.

Benefits of 1 expression per column:
- makes it super easy to understand whats happening
- makes data lineage simple
- reordering instructions will not cause problems

Cons:
- nested IF's are hard to read
- repetitive if has to be done for multiple cols


# Lookups
Simple data transformations within the same row are easy to write. However, more often than not you need to reference a value in another table, or aggregate it.

For example:
```go
tbl.Master = CSV(`dateTime,userId,message,
1 Jan 3pm, u1, hi
1 Jan 4pm, u1, hello
2 Jan 11am, u2, test message 1
5 Jan 11am, u3, test message 2
`)
```

You may wish to count the number of messages per user per day:
```go
tbl.Master.dayId = DAY(tbl.Master.dateTime)  // calculate day

tbl.NumMsgsPerUser = CARTESIAN(tbl.Master.dayId, tbl.Master.userId)
tbl.NumMsgsPerUser.numMsgs = COUNT( LOOKUP( tbl.Master, tbl.Master.dayId == tbl.NumMsgsPerUser.dayId, tbl.Master.userId == tbl.NumMsgsPerUser.userId ))
```


# How to run lexer

```bash
$ antlr4-parse ellie.g4 prog -gui
$ antlr4-parse ellie.g4 prog -tree

$ antlr4 -Dlanguage=Go -o parser ellie.g4
```
