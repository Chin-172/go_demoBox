package db

import "database/sql"

var (
	table_name = "user_list"
)

type User struct {
	ID       int
	Username string
	Password string
	Identity int
}

func GetAllUser() []User {

	sql_statement := "SELECT * FROM " + table_name + ";"
	var res []User
	rows, err := dbPool.Query(sql_statement)

	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.Identity)
		CheckError(err)

		res = append(res, user)
	}
	return res
}

func GetUserInfo(username string) User {
	var user User
	sql_statement := "SELECT * FROM " + table_name + " WHERE username = $1;"
	err := dbPool.QueryRow(sql_statement, username).Scan(&user.ID, &user.Username, &user.Password, &user.Identity)
	if err != nil && err != sql.ErrNoRows {
		CheckError(err)
	}
	return user
}
