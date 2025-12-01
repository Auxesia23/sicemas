package auth

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v3"
	_ "github.com/lib/pq"
)

//go:embed model.conf
var casbinModel string

func NewEnforcer() (*casbin.Enforcer, error) {
	uri := os.Getenv("DB_URI")
	a, err := xormadapter.NewAdapter("postgres", uri, true)
	if err != nil {
		return nil, err
	}

	m, err := model.NewModelFromString(casbinModel)
	if err != nil {
		return nil, err
	}

	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	e.EnableAutoSave(true)
	groups, _ := e.GetUsersForRole("dev")
	fmt.Println("GROUPS : ", groups)

	return e, nil
}
