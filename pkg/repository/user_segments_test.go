package repository_test

import (
	repository "avito/pkg/repository"
	"avito/pkg/structures"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserSegments_Patch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUserSegmentsDB(db)

	type args struct {
		UserSegments structures.UserSegments
	}

	type mockBehavior func(args args, userSegments structures.UserSegments)

	var validDateTime = "2023-08-28 20:00:00"
	var invalidDateTime = "2023-08-28 202:00:00"

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantUserID   int
		wantErr      bool
		expectError  string
	}{
		{
			name: "AddSegments_Success",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToAdd {
					mock.ExpectExec("INSERT").
						WithArgs(userSegments.UserId, segment, &validDateTime).
						WillReturnResult(sqlmock.NewResult(0, 1))
					mock.ExpectQuery("INSERT").
						WithArgs(1, segment, true).
						WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				}

				mock.ExpectCommit()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:                  1,
					SegmentsToAdd:           []string{"segment1", "segment2"},
					SegmentsToAddExpiration: &validDateTime,
					SegmentsToDelete:        nil,
				},
			},
			wantUserID:  1,
			wantErr:     false,
			expectError: "",
		},
		{
			name: "AddSegments_Success",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToAdd {
					mock.ExpectExec("INSERT").
						WithArgs(userSegments.UserId, segment).
						WillReturnResult(sqlmock.NewResult(0, 1))
					mock.ExpectQuery("INSERT").
						WithArgs(1, segment, true).
						WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				}

				mock.ExpectCommit()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    []string{"segment1", "segment2"},
					SegmentsToDelete: nil,
				},
			},
			wantUserID:  1,
			wantErr:     false,
			expectError: "",
		},
		{
			name: "AddSegments_ErrorWithTime",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:                  1,
					SegmentsToAdd:           []string{"segment1", "segment2"},
					SegmentsToAddExpiration: &invalidDateTime,
					SegmentsToDelete:        nil,
				},
			},
			wantUserID:  1,
			wantErr:     true,
			expectError: "parsing time \"2023-08-28 202:00:00\" as \"2006-01-02 15:04:05\": cannot parse \"2:00:00\" as \":\"",
		},
		{
			name: "AddSegments_ErrorWithValidTime",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToAdd {
					mock.ExpectExec("INSERT").
						WithArgs(userSegments.UserId, segment, &validDateTime).
						WillReturnError(errors.New("uwaaa"))
					break
				}

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:                  1,
					SegmentsToAdd:           []string{"segment1", "segment2"},
					SegmentsToAddExpiration: &validDateTime,
					SegmentsToDelete:        nil,
				},
			},
			wantUserID:  1,
			wantErr:     true,
			expectError: "error occurred while processing segment to add 'segment1': uwaaa",
		},
		{
			name: "DeleteSegments_Success",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToDelete {
					mock.ExpectQuery("SELECT").
						WithArgs(userSegments.UserId, segment).
						WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

					mock.ExpectExec("DELETE").
						WithArgs(userSegments.UserId, segment).
						WillReturnResult(sqlmock.NewResult(0, 1))
					mock.ExpectQuery("INSERT").
						WithArgs(1, segment, false).
						WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				}

				mock.ExpectCommit()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    nil,
					SegmentsToDelete: []string{"segment1", "segment2"},
				},
			},
			wantUserID:  1,
			wantErr:     false,
			expectError: "",
		},
		{
			name: "TransactionError",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    []string{"segment1", "segment2"},
					SegmentsToDelete: nil,
				},
			},
			wantUserID:  -1,
			wantErr:     true,
			expectError: "transaction error",
		},
		{
			name: "AddSegments_Error",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToAdd {
					mock.ExpectExec("INSERT").
						WithArgs(userSegments.UserId, segment).
						WillReturnError(errors.New("insert error"))
					break
				}

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    []string{"segment1", "segment2"},
					SegmentsToDelete: nil,
				},
			},
			wantUserID:  -1,
			wantErr:     true,
			expectError: "error occurred while processing segment to add 'segment1': insert error",
		},
		{
			name: "DeleteSegments_Error",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToDelete {
					mock.ExpectQuery("SELECT").
						WithArgs(userSegments.UserId, segment).
						WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

					mock.ExpectExec("DELETE").
						WithArgs(userSegments.UserId, segment).
						WillReturnError(errors.New("delete error"))
					break
				}

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    nil,
					SegmentsToDelete: []string{"segment1", "segment2"},
				},
			},
			wantUserID:  -1,
			wantErr:     true,
			expectError: "error occurred while processing segment to delete 'segment1': delete error",
		},
		{
			name: "DeleteSegments_ErrorNonExistance",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToDelete {
					mock.ExpectQuery("SELECT").
						WithArgs(userSegments.UserId, segment).
						WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(0))
					break
				}

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    nil,
					SegmentsToDelete: []string{"segment1", "segment2"},
				},
			},
			wantUserID:  -1,
			wantErr:     true,
			expectError: "error occurred while checking segment to delete existence 'segment1': user(1) is not in this segment",
		},
		{
			name: "DeleteSegments_ErrorNonExistanceError",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToDelete {
					mock.ExpectQuery("SELECT").
						WithArgs(userSegments.UserId, segment).
						WillReturnError(errors.New("existance select error"))
					break
				}

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    nil,
					SegmentsToDelete: []string{"segment1", "segment2"},
				},
			},
			wantUserID:  -1,
			wantErr:     true,
			expectError: "error occurred while checking segment to delete existence 'segment1': existance select error",
		},
		{
			name: "AddSegments_HistoryError",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToAdd {
					mock.ExpectExec("INSERT").
						WithArgs(userSegments.UserId, segment).
						WillReturnResult(sqlmock.NewResult(0, 1))
					mock.ExpectQuery("INSERT").
						WithArgs(1, segment, true).
						WillReturnError(errors.New("history error"))
					break
				}

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    []string{"segment1", "segment2"},
					SegmentsToDelete: nil,
				},
			},
			wantUserID:  -1,
			wantErr:     true,
			expectError: "history error",
		},
		{
			name: "DeleteSegments_HistoryError",
			mockBehavior: func(args args, userSegments structures.UserSegments) {
				mock.ExpectBegin()

				for _, segment := range userSegments.SegmentsToDelete {
					mock.ExpectQuery("SELECT").
						WithArgs(userSegments.UserId, segment).
						WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

					mock.ExpectExec("DELETE").
						WithArgs(userSegments.UserId, segment).
						WillReturnResult(sqlmock.NewResult(0, 1))
					mock.ExpectQuery("INSERT").
						WithArgs(1, segment, false).
						WillReturnError(errors.New("history error"))
				}

				mock.ExpectRollback()
			},
			args: args{
				UserSegments: structures.UserSegments{
					UserId:           1,
					SegmentsToAdd:    nil,
					SegmentsToDelete: []string{"segment1", "segment2"},
				},
			},
			wantUserID:  -1,
			wantErr:     true,
			expectError: "history error",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.UserSegments)

			gotUserID, err := repo.Patch(testCase.args.UserSegments)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.wantUserID, gotUserID)
			}
		})
	}
}

func TestUserSegments_GetUserSegments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUserSegmentsDB(db)

	type args struct {
		User structures.User
	}

	type mockBehavior func(args args, user structures.User)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantSlugs    []string
		wantErr      bool
		expectError  string
	}{
		{
			name: "Success",
			mockBehavior: func(args args, user structures.User) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"}).
					AddRow("segment1").
					AddRow("segment2")

				mock.ExpectQuery("SELECT segment").
					WithArgs(user.Id).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				User: structures.User{
					Id: 1,
				},
			},
			wantSlugs:   []string{"segment1", "segment2"},
			wantErr:     false,
			expectError: "",
		},
		{
			name: "NoRows",
			mockBehavior: func(args args, user structures.User) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"})

				mock.ExpectQuery("SELECT segment").
					WithArgs(user.Id).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				User: structures.User{
					Id: 1,
				},
			},
			wantSlugs:   nil,
			wantErr:     false,
			expectError: "",
		},
		{
			name: "TransactionError",
			mockBehavior: func(args args, user structures.User) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			args: args{
				User: structures.User{
					Id: 1,
				},
			},
			wantSlugs:   nil,
			wantErr:     true,
			expectError: "transaction error",
		},
		{
			name: "QueryError",
			mockBehavior: func(args args, user structures.User) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT segment").
					WithArgs(user.Id).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			args: args{
				User: structures.User{
					Id: 1,
				},
			},
			wantSlugs:   nil,
			wantErr:     true,
			expectError: "query error",
		},
		{
			name: "RowsScanError",
			mockBehavior: func(args args, user structures.User) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"}).
					AddRow(nil)

				mock.ExpectQuery("SELECT segment").
					WithArgs(user.Id).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args: args{
				User: structures.User{
					Id: 1,
				},
			},
			wantSlugs:   nil,
			wantErr:     true,
			expectError: "sql: Scan error on column index 0, name \"segment\": converting NULL to string is unsupported",
		},
		{
			name: "RowsError",
			mockBehavior: func(args args, user structures.User) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"}).
					AddRow(nil).RowError(0, errors.New("Row error"))

				mock.ExpectQuery("SELECT segment").
					WithArgs(user.Id).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args: args{
				User: structures.User{
					Id: 1,
				},
			},
			wantSlugs:   nil,
			wantErr:     true,
			expectError: "Row error",
		},
		{
			name: "CommitError",
			mockBehavior: func(args args, user structures.User) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"}).
					AddRow("segment1").
					AddRow("segment2")

				mock.ExpectQuery("SELECT segment").
					WithArgs(user.Id).
					WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("Commit error"))
			},
			args: args{
				User: structures.User{
					Id: 1,
				},
			},
			wantSlugs:   nil,
			wantErr:     true,
			expectError: "Commit error",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.User)

			gotSlugs, err := repo.GetUserSegments(testCase.args.User)

			assert.Equal(t, testCase.wantSlugs, gotSlugs)

			if testCase.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, testCase.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserSegments_GetSegmentUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewUserSegmentsDB(db)

	type args struct {
		Segment structures.Segment
	}

	type mockBehavior func(args args, segment structures.Segment)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		args         args
		wantUsers    []int
		wantErr      bool
		expectError  string
	}{
		{
			name: "Success",
			mockBehavior: func(args args, segment structures.Segment) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"user"}).
					AddRow(1).
					AddRow(2)

				mock.ExpectQuery("SELECT user_id").
					WithArgs(segment.Slug).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				Segment: structures.Segment{
					Slug: "example",
				},
			},
			wantUsers:   []int{1, 2},
			wantErr:     false,
			expectError: "",
		},
		{
			name: "NoRows",
			mockBehavior: func(args args, segment structures.Segment) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"})

				mock.ExpectQuery("SELECT user_id").
					WithArgs(segment.Slug).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				Segment: structures.Segment{
					Slug: "example",
				},
			},
			wantUsers:   nil,
			wantErr:     false,
			expectError: "",
		},
		{
			name: "TransactionError",
			mockBehavior: func(args args, segment structures.Segment) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			args: args{
				Segment: structures.Segment{
					Slug: "example",
				},
			},
			wantUsers:   nil,
			wantErr:     true,
			expectError: "transaction error",
		},
		{
			name: "QueryError",
			mockBehavior: func(args args, segment structures.Segment) {
				mock.ExpectBegin()

				mock.ExpectQuery("SELECT user_id").
					WithArgs(segment.Slug).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			args: args{
				Segment: structures.Segment{
					Slug: "example",
				},
			},
			wantUsers:   nil,
			wantErr:     true,
			expectError: "query error",
		},
		{
			name: "RowsScanError",
			mockBehavior: func(args args, segment structures.Segment) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"}).
					AddRow(nil)

				mock.ExpectQuery("SELECT user_id").
					WithArgs(segment.Slug).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args: args{
				Segment: structures.Segment{
					Slug: "example",
				},
			},
			wantUsers:   nil,
			wantErr:     true,
			expectError: "sql: Scan error on column index 0, name \"segment\": converting NULL to int is unsupported",
		},
		{
			name: "RowsError",
			mockBehavior: func(args args, segment structures.Segment) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"}).
					AddRow(nil).RowError(0, errors.New("Row error"))

				mock.ExpectQuery("SELECT user_id").
					WithArgs(segment.Slug).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			args: args{
				Segment: structures.Segment{
					Slug: "example",
				},
			},
			wantUsers:   nil,
			wantErr:     true,
			expectError: "Row error",
		},
		{
			name: "CommitError",
			mockBehavior: func(args args, segment structures.Segment) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"segment"}).
					AddRow(1).
					AddRow(2)

				mock.ExpectQuery("SELECT user_id").
					WithArgs(segment.Slug).
					WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("Commit error"))
			},
			args: args{
				Segment: structures.Segment{
					Slug: "example",
				},
			},
			wantUsers:   nil,
			wantErr:     true,
			expectError: "Commit error",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.Segment)

			gotSlugs, err := repo.GetSegmentUsers(testCase.args.Segment)

			assert.Equal(t, testCase.wantUsers, gotSlugs)

			if testCase.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, testCase.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
