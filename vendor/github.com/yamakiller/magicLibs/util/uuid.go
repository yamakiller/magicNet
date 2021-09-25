package util

import "github.com/google/uuid"

//SpawnUUID doc
//@Method SpawnUUID @Summary spawn uuid
//@Return (string) uuid
func SpawnUUID() string {
	guid := uuid.New()
	return guid.String()
}
