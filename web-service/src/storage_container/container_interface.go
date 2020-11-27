package storage_container

import (
	guuid "github.com/google/uuid"
	utils "web-service/src/utils"
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
	SaveSuccessClient(id guuid.UUID, res string)
}
