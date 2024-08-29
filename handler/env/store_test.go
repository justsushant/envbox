package env

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/justsushant/envbox/types"
)


func TestSaveContainer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error while opening stub database connection: %v", err)
	}
	defer db.Close()
	stubStore := NewStore(db)

	// save container test
	mock.ExpectExec("INSERT INTO containers_running").
		WithArgs("testContainerID", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err = stubStore.SaveContainer("testContainerID", 1); err != nil {
		t.Fatalf("Error while saving details in DB: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}

func TestDeleteContainer(t *testing.T) {
	tt := []struct {
		name     string
		id       string
		mockServiceOutput driver.Result
		expError error
	}{
		{
			name:     "delete container happy path",
			id:       "testID",
			mockServiceOutput: sqlmock.NewResult(1, 1),
			expError: nil,
		},
		{
			name:     "delete container unhappy path of error",
			id:       "testID",
			mockServiceOutput: nil,
			expError: fmt.Errorf("test error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// setting the mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error while opening stub database connection: %v", err)
			}
			defer db.Close()

			// setting the store
			stubStore := NewStore(db)

			// setting the expected value in the mock
			expQuery := mock.ExpectExec("UPDATE containers_running SET active = 0 WHERE id = ?")
			if tc.expError != nil {
				expQuery.
				WithArgs(tc.id).
				WillReturnError(tc.expError)
			} else {
				expQuery.
				WithArgs(tc.id).
				WillReturnResult(tc.mockServiceOutput)
			}
			
			// calling the function
			err = stubStore.DeleteContainer(tc.id)

			// checking the results
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if err != tc.expError {
					t.Errorf("Expected error %v, but got %v", tc.expError, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpedted error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetAllEnvs(t *testing.T) {
	tt := []struct{
		name string
		mockServiceOutput *sqlmock.Rows
		expEnv []types.Env
		expErr error
	}{
		{
			name: "get all envs happy path",
			mockServiceOutput: sqlmock.NewRows([]string{
				"id", "imageName", "containerID", "accessLink", "active", "createdAt",
			}).AddRow("testID", "testImageName", "testContainerID", "testAccessLink", 1, "2024-08-27 00:00:00"),
			expEnv: []types.Env{
				{
					ID: "testID",
					ImageName: "testImageName",
					ContainerID: "testContainerID",
					AccessLink: "testAccessLink",
					Active: true,
					CreatedAt: "2024-08-27 00:00:00",
				},
			},
			expErr: nil,
		},
		{
			name: "get all envs unhappy path of zero envs",
			mockServiceOutput: sqlmock.NewRows([]string{
				"id", "imageName", "containerID", "accessLink", "active", "createdAt",
			}),
			expEnv: []types.Env{},
			expErr: sql.ErrNoRows,
		},
		{
			name: "get all envs unhappy path of error",
			expEnv: nil,
			expErr: fmt.Errorf("test error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// setting the mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error while opening stub database connection: %v", err)
			}
			defer db.Close()

			// setting the store
			stubStore := NewStore(db)

			// setting the expected value in the mock
			expQuery := mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM containers_running WHERE active = 1"))
			if tc.expErr != nil {
				expQuery.
				WithArgs().
				WillReturnError(tc.expErr)
			} else {
				expQuery.
				WithArgs().
				WillReturnRows(tc.mockServiceOutput)
			}

			// calling the function
			envs, err := stubStore.GetAllEnvs()

			// checking the results
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if err != tc.expErr {
					t.Errorf("Expected error %v, but got %v", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(tc.expEnv) != len(envs) {
				t.Fatalf("Expected %d envs, but got %d", len(tc.expEnv), len(envs))
			}

			// if both slices have size of zero, then they're equal
			if len(tc.expEnv) == 0 && len(envs) == 0 {
				return
			}

			if !reflect.DeepEqual(tc.expEnv, envs) {
				t.Fatalf("Expected %v, but got %v", tc.expEnv, envs)
			}


			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestGetContainerByID(t *testing.T) {
	tt := []struct{
		name string
		id string
		mockServiceOutput *sqlmock.Rows
		expEnv types.Env
		expErr error
	}{
		{
			name: "get container by id happy path",
			id: "testID",
			mockServiceOutput: sqlmock.NewRows([]string{
				"id", "imageName", "containerID", "accessLink", "active", "createdAt",
			}).AddRow("testID", "testImageName", "testContainerID", "testAccessLink", 1, "2024-08-27 00:00:00"),
			expEnv: types.Env{
				ID: "testID",
				ImageName: "testImageName",
				ContainerID: "testContainerID",
				AccessLink: "testAccessLink",
				Active: true,
				CreatedAt: "2024-08-27 00:00:00",
			},
			expErr: nil,
		},
		{
			name: "get container by id unhappy path of no env",
			id: "testID",
			mockServiceOutput: sqlmock.NewRows([]string{
				"id", "imageName", "containerID", "accessLink", "active", "createdAt",
			}),
			expEnv: types.Env{},
			expErr: sql.ErrNoRows,
		},
		{
			name: "get container by id unhappy path of error",
			id: "testID",
			expErr: fmt.Errorf("test error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// setting the mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error while opening stub database connection: %v", err)
			}
			defer db.Close()

			// setting the store
			stubStore := NewStore(db)

			// setting the expected value in the mock
			expQuery := mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM containers_running WHERE id = ?"))
			if tc.expErr != nil {
				expQuery.
				WithArgs().
				WillReturnError(tc.expErr)
			} else {
				expQuery.
				WithArgs().
				WillReturnRows(tc.mockServiceOutput)
			}

			// calling the function
			env, err := stubStore.GetContainerByID(tc.id)

			// checking the results
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if err != tc.expErr {
					t.Errorf("Expected error %v, but got %v", tc.expErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(tc.expEnv, env) {
				t.Fatalf("Expected %v, but got %v", tc.expEnv, env)
			}


			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestUpdateContainerAccessLink(t *testing.T) {
	tt := []struct {
		name     string
		containerID       string
		accessLink string
		mockServiceOutput driver.Result
		expError error
	}{
		{
			name:     "update container link happy path",
			containerID:  "testID",
			accessLink: "testLink",
			mockServiceOutput: sqlmock.NewResult(1, 1),
			expError: nil,
		},
		{
			name:     "update container link unhappy path of error",
			containerID:  "testID",
			accessLink: "testLink",
			expError: fmt.Errorf("test error"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// setting the mock
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Error while opening stub database connection: %v", err)
			}
			defer db.Close()

			// setting the store
			stubStore := NewStore(db)

			// setting the expected value in the mock
			expQuery := mock.ExpectExec(regexp.QuoteMeta("UPDATE containers_running SET accessLink = ? WHERE containerID = ?"))
			if tc.expError != nil {
				expQuery.
				WithArgs(tc.accessLink, tc.containerID).
				WillReturnError(tc.expError)
			} else {
				expQuery.
				WithArgs(tc.accessLink, tc.containerID).
				WillReturnResult(tc.mockServiceOutput)
			}
			
			// calling the function
			err = stubStore.UpdateContainerAccessLink(tc.containerID, tc.accessLink)

			// checking the results
			if tc.expError != nil {
				if err == nil {
					t.Fatalf("Expected error but got nil")
				}

				if err != tc.expError {
					t.Errorf("Expected error %v, but got %v", tc.expError, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpedted error: %v", err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unfulfilled expectations: %v", err)
			}
		})
	}
}