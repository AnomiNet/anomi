package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"reflect"
	"testing"
)

type TestStruct struct {
	TestString string `json:"test_string"`
	TestInt64  int64  `json:"test_int64"`
}

var TypePrefixRegistry = map[reflect.Type]string{
	reflect.TypeOf(TestStruct{}): "teststruct",
	reflect.TypeOf(int64(42)):    "int64",
	reflect.TypeOf("42"):         "string",
}

var c Cache

func TestMain(m *testing.M) {
	c = &RedisCache{}
	err := c.Dial("172.17.0.26:6379")
	if err != nil {
		panic("Can't connect to redis test instance!")
	}
	c.SetSerializer(JsonSerialzer{})
	c.SetSeparator(":")
	c.SetTypePrefixRegistry(TypePrefixRegistry)

	c.SelectDb(13)
	c.FlushDb()

	exit := m.Run()

	//c.FlushDb()

	os.Exit(exit)
}

func Test_hellow(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func Test_set_get_string(t *testing.T) {
	o_s := "teststring"
	c.Set("0", o_s)
	var s string
	err := c.Get(&s, "0")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, o_s, s)

}

func Test_set_flush(t *testing.T) {
	o_s := "teststring"
	c.Set("1", o_s)

	c.FlushDb()

	var s string
	err := c.Get(&s, "1")
	assert.NotNil(t, err)

}

func Test_set_get_struct(t *testing.T) {
	o_s := TestStruct{}
	o_s.TestString = "hi"
	o_s.TestInt64 = 123429013233

	c.Set("0", o_s)

	s := TestStruct{}

	err := c.Get(&s, "0")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, o_s, s)

}

func Test_append_get_int64(t *testing.T) {
	var i int64
	o_l := make([]int64, 10)
	for i = 0; i < 10; i++ {
		c.Append("0", i)
		o_l[i] = i
	}
	var l []int64
	err := c.GetList(&l, "0")
	if err != nil {
		panic(err)
	}
	assert.Equal(t, o_l, l)

}

func Test_append_get_struct(t *testing.T) {
	o_s := TestStruct{}
	o_s.TestString = "hi"
	o_s.TestInt64 = 123429013233

	o_l := []TestStruct{o_s, o_s}

	c.Append("0", o_s)
	c.Append("0", o_s)

	l := []TestStruct{}

	err := c.GetList(&l, "0")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, o_l, l)

}

func Test_zadd_zrange_int64(t *testing.T) {
	var i int64
	o_s := []string{"four", "three", "two"}
	l := []string{"one", "two", "three", "four"}
	for i = 0; int(i) < len(l); i++ {
		c.ZAdd("testset", i+1, l[i])
	}

	s := []string{}
	scores, err := c.ZRangeByScore(&s, "testset", HIGH_TO_LOW, 3)
	if err != nil {
		panic(err)
	}

	assert.Equal(t, o_s, s)
	assert.Equal(t, []int64{4, 3, 2}, scores)
}
