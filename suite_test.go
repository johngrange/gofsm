package fsm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGofsm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gofsm Suite")
}
