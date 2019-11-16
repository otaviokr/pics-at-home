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

	// Database Setup
	dbSetting, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Fatalf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidateMissingRegistry(t *testing.T) {
	// Parameters used during the test.
	expectedID := 25
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, fmt.Sprintf("Picture not found in Database: %d", expectedID))

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
		WillReturnError(gorm.ErrRecordNotFound)

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Fatalf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidateConnectionError(t *testing.T) {
	// Parameters used during the test.
	expectedID := 25
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, "Connection error. Failed to validated picture ID.")

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
		WillReturnError(fmt.Errorf("Testing error"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Fatalf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidateNoIDUnexpectedError(t *testing.T) {
	// Parameters used during the test.
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, fmt.Sprintf("Unexpected error while validating picture path in database: %s", expectedPath))

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."path" = %s)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedPath))).
		WillReturnError(fmt.Errorf("Testing error"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	// Actual test.
	pic := &Picture{}
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Fatalf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidatePictureWithID(t *testing.T) {
	// Parameters used during the test.
	expectedID := 25
	expectedPath := "test.jpg"
	expectedValidated := true
	expectedMessage := utils.Message(false, "Validation passed")

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}).
				AddRow(expectedID, fmt.Sprintf("NOT_%s", expectedPath)))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	// Validation.
	if actualValidated != expectedValidated {
		t.Fatalf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidatePictureWithPathNotRegistered(t *testing.T) {
	expectedPath := "test.jpg"
	expectedValidated := true
	expectedMessage := utils.Message(false, "Validation passed")

	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((path = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`)).
		WithArgs(expectedPath).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	pic := &Picture{}
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	if actualValidated != expectedValidated {
		t.Fatalf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestValidatePictureWithPathAlreadyRegistered(t *testing.T) {
	expectedID := 25
	expectedPath := "test.jpg"
	expectedValidated := false
	expectedMessage := utils.Message(false, fmt.Sprintf("Picture path already stored in database: %s", expectedPath))

	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((path = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`)).
		WithArgs(expectedPath).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}).
				AddRow(expectedID, expectedPath))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	pic := &Picture{}
	pic.Path = expectedPath

	actualMessage, actualValidated := pic.Validate(db)

	if actualValidated != expectedValidated {
		t.Fatalf("Validation mismatch! actual[%v] - expected[%v]\n%v", actualValidated, expectedValidated, actualMessage)
	}

	if actualMessage["status"] != expectedMessage["status"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}

	if actualMessage["message"] != expectedMessage["message"] {
		t.Fatalf("Message mismatch!\n\tactual\n%s\n\n\texpected\n%s", actualMessage, expectedMessage)
	}
}

func TestCreateFailValidation(t *testing.T) {
	dbSetting, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)

	pic := &Picture{}
	actual := pic.Create(db)

	if actual == nil {
		t.Fatal("actual is not empty")
	}
}

func TestCreateFailInsert(t *testing.T) {
	// TODO Write test!
}

func TestCreateFailID(t *testing.T) {
	// TODO Write test!
}

func TestCreate(t *testing.T) {
	// TODO Write test!
}

func TestGetRandomPicture(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
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

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual == nil {
		t.Fatalf("actual is nil!")
	}
}

func TestGetRandomPictureFailAllIDs(t *testing.T) {
	// Parameters used during the test.
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
		WillReturnError(fmt.Errorf("failed"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Fatalf("actual is not nil!")
	}
}

func TestGetRandomPictureFailNoID(t *testing.T) {
	// Parameters used during the test.
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/not_a_picture.txt"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Fatalf("actual is not nil!")
	}
}

func TestGetRandomPictureFailSpecificID(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}).
				AddRow(expectedID, expectedPath))

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
		WillReturnError(fmt.Errorf("failed"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Fatalf("actual not is nil!")
	}
}

func TestGetRandomPictureFailFileNotFound(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath := "example.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
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

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Fatalf("actual is not nil!")
	}
}

func TestGetRandomPictureFailDecodeJPEG(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/not_a_picture.txt"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
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

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPicture(db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Fatalf("actual is not nil!")
	}
}

func TestGetRandomPictureInfo(t *testing.T) {
	// Parameters used during the test.
	expectedID := 33
	expectedPath := "test.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
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

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	if actual.ID != uint(expectedID) {
		t.Fatalf("ID mismatch! actual[%v] - expected[%v]", actual, expectedID)
	}

	if actual.Path != expectedPath {
		t.Fatalf("Path mismatch!\n\tactual\n%s\n\n\texpected\n%s", actual.Path, expectedPath)
	}
}

func TestGetRandomPictureInfoFailAllIDs(t *testing.T) {
	// Parameters used during the test.
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"
	expected := Picture{}

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
		WillReturnError(fmt.Errorf("failed"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	// TODO improve this test!
	if actual != expected {
		t.Fatalf("actual is not empty!")
	}
}

func TestGetRandomPictureInfoFailNoID(t *testing.T) {
	// Parameters used during the test.
	expected := Picture{}
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/not_a_picture.txt"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	// TODO improve this test!
	if actual != expected {
		t.Fatalf("actual is not nil!")
	}
}

func TestGetRandomPictureInfoFailSpecificID(t *testing.T) {
	// Parameters used during the test.
	expected := Picture{}
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY id ASC`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}).
				AddRow(expectedID, expectedPath))

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND (("pictures"."id" = %d)) ORDER BY "pictures"."id" ASC LIMIT 1`, expectedID))).
		WillReturnError(fmt.Errorf("failed"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRandomPictureInfo(db)

	// Validation.
	// TODO improve this test!
	if actual != expected {
		t.Fatalf("actual not is nil!")
	}
}

func TestGetRecentPics(t *testing.T) {
	// Parameters used during the test.
	expectedNumberOfPics := 5
	expectedID := 33
	expectedPath := "test.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT %d`, expectedNumberOfPics))).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}).
				AddRow(expectedID, expectedPath))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actual := GetRecentPics(uint(expectedNumberOfPics), db)

	// Validation.
	if actual[0].ID != uint(expectedID) {
		t.Fatalf("ID mismatch! actual[%v] - expected[%v]", actual, expectedID)
	}

	if actual[0].Path != expectedPath {
		t.Fatalf("Path mismatch!\n\tactual\n%s\n\n\texpected\n%s", actual[0].Path, expectedPath)
	}
}

func TestGetRecentPicsErrorDB(t *testing.T) {
	// Parameters used during the test.
	expectedNumberOfPics := 10

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL ORDER BY created_at DESC LIMIT %d`, expectedNumberOfPics))).
		WillReturnError(fmt.Errorf("test error"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetRecentPics(uint(expectedNumberOfPics), db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Fatalf("actual not is nil!")
	}
}

func TestGetPictureByID(t *testing.T) {
	// Parameters used during the test.
	expectedID := 33
	expectedPath := "test.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((id = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`))).
		WithArgs(expectedID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "path"}).
				AddRow(expectedID, expectedPath))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	pic := &Picture{}
	pic.ID = uint(expectedID)
	pic.Path = expectedPath

	actual := GetPictureByID(uint(expectedID), db)

	// Validation.
	if actual.ID != uint(expectedID) {
		t.Fatalf("ID mismatch! actual[%v] - expected[%v]", actual, expectedID)
	}

	if actual.Path != expectedPath {
		t.Fatalf("Path mismatch!\n\tactual\n%s\n\n\texpected\n%s", actual.Path, expectedPath)
	}
}

func TestGetPictureByIDFailErrorDB(t *testing.T) {
	// Parameters used during the test.
	expectedID := 89
	expectedPath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not define current location: %v", err)
	}
	expectedPath = expectedPath + "/../test_aux/example.jpg"

	// Database Setup
	dbSetting, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(
			fmt.Sprintf(`SELECT * FROM "pictures" WHERE "pictures"."deleted_at" IS NULL AND ((id = $1)) ORDER BY "pictures"."id" ASC LIMIT 1`))).
		WithArgs(expectedID).
		WillReturnError(fmt.Errorf("test error"))

	db, err := gorm.Open("postgres", dbSetting)
	if err != nil {
		t.Fatalf("Fail to connect to mock db: %v", err)
	}
	db.LogMode(true)
	defer db.Close()

	// Actual test.
	actual := GetPictureByID(uint(expectedID), db)

	// Validation.
	// TODO improve this test!
	if actual != nil {
		t.Fatalf("actual not is nil!")
	}
}
