package entities

import (
	"testing"
	"time"
)

func TestSetDateOfBirth(t *testing.T) {
	u := User{
		FirstName: "tester",
	}
	c := Customer{
		User: u,
	}
	err := c.SetDateOfBirth(time.Time{}, &c)
	if err != nil {
		t.Error(err)
	}
}
