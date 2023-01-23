package todo

import (
	"testing"
)

func TestCreateTodoNotAllowSleepTask(t *testing.T) {

	handler := NewTodoHandler(&TestDB{})
	c := &TestContext{}

	handler.NewTask(c)

	want := "not allowed"

	if got := c.v["error"]; got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}

type TestDB struct{}

func (TestDB) New(*Todo) error {
	return nil
}
func (TestDB) Find(*[]Todo) error {
	return nil
}
func (TestDB) Delete(*Todo, int) error {
	return nil
}

type TestContext struct {
	v map[string]interface{}
}

func (TestContext) Bind(v interface{}) error {
	*v.(*Todo) = Todo{
		Title: "sleep",
	}
	return nil
}

func (c *TestContext) JSON(statuscode int, v interface{}) {
	c.v = v.(map[string]interface{})
}

func (TestContext) TransactionID() string {
	return "TestTransactionID"
}

func (TestContext) Audience() string {
	return "Unit Test"
}

func (TestContext) Param(v string) string {
	return "Param"
}
