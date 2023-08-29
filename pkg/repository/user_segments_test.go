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
						WithArgs(userSegments.UserId, segment).
						WillReturnResult(sqlmock.NewResult(0, 1))
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

func TestUserSegments_GetUsersInSegment(t *testing.T) {
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
			name: "GetUsersInSegment_Success",
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
			name: "GetUsersInSegment_NoRows",
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

			gotSlugs, err := repo.GetUsersInSegment(testCase.args.User)

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
