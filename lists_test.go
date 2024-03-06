package main

import (
	"strings"
	"testing"
)

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
