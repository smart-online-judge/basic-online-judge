package storage_container

import (
	"database/sql"
	"encoding/json"
	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	utils "web-service/src/utils"
)

var (
	ErrorLogger *log.Logger
)

type DbClientContainer struct {
	dbConnection          *sql.DB
	createClientStatement *sql.Stmt
	updateClientStatement *sql.Stmt
	getResStatement       *sql.Stmt
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
			id INTEGER PRIMARY KEY,
			uuid VARCHAR,
			dir_path VARCHAR,
			acessed_time TIMESTAMP,
			status CHARACTER,
			result TEXT
	)`)
	// status - success, error, pending
	// result - JSON error_msg or computed result
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	self.createClientStatement, err = self.dbConnection.Prepare(
		"INSERT INTO CLIENTS (uuid, status, result) VALUES (?, ?, ?)")
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	self.updateClientStatement, err = self.dbConnection.Prepare(
		"UPDATE CLIENTS SET status = ?, result = ? WHERE uuid = ?;")
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	self.getResStatement, err = self.dbConnection.Prepare(
		"SELECT status, result FROM CLIENTS WHERE uuid == ?")
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

func (self *DbClientContainer) GetResValue(id guuid.UUID) (*utils.Pair, bool) {
	var result string
	var status ResStatus
	err := self.getResStatement.QueryRow(id.String()).Scan(&status, &result)
	if err != nil {
		ErrorLogger.Println(err)
		return nil, false
	}

	return &utils.Pair{status, result}, true
}

func (self *DbClientContainer) SavePendingClient(id guuid.UUID, msg string) {
	_, err := self.createClientStatement.Exec(
		id.String(),
		Pending,
		msg)

	if err != nil {
		ErrorLogger.Println(err)
	}
}

func (self *DbClientContainer) SaveErrorClient(id guuid.UUID, err_msg string) {
	_, err := self.updateClientStatement.Exec(
		Error,
		err_msg,
		id.String())

	if err != nil {
		ErrorLogger.Println(err)
	}
}

func (self *DbClientContainer) SaveResClient(id guuid.UUID, result_json string) {
	_, err := self.updateClientStatement.Exec(
		Success,
		result_json,
		id.String())

	if err != nil {
		ErrorLogger.Println(err)
	}
}
