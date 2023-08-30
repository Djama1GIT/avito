package repository

import (
	"avito/pkg/structures"
	"avito/pkg/utils"
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type UserDB struct {
	db *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{db: db}
}

func (r *UserDB) GetUserHistory(userHistory structures.UserHistory) (string, error) {
	if err := utils.ValidateYearMonth(userHistory.YearMonth); err != nil {
		return "", err
	}

	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}

	reportFolderPath := "reports"
	err = os.MkdirAll(reportFolderPath, os.ModePerm)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	reportFileName := fmt.Sprintf("%s/user_history_%s_%d.csv", reportFolderPath, userHistory.YearMonth, userHistory.Id)
	reportFile, err := os.Create(reportFileName)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	defer reportFile.Close()

	writer := csv.NewWriter(reportFile)
	defer writer.Flush()

	createSegmentQuery := fmt.Sprintf("SELECT user_id, segment, operation, operation_datetime FROM %s WHERE user_id = $1 AND to_char(operation_datetime, 'YYYY-MM') = $2", userSegmentsHistoryTable)
	rows, err := tx.Query(createSegmentQuery, userHistory.Id, userHistory.YearMonth)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		var segment string
		var operation bool
		var operationDatetime time.Time

		err := rows.Scan(&userID, &segment, &operation, &operationDatetime)
		if err != nil {
			tx.Rollback()
			return "", err
		}

		operationStr := "добавление"
		if !operation {
			operationStr = "удаление"
		}

		record := []string{
			fmt.Sprintf("%d", userID),
			segment,
			operationStr,
			operationDatetime.Format("2006-01-02 15:04:05"),
		}

		err = writer.Write(record)
		if err != nil {
			tx.Rollback()
			return "", err
		}
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return reportFileName, nil
}

func historyUpdate(tx *sql.Tx, segment string, userId int, operation bool) (int, error) {
	// operation:
	// 		true - insert
	//		false - delete
	var user_id int
	createUserSegmentsHistoryQuery := fmt.Sprintf(
		"INSERT INTO %s (user_id, segment, operation) VALUES ($1, $2, $3) RETURNING user_id",
		userSegmentsHistoryTable)

	row := tx.QueryRow(createUserSegmentsHistoryQuery, userId, segment, operation)
	if err := row.Scan(&user_id); err != nil {
		tx.Rollback()
		return -1, err
	}
	return user_id, nil
}
