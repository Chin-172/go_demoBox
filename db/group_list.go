package db

import "database/sql"

type Group struct {
	id           int
	group        string
	group_leader string
	desc         string
	auth         int
}

func GetUserGroup(username string) []Group {
	table = "group_list"
	var res []Group
	sql_statement := "SELECT * FROM " + table + " WHERE username = ?;"
	rows, err := dbPool.Query(sql_statement, username)
	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}

	for rows.Next() {
		var group Group
		err = rows.Scan(&group)
		CheckError(err)
		res = append(res, group)
	}
	return res
}
