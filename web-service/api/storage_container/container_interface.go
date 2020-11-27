package storage_container

import (
	guuid "github.com/google/uuid"
	utils "web-service/api/utils"
)

type ResStatus string
const (
	Success ResStatus = "Success"
	Error   ResStatus = "Error"
	Pending ResStatus = "Pending"
)

type ClientContainer interface {
	GetResValue(id guuid.UUID) (*utils.Pair, bool)
	SavePendingClient(id guuid.UUID, res string)
	SaveErrorClient(id guuid.UUID, res string)
	SaveResClient(id guuid.UUID, res string)
}
