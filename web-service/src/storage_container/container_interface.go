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
	GetResValue(id guuid.UUID) (*utils.Pair, error)
	SavePendingClient(id guuid.UUID, res string) error
	SaveErrorClient(id guuid.UUID, res string) error
	SaveSuccessClient(id guuid.UUID, res string) error
}
