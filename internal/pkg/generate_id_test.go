package pkg_test

import (
	"campaign/internal/pkg"
	"fmt"
	"testing"
)

func TestGenerateID(t *testing.T) {
	id, _ := pkg.GenerateUniqueID(10)
	fmt.Println(id)
}
