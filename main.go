package main

import (
	"crypto/sha1"
	"database/sql"
	_ "database/sql"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	_ "fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	_ "log"
	"net/http"
	_ "net/http"
	"os"
	_ "os"
	"strconv"
	_ "strconv"
	"strings"
	_ "github.com/google/uuid"
	_ "crypto/sha1"
	_ "github.com/gorilla/handlers"
)

// MESSAGES
const(
	GENERAL_MESSAGE_APP_PORT_LISTENING = "Listening on port: "
	EXCEPTION_ERROR_PREFIX = "Exception Error:"
	DB_ERROR_EXCEPTION_INSERT = " Could not insert new application into Database"

	SUCCESS_MESSAGE_PREFIX = "Success: "
	SUCCESS_MESSAGE_GET_ALL_APP = " Get all applications"
	SUCCESS_MESSAGE_ADD_NEW_APP = "Added new application"

	ERROR_MESSAGE_UNAUTHORIZED_PERMISSION_REQUIRE = "Access Denied, you need permission to access this site"

	ERROR_MESSAGE_INVALID_REQUEST_PARAMS = "Invalid request params"
	ERROR_MESSAGE_APP_KEY_AND_NAME_REQUIRED = "app_key And Name parameters are required"
	ERROR_MESSAGE_APP_KEY_REQUIRED = "app_key parameter is required"
	ERROR_MESSAGE_APP_NAME_REQUIRED = "app_name parameter is required"
	ERROR_MESSAGE_APP_NAME_VALIDATION = "app_name parameter was over the limit of %s"
	ERROR_MESSAGE_APP_KEY_VALIDATION = "app_key parameter must be a valid uuid key parameter"
)

const (
	INSERT_NEW_APP_QUERY = "insert into application(app_key, app_name, app_uuid) values($1, $2, $3)"
	SELECT_ALL_APPS_QUERY = "SELECT app_id, app_key, app_name, creation_time, app_uuid FROM application"
)

// GENERAL PROPERTIES
const(
	AUTHENTICATION_KEY = "Checkpointid"
	AUTHENTICATION_VALUE = "let-me-pass"
	CONTENT_TYPE_KEY = "Content-Type"
	X_REQUIRED_WITH_KEY = "X-Requested-With"
	AUTHORIZATION_KEY = "Authorization"
	CONTENT_TYPE_VALUE = "application/json"
	SERVICE_METHOD_OPTIONS = "OPTIONS"
	SERVICE_METHOD_HEAD = "HEAD"
	SERVICE_METHOD_POST = "POST"
	SERVICE_METHOD_PUT = "PUT"
	SERVICE_METHOD_GET = "GET"
	SERVICE_ADD_NEW_APPLICATION = "/addnewapplication"
	SERVICE_GET_ALL_APPLICATIONS = "/getallapplications"
)

// APPLICATION ENV PARAMETERS KEYS
const (
	APP_SERVICE_PORT = "APP_SERVICE_PORT"
	PG_URL = "PG_URL"
	PG_HOST = "PG_HOST"
	PG_PASSWORD = "PG_PASSWORD"
	PG_PORT = "PG_PORT"
	PG_DBNAME = "PG_DBNAME"
	PG_USER = "PG_USER"
)

// VALIDATION CONSTANTS
const (
	ASTERISK = "*"
	EMPTY_LENGTH = 0;
	SEPARATOR_CHAR = "-"
	SEPARATOR_COLON = ":"
	VALID_UUID_SEPARATOR_CHAR = 5;
	VALID_UUID_CHAR_LIMIT = 36;
	APPLICATION_TABLE_DB_PROPERTY_NAME_LIMIT = 45;
)

// Application manager
type AppManager struct {
	DB * sql.DB
}

// Application response structure
type AppResponeGeneral struct {
	Code int
	Desc string
}

// Application data response structure
type AppDataRespone struct {
	Code int
	Desc string
	Applications applications
}

// Application request structure
type NewAppRequest struct {
	AppName string  `json:"app_name"`
	AppKey string `json:"app_key"`
}

// Application response structure

type Application struct {
	AppUuid string `json:"app_uuid"`
	AppCreationTime string `json:"creation_time"`
	AppName string  `json:"app_name"`
	AppKey string `json:"app_key"`
	AppId string `json:"app_id"`
}
type applications [] Application

// Add new Application structure
type NewAppParams struct {
	AppUuid string
	AppName string
	AppKey string
	AppId string
}


// Add new Application
func (appManager * AppManager) addApplication(newApp NewAppParams) (error) {
	_ , err := appManager.DB.Exec(INSERT_NEW_APP_QUERY,
		newApp.AppKey, newApp.AppName, newApp.AppUuid);
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal(EXCEPTION_ERROR_PREFIX + DB_ERROR_EXCEPTION_INSERT)
		return err
	}

	defer appManager.DB.Close()

	return nil

}

// Get all Applications
func (appManager * AppManager)getApplications() ([]Application, error) {
	rows, err := appManager.DB.Query(SELECT_ALL_APPS_QUERY)

	if err != nil {
		return nil, err
	}

	applications := []Application{}

	for rows.Next() {

		var app Application
		if err := rows.Scan(&app.AppId, &app.AppKey, &app.AppName, &app.AppCreationTime, &app.AppUuid); err != nil {
			return nil, err
		}
		app.AppId = strings.TrimSpace(app.AppId)
		app.AppUuid = strings.TrimSpace(app.AppUuid)
		app.AppKey = strings.TrimSpace(app.AppKey)
		app.AppName = strings.TrimSpace(app.AppName)
		app.AppCreationTime = strings.TrimSpace(app.AppCreationTime)
		applications = append(applications, app)
	}

	defer appManager.DB.Close()

	return applications, nil
}

// Get Applications Process
func (appManager * AppManager) getApp(response http.ResponseWriter, request *http.Request) {

	if validateRequestHeaderCookieParam(request) == false {
		respondWithError(response, http.StatusUnauthorized, ERROR_MESSAGE_UNAUTHORIZED_PERMISSION_REQUIRE)
		return
	}

	appManager.connectToDb();
	applications, _ := appManager.getApplications();

	responseGetApplications(response, http.StatusOK, SUCCESS_MESSAGE_PREFIX + SUCCESS_MESSAGE_GET_ALL_APP, applications);
}

// Connect to the Database
func (appManager * AppManager) connectToDb() {

	value, err_conv := strconv.Atoi(os.Getenv(PG_PORT))
	if err_conv != nil {
		log.Fatal(err_conv)
		return
	}
	psqlInfo := fmt.Sprintf(
		os.Getenv(PG_URL),
		os.Getenv(PG_HOST),
		value,
		os.Getenv(PG_USER),
		os.Getenv(PG_PASSWORD),
		os.Getenv(PG_DBNAME))


	var err error

	appManager.DB, err = sql.Open(os.Getenv(PG_USER), psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

}

// Add new Application Process
func (appManager * AppManager) addApp(response http.ResponseWriter, request *http.Request) {
	var newAppReq NewAppRequest

	if validateRequestHeaderCookieParam(request) == false {
		respondWithError(response, http.StatusUnauthorized, ERROR_MESSAGE_UNAUTHORIZED_PERMISSION_REQUIRE)
		return
	}

	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&newAppReq); err != nil {
		respondWithError(response, http.StatusBadRequest, ERROR_MESSAGE_INVALID_REQUEST_PARAMS)
		return
	}
	defer request.Body.Close()

	if len(newAppReq.AppName) <= EMPTY_LENGTH && len(newAppReq.AppKey) <= EMPTY_LENGTH {
		respondWithError(response, http.StatusBadRequest, ERROR_MESSAGE_APP_KEY_AND_NAME_REQUIRED)
		return
	}

	if len(newAppReq.AppKey) <= EMPTY_LENGTH {
		respondWithError(response, http.StatusBadRequest, ERROR_MESSAGE_APP_KEY_REQUIRED)
		return
	}

	if len(newAppReq.AppName) <= EMPTY_LENGTH {
		respondWithError(response, http.StatusBadRequest, ERROR_MESSAGE_APP_NAME_REQUIRED)
		return
	}

	if len(newAppReq.AppName) > APPLICATION_TABLE_DB_PROPERTY_NAME_LIMIT {
		respondWithError(response, http.StatusBadRequest, fmt.Sprintf(ERROR_MESSAGE_APP_NAME_VALIDATION, APPLICATION_TABLE_DB_PROPERTY_NAME_LIMIT))
		return
	}

	app_key_splited := strings.Split(newAppReq.AppKey, SEPARATOR_CHAR);
	if  (len(newAppReq.AppKey) != VALID_UUID_CHAR_LIMIT || len(app_key_splited) != VALID_UUID_SEPARATOR_CHAR) {
		respondWithError(response, http.StatusBadRequest, ERROR_MESSAGE_APP_KEY_VALIDATION)
		return
	}

	// Encrypt with sha1 the app_key parameter before insert it into db
	encrypt := sha1.New()
	encrypt.Write([]byte(newAppReq.AppKey))
	encrypted := encrypt.Sum(nil)

	// Generate uuid for id + uuid
	uuid_str_app_id := uuid.New()
	uuid_str_app_uuid := uuid.New()
	appManager.connectToDb();

	newApp := NewAppParams{
		AppUuid: uuid_str_app_uuid.String(),
		AppKey:  fmt.Sprintf("%x", encrypted),
		AppName: newAppReq.AppName,
		AppId: uuid_str_app_id.String()}

	appManager.addApplication(newApp);

	respondGeneral(response, http.StatusOK, SUCCESS_MESSAGE_ADD_NEW_APP);
}

// Get all apps response
func responseGetApplications(
	response http.ResponseWriter,
	code int,
	message string,
	applications applications) {

	appRes := AppDataRespone{Code: code, Desc: message, Applications: applications}
	response.Header().Set(CONTENT_TYPE_KEY, CONTENT_TYPE_VALUE)
	response.WriteHeader(code)
	json.NewEncoder(response).Encode(appRes)
}

// Respond JSON with error
func respondWithError(response http.ResponseWriter, code int, message string) {
	respondGeneral(response, code, message)
}

// General respond JSON
func respondGeneral(response http.ResponseWriter, code int, desc string) {

	appRes := AppResponeGeneral{Code: code, Desc: desc}
	response.Header().Set(CONTENT_TYPE_KEY, CONTENT_TYPE_VALUE)
	response.WriteHeader(code)
	json.NewEncoder(response).Encode(appRes)
}

func validateRequestHeaderCookieParam(request *http.Request) bool {

	headerReq := request.Header.Get(AUTHENTICATION_KEY);
	if len(headerReq) <= EMPTY_LENGTH {
		return false
	}
	if strings.Compare(headerReq, AUTHENTICATION_VALUE) != 0 {
		return false
	}
	return true
}

// Main starter
func main() {

	// Get default port for the service
	appPort := os.Getenv(APP_SERVICE_PORT)

	// Get Application manager object
	appManager := AppManager{}

	// Get router object
	router := mux.NewRouter()

	// Handle routes
	router.HandleFunc(SERVICE_ADD_NEW_APPLICATION, appManager.addApp).Methods(SERVICE_METHOD_PUT)
	router.HandleFunc(SERVICE_GET_ALL_APPLICATIONS, appManager.getApp).Methods(SERVICE_METHOD_GET)

	fmt.Println(GENERAL_MESSAGE_APP_PORT_LISTENING + appPort)
	log.Fatal(http.ListenAndServe(SEPARATOR_COLON + appPort, handlers.CORS(handlers.AllowedHeaders([]string{X_REQUIRED_WITH_KEY, CONTENT_TYPE_KEY, AUTHORIZATION_KEY, AUTHENTICATION_KEY}), handlers.AllowedMethods([]string{SERVICE_METHOD_GET, SERVICE_METHOD_POST, SERVICE_METHOD_PUT, SERVICE_METHOD_HEAD, SERVICE_METHOD_OPTIONS}), handlers.AllowedOrigins([]string{ASTERISK}))(router)))
}