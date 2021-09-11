package user

import (
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestBlockedUserNameURL(t *testing.T) {
	url := createBlockURL("1")
	AssertEqual(t, url, "https://api.twitter.com/2/users/1/blocking")
}

func TestDeleteBlockedUserNamesURL(t *testing.T) {
	url := createDeleteBlockURL("1", "2")
	AssertEqual(t, url, "https://api.twitter.com/2/users/1/blocking/2")
}

func TestBlockUserId(t *testing.T) {
	response := BlockUserId("1")
	AssertEqual(t, response, nil)
}

func TestDeleteBlock(t *testing.T, id string, deleteId string) {
	DeleteBlock(id, deleteId)
}

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}
