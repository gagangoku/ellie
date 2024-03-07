package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
)

func Test_Engine_Basic(t *testing.T) {
	solver, _ := setupEngine("", "")

	var lhs *Node
	var code int
	var err error

	// Should get an error since tbl.test is not initialized
	_, err = solver.Evaluate("tbl.test._c1 = ((1+2)*3)/4")
	if err == nil {
		t.Fatalf("Expected error, got %s", lhs.Text)
	}

	code, err = solver.Evaluate("tbl.test = CSV(`id,Col2\n0,hi\n1,there`)")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if code != CODE_OK {
		t.Fatalf("Bad code: %v", code)
	}

	lhs, _, _ = solver.Parse("tbl.test._row = ROW()+2")
	if lhs.Text != "tbl.test._row" {
		t.Fatalf("Got %s", lhs.Text)
	}

	lhs, _, err = solver.Parse("tbl.test._c1 = ((1+2)*3)/4")
	noop(err)
	if lhs.Text != "tbl.test._c1" {
		t.Fatalf("Got %s", lhs.Text)
	}

	code, err = solver.Evaluate("tbl.test._c1 = ((1+2)*3)/4")
	noop(err)
	if err != nil || code != CODE_OK {
		t.Fatalf("Expected success")
	}

	code, err = solver.Evaluate("tbl.test._c1 = 2")
	noop(err)
	if err == nil || code == CODE_OK {
		t.Fatalf("Expected failure")
	}

	code, _ = solver.Evaluate("tbl.test.Col3 = 'Col3'")
	if code != CODE_OK {
		t.Fatalf("Bad code: %v", code)
	}

	code, _ = solver.Evaluate("tbl.test.Col4 = 1")
	if code != CODE_OK {
		t.Fatalf("Bad code: %v", code)
	}

	lhs, _, _ = solver.Parse("tbl.test._c4 = tbl.test.Col2 + 1.1 + 'hello'")
	if lhs.Text != "tbl.test._c4" {
		t.Fatalf("Got %s", lhs.Text)
	}

	// Reference to a non-existing column
	_, err = solver.Evaluate("tbl.test._c2 = tbl.test.Col1")
	noop(err)
	if err == nil {
		t.Fatalf("Expected error")
	}

	lhs, _, _ = solver.Parse("tbl.test._c2 = CONCAT(tbl.test.Col3, CONCAT(tbl.test.Col2, 'hello'))")
	if lhs.Text != "tbl.test._c2" {
		t.Fatalf("Got %s", lhs.Text)
	}

	lhs, _, _ = solver.Parse("tbl.test._c3 = IF(tbl.test.Col3 == 'hello', 1, 2) + 0.5")
	if lhs.Text != "tbl.test._c3" {
		t.Fatalf("Got %s", lhs.Text)
	}

	// IF string comparison
	lhs, _, _ = solver.Parse("tbl.test._c5 = IF(tbl.test.Col2 == 'there', 1, 0)")
	if lhs.Text != "tbl.test._c5" {
		t.Fatalf("Got %s", lhs.Text)
	}

	noop(solver.tables)
	got := strings.Join(solver.prettyTable("test", ",", false), "\n")
	expected := trimLines(`id,Col2,_c1,Col3,Col4
	0,hi,2.25,Col3,1
	1,there,2.25,Col3,1`)
	if expected != got {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_MathFunctions(t *testing.T) {
	solver, _ := setupEngine("", "")

	// numbers db
	team := `col1,col2
		1,-2
		3,4
		`
	solver.Evaluate(fmt.Sprintf("db.Numbers = CSV(`%s`)", trimLines(team)))

	solver.Evaluate("tbl.Numbers = db.Numbers")
	solver.Evaluate("tbl.Numbers.col3 = (tbl.Numbers.col1 * 2) + (tbl.Numbers.col2 - 1) - 3 + IF(tbl.Numbers.col2 < 0 , -1, 1)")

	noop(solver.tables)
	got := strings.Join(solver.prettyTable("Numbers", ",", false), "\n")
	expected := trimLines(`col1,col2,col3
	1,-2,-5
	3,4,7`)
	if expected != got {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_AndOrExpression(t *testing.T) {
	solver, _ := setupEngine("", "")

	// numbers db
	team := `col1,col2
		1,-2
		3,4
		1,4
		2,2
		`
	solver.Evaluate(fmt.Sprintf("db.Numbers = CSV(`%s`)", trimLines(team)))

	solver.Evaluate("tbl.Numbers = db.Numbers")
	solver.Evaluate("tbl.Numbers.col3 = IF ( (tbl.Numbers.col1 == 1) AND (tbl.Numbers.col2 == -2) , 1, 0 )")
	solver.Evaluate("tbl.Numbers.col4 = IF ( (tbl.Numbers.col1 == '') OR (tbl.Numbers.col2 == 4) , 1, 0 )")

	noop(solver.tables)
	got := strings.Join(solver.prettyTable("Numbers", ",", false), "\n")
	expected := trimLines(`col1,col2,col3,col4
	1,-2,1,0
	3,4,0,1
	1,4,0,1
	2,2,0,0`)
	if expected != got {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_Filter(t *testing.T) {
	solver, _ := setupEngine("", "")

	var lhs string
	var res Result
	var code int
	var err error
	noop(lhs, res, code, err)

	// Messages db
	msgs := `userId,groupId,user,groupName,msgId,type,timestamp,time,date,month,body,quotedMsg
	u1,g1,user1,group1,msgId1,text,1680275430000,31 March 2023 20:40:30,31,Mar,hello,
	u1,g1,user1,group1,msgId2,text,1680275430000,31 March 2023 20:40:30,31,Mar,hi,
	u1,g1,user1,group1,msgId3,image,1680275430000,31 March 2023 20:40:30,31,Mar,,
	u2,g2,user2,group2,msgId4,image,1680275430000,31 March 2023 20:40:30,31,Mar,,
	`
	code, err = solver.Evaluate(fmt.Sprintf("db.Messages = CSV(`%s`)", trimLines(msgs)))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if code != CODE_OK {
		t.Fatalf("Bad code: %v", code)
	}

	// Team db
	team := `center,groupId,batch,Nutritionist,fromDate,toDate
		center1,g1,Aug 15,Nut1,2023-01-01,2023-02-02
		center2,g2,Aug 15,Nut2,2023-01-01,2023-02-02
		`
	code, err = solver.Evaluate(fmt.Sprintf("db.Team = CSV(`%s`)", trimLines(team)))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if code != CODE_OK {
		t.Fatalf("Bad code: %v", code)
	}

	// Image tags db
	imageTags := `msgId,tags
		msgId3,Grocery|Food
		msgId4,Person|abcd
		msgId5,
		`
	code, err = solver.Evaluate(fmt.Sprintf("db.ImageTags = CSV(`%s`)", trimLines(imageTags)))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if code != CODE_OK {
		t.Fatalf("Bad code: %v", code)
	}

	// Master table
	solver.Evaluate("tbl.Master = db.Messages")

	// LOOKUP - multiple tables not supported
	code, err = solver.Evaluate("tbl.Master.batchId = FIRST ( LOOKUP ( db.Team.Center , db.Team1.groupId == tbl.Master.groupId ) )")
	noop(err)
	if err == nil || code == CODE_OK {
		t.Fatalf("Expected error")
	}

	// LOOKUP - single table
	solver.Evaluate("tbl.Master.batchId = FIRST ( LOOKUP ( db.Team.batch , db.Team.groupId == tbl.Master.groupId ) )")
	solver.Evaluate("tbl.Master.Nutritionist = FIRST ( LOOKUP ( db.Team.Nutritionist , tbl.Master.groupId == db.Team.groupId ) )")

	// TODO: date conditions check
	// solver.Execute("tbl.Master.Nutritionist = FIRST ( LOOKUP ( db.Team.Nutritionist , tbl.Master.groupName == db.Team.center , tbl.Master.dateTime >= db.Team.fromDate , tbl.Master.dateTime <= db.Team.toDate ) )")

	got := strings.Join(solver.dumpAllTablesAndDBs(false), "\n")
	expected := `DB: ImageTags
	msgId,tags
	msgId3,Grocery|Food
	msgId4,Person|abcd
	msgId5,

	DB: Messages
	userId,groupId,user,groupName,msgId,type,timestamp,time,date,month,body,quotedMsg
	u1,g1,user1,group1,msgId1,text,1680275430000,31 March 2023 20:40:30,31,Mar,hello,
	u1,g1,user1,group1,msgId2,text,1680275430000,31 March 2023 20:40:30,31,Mar,hi,
	u1,g1,user1,group1,msgId3,image,1680275430000,31 March 2023 20:40:30,31,Mar,,
	u2,g2,user2,group2,msgId4,image,1680275430000,31 March 2023 20:40:30,31,Mar,,

	DB: Team
	center,groupId,batch,Nutritionist,fromDate,toDate
	center1,g1,Aug 15,Nut1,2023-01-01,2023-02-02
	center2,g2,Aug 15,Nut2,2023-01-01,2023-02-02

	Table: Master
	userId,groupId,user,groupName,msgId,type,timestamp,time,date,month,body,quotedMsg,batchId,Nutritionist
	u1,g1,user1,group1,msgId1,text,1680275430000,31 March 2023 20:40:30,31,Mar,hello,,Aug 15,Nut1
	u1,g1,user1,group1,msgId2,text,1680275430000,31 March 2023 20:40:30,31,Mar,hi,,Aug 15,Nut1
	u1,g1,user1,group1,msgId3,image,1680275430000,31 March 2023 20:40:30,31,Mar,,,Aug 15,Nut1
	u2,g2,user2,group2,msgId4,image,1680275430000,31 March 2023 20:40:30,31,Mar,,,Aug 15,Nut2
	`
	if trimLines(expected) != trimLines(got) {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_Filter_Parallel(t *testing.T) {
	solver, _ := setupEngine("", "")
	solver.nRowsPerLookupBatch = 5

	// db 1
	msgs := `idx,col1
	1,a
	2,b
	3,c
	4,d
	5,e
	6,f
	7,g
	8,h
	9,i
	10,j
	11,k
	`
	solver.Evaluate(fmt.Sprintf("db.db1 = CSV(`%s`)", trimLines(msgs)))

	// LOOKUP - multiple conditions
	solver.Evaluate("tbl.Master = db.db1")
	solver.Evaluate("tbl.Master.col2 = FIRST ( LOOKUP ( db.db1.col1 , db.db1.col1 == tbl.Master.col1 ) )")

	got := strings.Join(solver.prettyTable("Master", ",", false), "\n")
	expected := `idx,col1,col2
	1,a,a
	2,b,b
	3,c,c
	4,d,d
	5,e,e
	6,f,f
	7,g,g
	8,h,h
	9,i,i
	10,j,j
	11,k,k
	`
	if trimLines(expected) != trimLines(got) {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_FF_LOOKUP(t *testing.T) {
	solver, _ := setupEngine("", "")

	var lhs string
	var res Result
	var code int
	var err error
	noop(lhs, res, code, err)

	// db 1
	msgs := `idx,col1,col2,float,date
	1,a1,a2,1.001,2023-01-01
	2,b1,b2,2.5,2023-01-02
	3,c1,c2,2.6,2023-01-03
	4,d1,,2.6,2023-01-03
	`
	solver.Evaluate(fmt.Sprintf("db.db1 = CSV(`%s`)", trimLines(msgs)))

	// Team db
	team := `col2,col3,col4,col5
		a2,x1,1,1
		a2,x2,2,2
		a2,x3,3,3
		b2,x4,4,4
		d2,x5,5,5
		`
	solver.Evaluate(fmt.Sprintf("db.db2 = CSV(`%s`)", trimLines(team)))

	// LOOKUP - multiple conditions
	solver.Evaluate("tbl.Master = db.db1")
	solver.Evaluate("tbl.Master.ceil = CEIL ( tbl.Master.float ) ")
	solver.Evaluate("tbl.Master.epochMs = EPOCH_MS ( tbl.Master.date ) ")
	solver.Evaluate("tbl.Master._c4 = LAST ( LOOKUP ( db.db2.col4 , db.db2.col2 == tbl.Master.col2, db.db2.col5 <= 2 ) )")
	solver.Evaluate("tbl.Master._count = COUNT ( LOOKUP ( db.db2.col4 , db.db2.col2 == tbl.Master.col2 ) )")
	solver.Evaluate("tbl.Master._sContains = IF ( STR_CONTAINS ( tbl.Master.idx, '2' ), 1, 0 )")

	got := strings.Join(solver.dumpAllTablesAndDBs(false), "\n")
	expected := `DB: db1
	idx,col1,col2,float,date
	1,a1,a2,1.001,2023-01-01
	2,b1,b2,2.5,2023-01-02
	3,c1,c2,2.6,2023-01-03
	4,d1,,2.6,2023-01-03

	DB: db2
	col2,col3,col4,col5
	a2,x1,1,1
	a2,x2,2,2
	a2,x3,3,3
	b2,x4,4,4
	d2,x5,5,5

	Table: Master
	idx,col1,col2,float,date,ceil,epochMs,_c4,_count,_sContains
	1,a1,a2,1.001,2023-01-01,2,1672531200000,2,3,0
	2,b1,b2,2.5,2023-01-02,3,1672617600000,,1,1
	3,c1,c2,2.6,2023-01-03,3,1672704000000,,0,0
	4,d1,,2.6,2023-01-03,3,1672704000000,,0,0
	`
	if trimLines(expected) != trimLines(got) {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_NewFunctions(t *testing.T) {
	solver, _ := setupEngine("", "")

	l1 := `-> db.Weeks_52 = XRANGE ( 'weekId', 1, 52 )
			-> db.Days_52 = XRANGE ( 'dayId', 1, 10 )
			-> tbl._tb1 = CARTESIAN ( db.Weeks_52.weekId , db.Days_52.dayId )
			-> tbl.tb1 = FILTER ( tbl._tb1, tbl._tb1.weekId == 1 , tbl._tb1.dayId == 1 )
			-> tbl._tb2 = CARTESIAN ( db.Weeks_52.weekId , db.Days_52.dayId )
			-> tbl.tb2 = FILTER ( tbl._tb2, tbl._tb2.weekId == 2 , tbl._tb2.dayId == 2 )
			-> tbl._tb3 = CONCAT_TABLE ( tbl.tb1, tbl.tb2 )
			-> tbl._tb3.col1 = 'col1'
			-> tbl.tb3 = RETAIN_COLS ( tbl._tb3, 'col1,weekId' )
	`
	lines := strings.Split(l1, "\n")
	payload := &SolveHandlerReq{Script: lines}
	res, _, err := solver.app.solveScript(payload)
	if err != nil {
		t.Fatalf("Script eval failed: %s", err)
	}
	fmt.Println(res)

	var got, expected string
	got = strings.Join(solver.prettyTable("tb3", ",", true), "\n")
	expected = `col1,weekId
	"col1","1"
	"col1","2"`
	if trimLines(expected) != trimLines(got) {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_AntlrParser_1(t *testing.T) {
	cmd := "tbl.Master._firstNutMsgsAfterCurrentMsg = FIRST ( LOOKUP ( db.RMsg.timeSec, tbl.RMsg.author == tbl.Master._nutNumber, tbl.RMsg.type == 'chat', tbl.RMsg.timeSec > tbl.Master.timeSec ) )"
	engine := Engine{}

	lhs, rhs, _ := engine.Parse(cmd)
	fmt.Println("lhs: ", lhs)
	_bytes, err := json.MarshalIndent(rhs, "", "    ")
	noop(_bytes, err)
	fmt.Println(string(_bytes))
}

func Test_AntlrParser_Basic(t *testing.T) {
	engine := Engine{}
	engine.Init(&EtlApp{})
	var cmd string
	var err error

	_db1 := `_f1,_f2,_f3,_f4,_f5,_f6
	1,2,3,4,5,6
	1.1,2.1,3.1,4.1,5.1,6.1
	`
	cmd = fmt.Sprintf("db.db1 = CSV(`%s`)", trimLines(_db1))
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = "db.ImageTags = LOAD_CSV ('tests/data/tags.csv')"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = "tbl.Master = db.db1"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = "tbl.Master._f7 = tbl.Master._f1 + tbl.Master._f2 + tbl.Master._f3 + tbl.Master._f4 + CONCAT(tbl.Master._f5, tbl.Master._f6, 'hello') + 1.2 + TRUE"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = "tbl.Master._f8 = 'hello'"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = "tbl.Master._f9 = FIRST ( LOOKUP( db.db1._f1 , db.db1._f1 == tbl.Master._f1 ) )"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	got := strings.Join(engine.dumpAllTablesAndDBs(false), "\n")
	expected := `DB: ImageTags
	shortMsgId,tags
	3A8B9CB05FE774DFE94D,Food|fruits
	06B341CD079BC469C94F2B177C72237F,Food|roti
	xxxmsg2,clothes

	DB: db1
	_f1,_f2,_f3,_f4,_f5,_f6
	1,2,3,4,5,6
	1.1,2.1,3.1,4.1,5.1,6.1

	Table: Master
	_f1,_f2,_f3,_f4,_f5,_f6,_f7,_f8,_f9
	1,2,3,4,5,6,12.2,hello,1
	1.1,2.1,3.1,4.1,5.1,6.1,12.6,hello,1.1
	`
	if trimLines(expected) != trimLines(got) {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_Cartesian_SingleColPerGroup(t *testing.T) {
	engine := Engine{}
	engine.Init(&EtlApp{})

	db1 := `dayId
	1
	2
	3
	`
	db2 := `weekId
	a
	b
	`
	var cmd string
	var err error

	cmd = fmt.Sprintf("db.days = CSV(`%s`)", trimLines(db1))
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = fmt.Sprintf("db.weeks = CSV(`%s`)", trimLines(db2))
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = "tbl.cart = CARTESIAN( db.days.dayId, db.weeks.weekId )"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	cmd = "tbl.cart.one = 1"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	got := strings.Join(engine.dumpAllTablesAndDBs(false), "\n")
	expected := `DB: days
	dayId
	1
	2
	3

	DB: weeks
	weekId
	a
	b

	Table: cart
	dayId,weekId,one
	1,a,1
	1,b,1
	2,a,1
	2,b,1
	3,a,1
	3,b,1
	`
	if trimLines(expected) != trimLines(got) {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Engine_Cartesian_MultipleColsPerGroup(t *testing.T) {
	engine := Engine{}
	engine.Init(&EtlApp{})

	db1 := `dayId,weekId
	1,1
	2,1
	3,1
	4,2
	5,2
	6,2
	7,3
	`
	db2 := `animal,kind
	horse,herbivore
	cat,omnivore
	lion,carnivore
	tiger,carnivore
	`
	db3 := `idx
	5
	6
	`

	var cmd string
	var err error

	cmd = fmt.Sprintf("db.days = CSV(`%s`)", trimLines(db1))
	engine.Evaluate(cmd)

	cmd = fmt.Sprintf("db.animals = CSV(`%s`)", trimLines(db2))
	engine.Evaluate(cmd)

	cmd = fmt.Sprintf("db.nums = CSV(`%s`)", trimLines(db3))
	engine.Evaluate(cmd)

	cmd = "tbl.cart = CARTESIAN( db.days.dayId, db.animals.kind, db.nums.idx, db.animals.animal, db.days.weekId )"
	_, err = engine.Evaluate(cmd)
	if err != nil {
		t.Fatalf("failed with error: %s", err)
	}

	got := strings.Join(engine.dumpAllTablesAndDBs(false), "\n")
	expected := `DB: animals
	animal,kind
	horse,herbivore
	cat,omnivore
	lion,carnivore
	tiger,carnivore

	DB: days
	dayId,weekId
	1,1
	2,1
	3,1
	4,2
	5,2
	6,2
	7,3

	DB: nums
	idx
	5
	6

	Table: cart
	kind,animal,dayId,weekId,idx
	herbivore,horse,1,1,5
	herbivore,horse,1,1,6
	herbivore,horse,2,1,5
	herbivore,horse,2,1,6
	herbivore,horse,3,1,5
	herbivore,horse,3,1,6
	herbivore,horse,4,2,5
	herbivore,horse,4,2,6
	herbivore,horse,5,2,5
	herbivore,horse,5,2,6
	herbivore,horse,6,2,5
	herbivore,horse,6,2,6
	herbivore,horse,7,3,5
	herbivore,horse,7,3,6
	omnivore,cat,1,1,5
	omnivore,cat,1,1,6
	omnivore,cat,2,1,5
	omnivore,cat,2,1,6
	omnivore,cat,3,1,5
	omnivore,cat,3,1,6
	omnivore,cat,4,2,5
	omnivore,cat,4,2,6
	omnivore,cat,5,2,5
	omnivore,cat,5,2,6
	omnivore,cat,6,2,5
	omnivore,cat,6,2,6
	omnivore,cat,7,3,5
	omnivore,cat,7,3,6
	carnivore,lion,1,1,5
	carnivore,lion,1,1,6
	carnivore,lion,2,1,5
	carnivore,lion,2,1,6
	carnivore,lion,3,1,5
	carnivore,lion,3,1,6
	carnivore,lion,4,2,5
	carnivore,lion,4,2,6
	carnivore,lion,5,2,5
	carnivore,lion,5,2,6
	carnivore,lion,6,2,5
	carnivore,lion,6,2,6
	carnivore,lion,7,3,5
	carnivore,lion,7,3,6
	carnivore,tiger,1,1,5
	carnivore,tiger,1,1,6
	carnivore,tiger,2,1,5
	carnivore,tiger,2,1,6
	carnivore,tiger,3,1,5
	carnivore,tiger,3,1,6
	carnivore,tiger,4,2,5
	carnivore,tiger,4,2,6
	carnivore,tiger,5,2,5
	carnivore,tiger,5,2,6
	carnivore,tiger,6,2,5
	carnivore,tiger,6,2,6
	carnivore,tiger,7,3,5
	carnivore,tiger,7,3,6
	`
	if trimLines(expected) != trimLines(got) {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_ListsTest(t *testing.T) {
	solver, _ := setupEngine("", "")

	var err error
	_, err = solver.Evaluate("tbl.EmployeeRoster = CSV(`EmpId,Name,Mgr\ne1,Employee 1,M1\ne2,Employee 2,M2`)")
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	_, err = solver.Evaluate("tbl.EmployeeRoster.col1 = CONCAT(tbl.EmployeeRoster.Name , '-', tbl.EmployeeRoster.Mgr)")
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	got := strings.Join(solver.prettyTable("EmployeeRoster", ",", false), "\n")
	expected := trimLines(`EmpId,Name,Mgr,col1
	e1,Employee 1,M1,Employee 1-M1
	e2,Employee 2,M2,Employee 2-M2`)
	if expected != got {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_TotalArrTest(t *testing.T) {
	solver, _ := setupEngine("", "")

	var err error
	solver.Evaluate("tbl.ArrPlan = CSV(`CustomerId,Customer,Region\nc1,Customer 1,R1\nc2,Customer 2,R2`)")
	_, err = solver.Evaluate("tbl.ArrPlan.Churn_ARR = ROW() +1")
	_, err = solver.Evaluate("tbl.ArrPlan.Total_ARR = CUMULATIVE_SUM ( tbl.ArrPlan.Churn_ARR )")
	noop(err)

	got := strings.Join(solver.prettyTable("ArrPlan", ",", false), "\n")
	expected := trimLines(`CustomerId,Customer,Region,Churn_ARR,Total_ARR
	c1,Customer 1,R1,1,1
	c2,Customer 2,R2,2,3
`)
	if expected != got {
		t.Fatalf("Failed: %s", got)
	}
}

func Test_Benchmark_AntlrParser_Manual(t *testing.T) {
	cmd := "tbl.Master._firstNutMsgsAfterCurrentMsg = FIRST ( LOOKUP ( tbl.RMsg.timeSec, tbl.RMsg.author == tbl.Master._nutNumber, tbl.RMsg.type == 'chat', tbl.RMsg.timeSec > tbl.Master.timeSec ) )"
	engine := Engine{}

	startTime := time.Now()
	N := 1000 * 10
	for i := 0; i < N; i++ {
		lhs, rhs, _ := engine.Parse(cmd)
		noop(lhs, rhs)
	}
	fmt.Println("Num iterations: ", N)
	fmt.Println("Time taken milliseconds: ", time.Since(startTime).Milliseconds())
}

func Benchmark_AntlrParser(b *testing.B) {
	cmd := "tbl.Master._firstNutMsgsAfterCurrentMsg = FIRST ( LOOKUP ( tbl.RMsg.timeSec, tbl.RMsg.author == tbl.Master._nutNumber, tbl.RMsg.type == 'chat', tbl.RMsg.timeSec > tbl.Master.timeSec ) )"
	engine := Engine{}
	for i := 0; i < b.N; i++ {
		lhs, rhs, _ := engine.Parse(cmd)
		noop(lhs, rhs)
	}
}

func Benchmark_CRC(b *testing.B) {
	N := 1000 * 1000

	vals := make([]ResultValueType, 0, N)
	for i := 0; i < N; i++ {
		str := fmt.Sprintf("random string: %d", rand.Intn(N))
		vals = append(vals, str)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := Checksum(vals)
		noop(c)
	}
}

func (engine *Engine) dumpAllTablesAndDBs(quote bool) []string {
	lines := make([]string, 0)

	dbNames := lo.Keys(engine.dbs)
	sort.Strings(dbNames)
	for _, dbName := range dbNames {
		lines = append(lines, "DB: "+dbName)
		lines = append(lines, engine.prettyDB(dbName, ",", quote)...)
		lines = append(lines, "")
	}

	tblNames := lo.Keys(engine.tables)
	sort.Strings(tblNames)
	for _, tblName := range tblNames {
		lines = append(lines, "Table: "+tblName)
		lines = append(lines, engine.prettyTable(tblName, ",", quote)...)
		lines = append(lines, "")
	}
	return lines
}

func (engine *Engine) prettyDB(tblName, sep string, quote bool) []string {
	table, found := engine.dbs[tblName]
	if !found {
		return nil
	}
	list := _pretty(table, sep, quote)
	return list
}

func (engine *Engine) prettyTable(tblName, sep string, quote bool) []string {
	table, found := engine.tables[tblName]
	if !found {
		return nil
	}
	list := _pretty(table, sep, quote)
	return list
}

func setupEngine(backupDir, appId string) (*Engine, *EtlApp) {
	serverLogger := log.New(os.Stdout, "[test]: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile|log.Lmsgprefix)
	engine := &Engine{}
	app := &EtlApp{
		logger:    serverLogger,
		engine:    engine,
		appId:     appId,
		backupDir: backupDir,
	}
	app.engine.Init(app)
	return engine, app
}
