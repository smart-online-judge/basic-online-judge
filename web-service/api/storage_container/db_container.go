package storage_container

import (
	"database/sql"
	"fmt"
	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	utils "web-service/api/utils"
	"encoding/json"
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

func (self *DbClientContainer) GetValue(id guuid.UUID) (string, bool) {
	var result string
	err := self.getStatement.QueryRow(id.String()).Scan(&result)
	if err != nil {
		return "", false
	}
	var result_internal [][]float32
	err = json.Unmarshal([]byte(result), &result_internal)
	if err != nil {
	        return "", false
	}
	return result, true
}

func (self *DbClientContainer) SaveClient(id guuid.UUID, res [][]float32) {
       result_binary, err := json.Marshal(res)
	if err != nil {
	        # TODO
	}
	_, err = self.saveStatement.Exec(id.String(), result_binary)
	if err != nil {
		ErrorLogger.Println(err)
	}
}
