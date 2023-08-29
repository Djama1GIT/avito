package repository_test

import (
	"avito/pkg/repository"
	"avito/pkg/structures"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_NewRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := repository.NewRepository(db)

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
			name: "RepositoryTestOK",
			mockBehavior: func(args args, slug string) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"slug"}).AddRow(args.Slug)
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
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.args.Slug)

			got, err := repo.Segment.Create(structures.Segment(testCase.args))
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.args.Slug, got)
			}
		})
	}
}
