package models

import (
	"testing"
	"io/ioutil"
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatePicAPIFailGetMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/pic/create", nil)
	if err != nil {
		t.Errorf("Failed to create new request: %v", err)
	}

	recorder := httptest.NewRecorder()

	a := App{}

	a.router = mux.NewRouter()
	a.router.HandleFunc("/api/pic/create", a.CreatePicAPI).Methods("POST")
	a.router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Errorf("HTTP code expected 405 (Method Not Allowed), but got instead: %d", recorder.Code)
	}
}

func TestCreatePicAPI(t *testing.T) {
	expectedPath := "test.jpg"
	expectedDescription := "this is test"

	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		// TODO improve mocks!
		// mock.ExpectQuery(
		// 	regexp.QuoteMeta(fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((path = '%s')) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedPath))).
		// 	WillReturnRows(sqlmock.NewRows([]string{"id", "path"}))
		// mock.ExpectQuery(
		// 	regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((path = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`)).
		// 	WithArgs(expectedPath).
		// 	WillReturnRows(sqlmock.NewRows([]string{"id", "path"}))
		mock.ExpectQuery(`SELECT . FROM .+`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "path"}))
		mock.ExpectQuery(`SELECT . FROM .+`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "path"}))
		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "pictures"*`).
			WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), expectedPath, expectedDescription).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("154"))
		mock.ExpectCommit()

		db, err := gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()

		a.SetDB(db)
		a.StartRouter()
	}

	p := Picture{Path: expectedPath, Description: expectedDescription}

	picJSON, err := json.Marshal(p)
	if err != nil {
		t.Errorf("Failed: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "/api/pic/create", bytes.NewBuffer(picJSON))
	if err != nil {
		t.Errorf("Failed to create new request: %v", err)
	}

	recorder := httptest.NewRecorder()

	a.GetRouter().ServeHTTP(recorder, req)

	if recorder.Code == http.StatusOK {
		var body map[string]interface{}
		raw, err := ioutil.ReadAll(recorder.Result().Body)
		if err != nil {
			t.Errorf("Failed to read body from create response.")
		}
		
		err = json.Unmarshal(raw, &body)
		if err != nil {
			t.Errorf("Failed to unmarshalled")
		}

		if body["picture"] == nil {
			t.Errorf("Response does not contain picture entity.")
		}
	}
}

func TestGetRandomPicAPINoPics(t *testing.T) {
	a := App{}

	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		// TODO improve mocks!
		mock.ExpectQuery(`SELECT . FROM .+`).WillReturnError(gorm.ErrRecordNotFound)

		db, err := gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()

		a.SetDB(db)
		a.StartRouter()
	}

	req, err := http.NewRequest(http.MethodGet, "/api/pic/random", nil)
	if err != nil {
		t.Errorf("Failed to create new request: %v", err)
	}

	recorder := httptest.NewRecorder()

	a.GetRouter().ServeHTTP(recorder, req)

	if recorder.Code == http.StatusOK {
		var body map[string]interface{}
		raw, err := ioutil.ReadAll(recorder.Result().Body)
		if err != nil {
			t.Errorf("Failed to read body from create response.")
		}
		
		err = json.Unmarshal(raw, &body)
		if err != nil {
			t.Errorf("Failed to unmarshalled")
		}

		if body["picture"] != nil {
			var picReturned map[string]interface{}
			picReturned = body["picture"].(map[string]interface{})
			if int(picReturned["ID"].(float64)) != 0  {
				t.Errorf("Response does not contain picture entity! -Picture: %v", picReturned)
			}
		}
	}
}