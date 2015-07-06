package model

import (
	"github.com/anominet/anomi/env"
	"reflect"
)

type ModelEnv struct {
	*env.Env
}

var TypePrefixRegistry = map[reflect.Type]string{
	reflect.TypeOf(User{}):      "user",
	reflect.TypeOf(Vote{}):      "vote",
	reflect.TypeOf(Post{}):      "post",
	reflect.TypeOf(PostScore{}): "postscore",
	reflect.TypeOf(int64(42)):    "int64",
	reflect.TypeOf("42"):          "string",
}
