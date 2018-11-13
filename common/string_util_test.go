package common

import (
	"testing"
	"fmt"
)

func TestGetUUID(t *testing.T) {
	uuid := GetUUID()
	fmt.Println(uuid, len(uuid))
}