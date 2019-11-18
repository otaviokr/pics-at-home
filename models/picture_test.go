package models

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/otaviokr/pics-at-home/utils"
)

func TestValidateMissingPath(t *testing.T) {
	// Parameters used during the test.
	expectedID := 25
	expectedValidated := false
	expectedMessage := utils.Message(false, "Pic entity does not have a valid Path attribute!")

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, _, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Errorf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidateMissingRegistry(t *testing.T) {
	// Parameters used during the test.
	expectedID := 25
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, fmt.Sprintf("Picture not found in Database: %d", expectedID))

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			WillReturnError(gorm.ErrRecordNotFound)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Errorf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidateConnectionError(t *testing.T) {
	// Parameters used during the test.
	expectedID := 25
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, "Connection error. Failed to validated picture ID.")

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			WillReturnError(gorm.ErrCantStartTransaction)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		//db = a.GetDB()
		
		t.Log("Ignoring this test for integration until a way to simulate connection error can be defined.")
		return
	}

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Errorf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidateNoIDUnexpectedError(t *testing.T) {
	// Parameters used during the test.
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, 
		fmt.Sprintf("Unexpected error while validating picture path in database: %s - Error: %s", expectedPath, gorm.ErrCantStartTransaction))

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(`SELECT . FROM *`).WillReturnError(gorm.ErrCantStartTransaction)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()

		t.Log("Ignoring this test for integration until a way to simulate expected error can be defined.")
		return
	}

	// Actual test.
	pic := &Picture{}
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Errorf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidatePictureWithID(t *testing.T) {
	// Parameters used during the test.
	expectedID := 1
	expectedPath := "test.jpg"
	expectedValidated := true
	expectedMessage := utils.Message(false, "Validation passed")

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, fmt.Sprintf("NOT_%s", expectedPath)))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()

		db.Delete(Picture{})
		p := Picture{Path: "test.jpg"}
		db.Create(&p)
		expectedID = int(p.ID)
	}

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Errorf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidatePictureWithPathNotRegistered(t *testing.T) {
	expectedPath := "test.jpg"
	expectedValidated := true
	expectedMessage := utils.Message(false, "Validation passed")

	var db *gorm.DB
	if testing.Short() {
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((path = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`)).
			WithArgs(expectedPath).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	pic := &Picture{}
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	if actualValidated != expectedValidated {
		t.Errorf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidatePictureWithPathAlreadyRegistered(t *testing.T) {
	expectedID := 25
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, fmt.Sprintf("Picture path already stored in database: %s", expectedPath))

	var db *gorm.DB
	if testing.Short() {
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((path = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`)).
			WithArgs(expectedPath).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()

		db.Delete(Picture{})
		p := Picture{Path: "test.jpg"}
		db.Create(&p)
	}

	pic := &Picture{}
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	if actualValidated != expectedValidated {
		t.Errorf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Errorf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestCreateFailValidation(t *testing.T) {
	var db *gorm.DB
	if testing.Short() {
		dbSetting, _, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		a.GetDB()
	}

	pic := &Picture{}
	actual := pic.Create(db)

	if actual == nil {
		t.Error("actual is not empty")
	}
}

func TestCreateFailInsert(t *testing.T) {
	// TODO Write test!
	t.Log("TestCreateFailInsert not implemented. Moving on...")
}

func TestCreateFailID(t *testing.T) {
	// TODO Write test!
	t.Log("TestCreateFailID not implemented. Moving on...")
}

func TestCreate(t *testing.T) {
	// TODO Write test!
	t.Log("TestCreate not implemented. Moving on...")
}

func TestGetRandomPicture(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()

		db.Delete(Picture{})
		p := Picture{Path: "test.jpg"}
		db.Create(&p)
	}

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual == nil {
		t.Errorf("actual is nil!")
	}
}

func TestGetRandomPictureFailAllIDs(t *testing.T) {
	// Parameters used during the test.
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			//WillReturnError(fmt.Errorf("failed"))
			WillReturnError(gorm.ErrRecordNotFound)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Errorf("actual is not nil!")
	}
}

func TestGetRandomPictureFailNoID(t *testing.T) {
	// Parameters used during the test.
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/not_a_picture.txt"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Errorf("actual is not nil!")
	}
}

func TestGetRandomPictureFailSpecificID(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			//WillReturnError(fmt.Errorf("failed"))
			WillReturnError(gorm.ErrRecordNotFound)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Errorf("actual not is nil!")
	}
}

func TestGetRandomPictureFailFileNotFound(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath := "example.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Errorf("actual is not nil!")
	}
}

func TestGetRandomPictureFailDecodeJPEG(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/not_a_picture.txt"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Errorf("actual is not nil!")
	}
}

func TestGetRandomPictureInfo(t *testing.T) {
	// Parameters used during the test.
	expectedID := 33
	expectedPath := "test.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()

		db.Delete(Picture{})
		p := Picture{Path: "test.jpg"}
		db.Create(&p)
		expectedID = int(p.ID)
	}

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	if actual.ID != uint(expectedID) {
		t.Errorf("ID mismatch! actual[%v] - expected[%v]", actual, expectedID)
	}

	if actual.Path != expectedPath {
		t.Errorf("Path mismatch!\n\tactual\n%s\n\n\texpected\n%s", actual.Path, expectedPath)
	}
}

func TestGetRandomPictureInfoFailAllIDs(t *testing.T) {
	// Parameters used during the test.
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"
	expected := Picture{}

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			//WillReturnError(fmt.Errorf("failed"))
			WillReturnError(gorm.ErrRecordNotFound)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	// TODO improve this test!
	if actual != expected {
		t.Errorf("actual is not empty!")
	}
}

func TestGetRandomPictureInfoFailNoID(t *testing.T) {
	// Parameters used during the test.
	expected := Picture{}
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/not_a_picture.txt"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	// TODO improve this test!
	if actual != expected {
		t.Errorf("actual is not nil!")
	}
}

func TestGetRandomPictureInfoFailSpecificID(t *testing.T) {
	// Parameters used during the test.
	expected := Picture{}
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
			//WillReturnError(fmt.Errorf("failed"))
			WillReturnError(gorm.ErrRecordNotFound)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	// TODO improve this test!
	if actual != expected {
		t.Errorf("actual not is nil!")
	}
}

func TestGetRecentPics(t *testing.T) {
	// Parameters used during the test.
	expectedNumberOfPics := 5
	expectedID := 33
	expectedPath := "test.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT %d`, expectedNumberOfPics))).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
		db.Delete(Picture{})
		p := Picture{Path: "test.jpg"}
		db.Create(&p)
		expectedID = int(p.ID)
	}

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actual := GetRecentPics(uint(expectedNumberOfPics), db)

	// Validation.
	if len(actual) < 1 {
		t.Error("No valid picture retrieved!")
	}
	if actual[0].ID != uint(expectedID) {
		t.Errorf("ID mismatch! actual[%v] - expected[%v]", actual, expectedID)
	}

	if actual[0].Path != expectedPath {
		t.Errorf("Path mismatch!\n\tactual\n%s\n\n\texpected\n%s", actual[0].Path, expectedPath)
	}
}

func TestGetRecentPicsErrorDB(t *testing.T) {
	// Parameters used during the test.
	expectedNumberOfPics := 10

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT %d`, expectedNumberOfPics))).
			//WillReturnError(fmt.Errorf("failed"))
			WillReturnError(gorm.ErrRecordNotFound)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetRecentPics(uint(expectedNumberOfPics), db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Errorf("actual not is nil!")
	}
}

func TestGetPictureByID(t *testing.T) {
	// Parameters used during the test.
	expectedID := 33
	expectedPath := "test.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((id = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`))).
			WithArgs(expectedID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "path"}).
					AddRow(expectedID, expectedPath))

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actual := GetPictureByID(uint(expectedID), db)

	// Validation.
	if actual.ID != uint(expectedID) {
		t.Errorf("ID mismatch! actual[%v] - expected[%v]", actual, expectedID)
	}

	if actual.Path != expectedPath {
		t.Errorf("Path mismatch!\n\tactual\n%s\n\n\texpected\n%s", actual.Path, expectedPath)
	}
}

func TestGetPictureByIDFailErrorDB(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	var db *gorm.DB
	if testing.Short() {
		// Database Setup
		dbSetting, mock, err := sqlmock.New()
		if err != nil {
			t.Errorf("An error '%s' was not expected when opening a stub database connection", err)
		}

		mock.ExpectQuery(
			regexp.QuoteMeta(
				fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((id = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`))).
			WithArgs(expectedID).
			//WillReturnError(fmt.Errorf("failed"))
			WillReturnError(gorm.ErrRecordNotFound)

		db, err = gorm.Open("postgres", dbSetting)
		if err != nil {
			t.Errorf("Fail to connect to mock db: %v", err)
		}
		//db.LogMode(true)
		defer db.Close()
	} else {
		db = a.GetDB()
	}

	// Actual test.
	actual := GetPictureByID(uint(expectedID), db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Errorf("actual not is nil!")
	}
}
