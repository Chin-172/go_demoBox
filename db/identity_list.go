package db

import (
	"database/sql"
	"fmt"
)

type Identity struct {
	identity    int
	description string
	data_center int
	log_history int
	group_list  int
	user_list   int
}

func GetUserActionList(identity int) []string {
	table = "identity_list"
	var res Identity
	var actionList []string
	sql_statment := "Select * From " + table + " WHERE identity = $1;"
	err := dbPool.QueryRow(sql_statment, identity).Scan(&res.identity, &res.description, &res.data_center, &res.log_history, &res.group_list, &res.user_list)
	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}
	fmt.Printf("GetUserActionList res.identity:%d\n", res.identity)
	if res.data_center == 1 {
		actionList = append(actionList, "Data Center")
	}
	if res.log_history == 1 {
		actionList = append(actionList, "Log")
	}
	if res.group_list == 1 {
		actionList = append(actionList, "Group List")
	}
	if res.user_list == 1 {
		actionList = append(actionList, "User List")
	}

	for _, action := range actionList {
		fmt.Println(action)
	}
	return actionList
}
