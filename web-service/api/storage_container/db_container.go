package storage_container

import (
	"database/sql"
	// "fmt"
	"encoding/json"
	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	utils "web-service/api/utils"
)

var (
	ErrorLogger *log.Logger
)

type DbClientContainer struct {
	dbConnection  *sql.DB
	saveStatement *sql.Stmt
	getStatement  *sql.Stmt
}

func NewDB() *DbClientContainer {
	return new(DbClientContainer)
}

func (self *DbClientContainer) Initialize(db_path string) {
	ErrorLogger = utils.GetLogger("ERROR: ")
	dir, _ := filepath.Split(db_path)

	var err error
	if err = os.MkdirAll(dir, 0777); err != nil {
		ErrorLogger.Fatal(err)
	}

	if self.dbConnection, err = sql.Open("sqlite3", db_path); err != nil {
		ErrorLogger.Fatal(err)
	}

	_, err = self.dbConnection.Exec(
		`CREATE TABLE IF NOT EXISTS CLIENTS (
            id INTEGER PRIMARY KEY, uuid VARCHAR, dir_path VARCHAR,
            acessed_time TIMESTAMP, result TEXT
    )`)
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	self.saveStatement, err = self.dbConnection.Prepare(
		"INSERT INTO CLIENTS (uuid, result) VALUES (?, ?)")
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	self.getStatement, err = self.dbConnection.Prepare(
		"SELECT result FROM CLIENTS WHERE uuid == ?")
	if err != nil {
		ErrorLogger.Fatal(err)
	}
}

func __temporaryHelper__unmarshal_json_into_2d_slice(result string) [][]float32 {
	var result_internal [][]float32
	err := json.Unmarshal([]byte(result), &result_internal)
	if err != nil {
		ErrorLogger.Println("Cannot unmarshal")
	}
	return result_internal
}

func (self *DbClientContainer) GetValue(id guuid.UUID) (string, bool) {
	var result string
	err := self.getStatement.QueryRow(id.String()).Scan(&result)
	if err != nil {
		ErrorLogger.Println(err)
		return "", false
	}
	return result, true
}

func (self *DbClientContainer) SaveClient(id guuid.UUID, result_json string) {
	_, err := self.saveStatement.Exec(id.String(), result_json)
	if err != nil {
		ErrorLogger.Println(err)
	}
}
