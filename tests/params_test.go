package test

import (
	"fmt"
	"mas/utils/params"
	"testing"
)


type test struct {
	test int
}

func TestParams(t *testing.T) {


	params := paramsUtils.NewParamsParser(test{})
	fmt.Println(params.Has("test"))
	a, b := params.Int("test1", "123", paramsUtils.Config{Require: false, DefaultValue: 11111,})
	fmt.Println(a, b)

}




