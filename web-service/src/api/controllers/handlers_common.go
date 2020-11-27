package controllers

import (
	"log"
	utils "web-service/src/utils"
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

func InitializeHandlersCommon() {
	WarningLogger = utils.GetLogger("WARNING: ")
	ErrorLogger = utils.GetLogger("ERROR: ")
	DebugLogger = utils.GetLogger("DEBUG: ")
}
