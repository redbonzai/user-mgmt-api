package user_management_api_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUserManagementApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "UserManagementApi Suite")
}
