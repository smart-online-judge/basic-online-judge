package storage_container

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	utils "web-service/src/utils"
)

type DbClientService struct {
	dbConnection *sql.DB
}

func NewDbClientService(db_path string) (*DbClientService, error) {
	db := new(DbClientService)
	err := db.initialize(db_path)
	return db, err
}

func (self *DbClientService) Close() error {
	err := self.dbConnection.Close()
	return err
}

func (self *DbClientService) initialize(db_path string) error {
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
	);`)
	// status - success, error, pending
	// result - JSON error_msg or computed result
	if err != nil {
		return err
	}

	return nil
}

func (self *DbClientService) GetResValue(id guuid.UUID) (*utils.Pair, error) {
	sqlStmt := "SELECT status, result FROM CLIENTS WHERE uuid == ?;"
	var result string
	var status ResStatus
	err := self.dbConnection.QueryRow(sqlStmt, id.String()).Scan(&status, &result)
	if err != nil {
		return nil, err
	}

	return &utils.Pair{status, result}, err
}

func (self *DbClientService) SavePendingClient(id guuid.UUID, msg string) error {
	sqlStmt := "INSERT INTO CLIENTS (uuid, creation_time, status, result) VALUES (?, ?, ?, ?);"
	_, err := self.dbConnection.Exec(
		sqlStmt, id.String(), time.Now().Unix(), Pending, msg)

	return err
}

func (self *DbClientService) SaveErrorClient(id guuid.UUID, err_msg string) error {
	sqlStmt := "UPDATE CLIENTS SET status = ?, result = ? WHERE uuid = ?;"
	_, err := self.dbConnection.Exec(
		sqlStmt, Error, err_msg, id.String())

	return err
}

func (self *DbClientService) SaveSuccessClient(id guuid.UUID, result_json string) error {
	sqlStmt := "UPDATE CLIENTS SET status = ?, result = ? WHERE uuid = ?;"
	_, err := self.dbConnection.Exec(
		sqlStmt, Success, result_json, id.String())

	return err
}

func (self *DbClientService) ClientExists(id guuid.UUID) bool {
	sqlStmt := "SELECT EXISTS(SELECT id FROM CLIENTS WHERE uuid == ?);"

	var exists bool
	err := self.dbConnection.QueryRow(sqlStmt, id.String()).Scan(&exists)
	return err == nil && exists
}

func (self *DbClientService) GetRemoveStaleClients(back_period int64) ([]string, error) {
	selectStmt := `SELECT uuid FROM CLIENTS WHERE ABS(? - creation_time) > ?;`
	deleteStmt := `DELETE FROM CLIENTS WHERE ABS(? - creation_time) > ?;`

	cur_time := time.Now().Unix()
	rows, err := self.dbConnection.Query(selectStmt, cur_time, back_period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resIds []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err == nil {
			resIds = append(resIds, id)
		}
	}

	_, err = self.dbConnection.Exec(deleteStmt, cur_time, back_period)
	if err != nil {
		return nil, err
	}

	return resIds, nil
}
