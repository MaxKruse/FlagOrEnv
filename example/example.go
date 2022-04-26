package main

import (
	"fmt"

	"github.com/maxkruse/flagorenv"
)

type ExampleStruct struct {
	StringField string `default:"my default string, as per the tag"`
	IntField    int64
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
