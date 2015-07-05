package cache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type TestStruct struct {
	TestString string `json:"test_string"`
	TestInt64  int64  `json:"test_int64"`
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

	c.SelectDb(13)
	c.FlushDb()

	exit := m.Run()

	c.FlushDb()

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
