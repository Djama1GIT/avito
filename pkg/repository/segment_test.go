package repository_test

import (
	"avito/pkg/repository"
	"avito/pkg/structures"
	"errors"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSegment_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewSegmentDB(db)

	type args struct {
		Slug string
	}

	type mockBehavior func(args args, slug string)

	tests := []struct {
		name          string
		mockBehavior  mockBehavior
		args          args
		wantErr       bool
		expectedError string
	}{
		{
			name: "OK",
			mockBehavior: func(args args, slug string) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"slug"}).AddRow(slug)
				mock.ExpectQuery("INSERT INTO segments").
					WithArgs(args.Slug).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			args: args{
				Slug: "example",
			},
			wantErr: false,
		},
		{
			name: "DuplicateSlug",
			mockBehavior: func(args args, slug string) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO segments").
					WithArgs(args.Slug).
					WillReturnError(errors.New("duplicate slug"))

				mock.ExpectRollback()
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "duplicate slug",
		},
		{
			name: "TransactionError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "transaction error",
		},
		{
			name: "CommitError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"slug"}).AddRow(slug)
				mock.ExpectQuery("INSERT INTO segments").
					WithArgs(args.Slug).
					WillReturnRows(rows)

				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "commit error",
		},
		{
			name: "QueryError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectBegin()

				mock.ExpectQuery("INSERT INTO segments").
					WithArgs(args.Slug).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "query error",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.Slug)

			got, err := repo.Create(structures.Segment(testCase.args))
			if testCase.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.Slug, got)
			}
		})
	}
}

func TestSegment_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewSegmentDB(db)

	type args struct {
		Slug string
	}

	type mockBehavior func(args args, slug string)

	tests := []struct {
		name          string
		mockBehavior  mockBehavior
		args          args
		wantErr       bool
		expectedError string
	}{
		{
			name: "DeleteExistingSegment",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(args.Slug).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectBegin()
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				mock.ExpectCommit()
				mock.ExpectQuery("INSERT").
					WithArgs(1, args.Slug, false).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				mock.ExpectCommit()
			},
			args: args{
				Slug: "example",
			},
			wantErr: false,
		},
		{
			name: "DeleteNonExistingSegment",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(0))
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "segment with slug example does not exist",
		},
		{
			name: "TransactionError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "transaction error",
		},
		{
			name: "SelectError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnError(errors.New("error reading segment"))
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "error reading segment",
		},
		{
			name: "DeleteError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin()

				mock.ExpectExec("DELETE").
					WithArgs(args.Slug).
					WillReturnError(errors.New("delete error"))

				mock.ExpectRollback()
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "delete error",
		},
		{
			name: "CreateRepoForHistoryError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin()

				mock.ExpectExec("DELETE").
					WithArgs(args.Slug).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectBegin().WillReturnError(errors.New("repo for history error"))
				mock.ExpectRollback()
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "repo for history error",
		},
		{
			name: "HistoryError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin()
				mock.ExpectExec("DELETE").
					WithArgs(args.Slug).
					WillReturnResult(sqlmock.NewResult(0, 1))

				mock.ExpectBegin()
				mock.ExpectQuery("SELECT").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))
				mock.ExpectCommit()
				mock.ExpectQuery("INSERT").
					WithArgs(1, args.Slug, false).
					WillReturnError(errors.New("history error"))
				mock.ExpectRollback()
			},
			args: args{
				Slug: "example",
			},
			wantErr:       true,
			expectedError: "history error",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.Slug)

			got, err := repo.Delete(structures.Segment(testCase.args))
			if testCase.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, testCase.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.Slug, got)
			}
		})
	}
}
