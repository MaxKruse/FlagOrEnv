package main

import (
	"fmt"

	"github.com/maxkruse/flagorenv"
)

type ExampleStruct struct {
	StringField string
	Int64Field  int64
}

func main() {
	c, err := flagorenv.LoadFlagsOrEnv[ExampleStruct](&flagorenv.Config{
		Prefix:     "test",
		PreferFlag: true,
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", c)
}
