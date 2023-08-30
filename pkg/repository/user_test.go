package repository_test

import (
	repository "avito/pkg/repository"
	"avito/pkg/structures"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUser_GetUserHistory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUserDB(db)

	type args struct {
		UserHistory structures.UserHistory
	}

	type mockBehavior func(args args, userHistory structures.UserHistory)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantReport   string
		wantErr      bool
		expectError  string
	}{
		{
			name: "Success",
			mockBehavior: func(args args, userHistory structures.UserHistory) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"user_id", "segment", "operation", "operation_datetime"}).
					AddRow(1, "segment1", false, time.Now())

				mock.ExpectQuery("SELECT").
					WithArgs(1, "2023-08").
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				UserHistory: structures.UserHistory{
					Id:        1,
					YearMonth: "2023-08",
				},
			},
			wantReport:  "reports/user_history_2023-08_1.csv",
			wantErr:     false,
			expectError: "",
		},
		{
			name: "QueryError",
			mockBehavior: func(args args, userHistory structures.UserHistory) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT").
					WithArgs(1, "2023-08").
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			args: args{
				UserHistory: structures.UserHistory{
					Id:        1,
					YearMonth: "2023-08",
				},
			},
			wantReport:  "",
			wantErr:     true,
			expectError: "query error",
		},
		{
			name: "BeginError",
			mockBehavior: func(args args, userHistory structures.UserHistory) {
				mock.ExpectBegin().WillReturnError(errors.New("db error"))
			},
			args: args{
				UserHistory: structures.UserHistory{
					Id:        1,
					YearMonth: "2023-08",
				},
			},
			wantReport:  "",
			wantErr:     true,
			expectError: "db error",
		},
		{
			name: "RowsError",
			mockBehavior: func(args args, userHistory structures.UserHistory) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"user"}).
					AddRow(nil).RowError(0, errors.New("Row error"))

				mock.ExpectQuery("SELECT").
					WithArgs(1, "2023-08").
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args: args{
				UserHistory: structures.UserHistory{
					Id:        1,
					YearMonth: "2023-08",
				},
			},
			wantReport:  "",
			wantErr:     true,
			expectError: "Row error",
		},
		{
			name: "RowsError",
			mockBehavior: func(args args, userHistory structures.UserHistory) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"user_id", "segment", "operation", "operation_datetime"}).
					AddRow(nil, nil, nil, nil)

				mock.ExpectQuery("SELECT").
					WithArgs(1, "2023-08").
					WillReturnRows(rows)
				mock.ExpectRollback()
			},
			args: args{
				UserHistory: structures.UserHistory{
					Id:        1,
					YearMonth: "2023-08",
				},
			},
			wantReport:  "",
			wantErr:     true,
			expectError: "sql: Scan error on column index 0, name \"user_id\": converting NULL to int is unsupported",
		},
		{
			name: "CommitError",
			mockBehavior: func(args args, userHistory structures.UserHistory) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"user_id", "segment", "operation", "operation_datetime"}).
					AddRow(1, "segment1", false, time.Now())

				mock.ExpectQuery("SELECT").
					WithArgs(1, "2023-08").
					WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("Commit error"))
			},
			args: args{
				UserHistory: structures.UserHistory{
					Id:        1,
					YearMonth: "2023-08",
				},
			},
			wantReport:  "",
			wantErr:     true,
			expectError: "Commit error",
		},
		{
			name: "InvalidYearMonth",
			mockBehavior: func(args args, userHistory structures.UserHistory) {

			},
			args: args{
				UserHistory: structures.UserHistory{
					Id:        1,
					YearMonth: "2023/8",
				},
			},
			wantReport:  "",
			wantErr:     true,
			expectError: "invalid YearMonth",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.UserHistory)

			gotReport, err := repo.GetUserHistory(testCase.args.UserHistory)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantReport, gotReport)
			}
		})
	}
}

func TestUser_DeleteExpiredSegments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUserDB(db)

	type mockBehavior func()

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		wantErr      bool
		expectError  string
	}{
		{
			name: "Success",
			mockBehavior: func() {
				mock.ExpectBegin()

				mock.ExpectExec("DELETE").
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectCommit()
			},
			wantErr:     false,
			expectError: "",
		},
		{
			name: "BeginError",
			mockBehavior: func() {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			wantErr:     true,
			expectError: "begin error",
		},
		{
			name: "ExecError",
			mockBehavior: func() {
				mock.ExpectBegin()

				mock.ExpectExec("DELETE").
					WillReturnError(errors.New("exec error"))

				mock.ExpectRollback()
			},
			wantErr:     true,
			expectError: "exec error",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior()

			err := repo.DeleteExpiredSegments()
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
