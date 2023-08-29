package repository_test

import (
	"avito/pkg/repository"
	"avito/pkg/structures"
	"database/sql"
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
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
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
					WillReturnError(sql.ErrNoRows)

				mock.ExpectRollback()
			},
			args: args{
				Slug: "example",
			},
			wantErr: true,
		},
		{
			name: "TransactionError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			args: args{
				Slug: "example",
			},
			wantErr: true,
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
			wantErr: true,
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
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.Slug)

			got, err := repo.Create(structures.Segment(testCase.args))
			if testCase.wantErr {
				assert.Error(t, err)
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
		name         string
		mockBehavior mockBehavior
		args         args
		wantErr      bool
	}{
		{
			name: "DeleteExistingSegment",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT (.+) FROM segments").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM segments").
					WithArgs(args.Slug).
					WillReturnResult(sqlmock.NewResult(0, 1))
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
				mock.ExpectQuery("SELECT (.+) FROM segments").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(0))
			},
			args: args{
				Slug: "example",
			},
			wantErr: true,
		},
		{
			name: "TransactionError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT (.+) FROM segments").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin().WillReturnError(errors.New("transaction error"))
			},
			args: args{
				Slug: "example",
			},
			wantErr: true,
		},
		{
			name: "SelectError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT (.+) FROM segments").
					WithArgs(args.Slug).
					WillReturnError(errors.New("error reading"))
			},
			args: args{
				Slug: "example",
			},
			wantErr: true,
		},
		{
			name: "DeleteError",
			mockBehavior: func(args args, slug string) {
				mock.ExpectQuery("SELECT (.+) FROM segments").
					WithArgs(args.Slug).
					WillReturnRows(sqlmock.NewRows([]string{"slug"}).AddRow(1))

				mock.ExpectBegin()

				mock.ExpectExec("DELETE FROM segments").
					WithArgs(args.Slug).
					WillReturnError(errors.New("delete error"))

				mock.ExpectRollback()
			},
			args: args{
				Slug: "example",
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.Slug)

			got, err := repo.Delete(structures.Segment(testCase.args))
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.Slug, got)
			}
		})
	}
}
