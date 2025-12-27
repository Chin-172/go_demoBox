package db

import (
	"database/sql"
	"fmt"
	"slices"
)

type DataAuth struct {
	ID    int
	CaseA int
	CaseB int
	CaseC int
}

func CheckAuth(groupID int, caseLevel []string) bool {
	table = "data_auth"
	var auth DataAuth
	var group Group

	sql_statement := "SELECT * FROM group_list WHERE id = $1;"
	err := dbPool.QueryRow(sql_statement, groupID).Scan(&group.id, &group.group, &group.group_leader, &group.desc, &group.auth)
	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}
	fmt.Printf("group.auth: %d\n", group.auth)

	sql_statement = "SELECT * FROM " + table + " WHERE id = $1;"
	err = dbPool.QueryRow(sql_statement, group.auth).Scan(&auth.ID, &auth.CaseA, &auth.CaseB, &auth.CaseC)
	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}

	fmt.Printf("CaseA: %d\n", auth.CaseA)
	fmt.Printf("slices.Contains: %s\n", caseLevel)

	if slices.Contains(caseLevel, "A") {
		return auth.CaseA == 1
	} else if slices.Contains(caseLevel, "B") {
		return auth.CaseB == 1
	} else if slices.Contains(caseLevel, "C") {
		return auth.CaseC == 1
	}
	return false
}
