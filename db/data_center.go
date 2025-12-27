package db

import (
	"database/sql"
	"time"
)

type DataCenter struct {
	ID        int
	FileName  string
	Year      int
	Signatory string
	FileType  string
	AddDt     time.Time
	Operator  string
}

func GetAllDataList(caseLevel []string) []DataCenter {
	var res []DataCenter
	table = "data_center"
	var rows *sql.Rows
	var err error
	var sql_statement string
	if len(caseLevel) == 1 {
		sql_statement = "Select * From " + table + " WHERE file_type = $1;"
		rows, err = dbPool.Query(sql_statement, "Case "+caseLevel[0])
		if err != nil && err != sql.ErrNoRows {
			CheckError(err)
		}
	} else {
		sql_statement = "Select * From " + table + " WHERE "
		for i, level := range caseLevel {
			sql_statement += "file_type = " + level
			if i+1 != len(caseLevel) {
				sql_statement += " AND "
			}
		}
		rows, err = dbPool.Query(sql_statement)
		if err != nil && err != sql.ErrNoRows {
			CheckError(err)
		}
	}

	for rows.Next() {
		var data DataCenter
		err = rows.Scan(&data.ID, &data.FileName, &data.Year, &data.Signatory, &data.FileType, &data.AddDt, &data.Operator)
		CheckError(err)

		res = append(res, data)
	}

	return res
}
