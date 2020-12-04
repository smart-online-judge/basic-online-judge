package storage_container

import (
	"database/sql"
	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"time"
	utils "web-service/src/utils"
)

type DbClientContainer struct {
	dbConnection *sql.DB

	createClientStmt     *sql.Stmt
	updateClientResStmt  *sql.Stmt
	updateClientTimeStmt *sql.Stmt
	getClientResStmt     *sql.Stmt
}

func NewDB() *DbClientContainer {
	return new(DbClientContainer)
}

func (self *DbClientContainer) Close() error {
	self.createClientStmt.Close()
	self.updateClientResStmt.Close()
	self.getClientResStmt.Close()
	err := self.dbConnection.Close()
	return err
}

func (self *DbClientContainer) Initialize(db_path string) error {
	dir, _ := filepath.Split(db_path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	var err error
	if self.dbConnection, err = sql.Open("sqlite3", db_path); err != nil {
		return err
	}

	_, err = self.dbConnection.Exec(
		`CREATE TABLE IF NOT EXISTS CLIENTS (
			id INTEGER PRIMARY KEY,
			uuid VARCHAR,
			creation_time INTEGER,
			status CHARACTER,
			result TEXT
	)`)
	// status - success, error, pending
	// result - JSON error_msg or computed result
	if err != nil {
		return err
	}

	self.createClientStmt, err = self.dbConnection.Prepare(
		"INSERT INTO CLIENTS (uuid, creation_time, status, result) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	self.updateClientResStmt, err = self.dbConnection.Prepare(
		"UPDATE CLIENTS SET status = ?, result = ? WHERE uuid = ?;")
	if err != nil {
		return err
	}

	self.getClientResStmt, err = self.dbConnection.Prepare(
		"SELECT status, result FROM CLIENTS WHERE uuid == ?")
	if err != nil {
		return err
	}

	return nil
}

func (self *DbClientContainer) GetResValue(id guuid.UUID) (*utils.Pair, error) {
	var result string
	var status ResStatus
	err := self.getClientResStmt.QueryRow(id.String()).Scan(&status, &result)
	if err != nil {
		return nil, err
	}

	return &utils.Pair{status, result}, err
}

func (self *DbClientContainer) SavePendingClient(id guuid.UUID, msg string) error {
	_, err := self.createClientStmt.Exec(
		id.String(), time.Now().Unix(), Pending, msg)

	return err
}

func (self *DbClientContainer) SaveErrorClient(id guuid.UUID, err_msg string) error {
	_, err := self.updateClientResStmt.Exec(
		Error, err_msg, id.String())

	return err
}

func (self *DbClientContainer) SaveSuccessClient(id guuid.UUID, result_json string) error {
	_, err := self.updateClientResStmt.Exec(
		Success,
		result_json,
		id.String())

	return err
}

func (self *DbClientContainer) ClientExists(id guuid.UUID) bool {
	sqlStmt := "SELECT EXISTS(SELECT * FROM CLIENTS WHERE uuid == ?);"

	res, err := self.dbConnection.Exec(sqlStmt, id)
	if err != nil {
		return false;
	}

	rows, err := res.RowsAffected()
	return err == nil && rows == 1
}
