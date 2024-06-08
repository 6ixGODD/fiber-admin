package service_test

import (
	"os"
	"testing"

	"fiber-admin/test/common"
)

func TestMain(m *testing.M) {
	if err := common.Setup(); err != nil {
		panic(err)
	}
	code := m.Run()
	if err := common.Teardown(); err != nil {
		panic(err)
	}
	os.Exit(code)
}
