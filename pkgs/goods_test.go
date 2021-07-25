package pkgs

import (
	"testing"

	"github.com/amead24/gotraders/pkgs/user"
)

func listGoods_test(t *testing.T) {
	creds, err := user.GetUsernameAndToken()

	output := ListGoods()

}
