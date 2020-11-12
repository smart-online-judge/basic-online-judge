package storage_container

import (
	guuid "github.com/google/uuid"
)

type ClientContainer interface {
	GetValue(id guuid.UUID) (string, bool)
	SaveClient(id guuid.UUID, res string)
}
