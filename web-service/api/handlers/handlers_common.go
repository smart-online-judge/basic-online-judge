package handlers

import (
	"log"
	utils "web-service/api/utils"
)

var (
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

func PrepareHandlersCommon()  {
	WarningLogger = utils.GetLogger("WARNING: ")
	ErrorLogger = utils.GetLogger("ERROR: ")
	DebugLogger = utils.GetLogger("DEBUG: ")
}