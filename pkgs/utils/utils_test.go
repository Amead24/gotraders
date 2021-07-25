package utils

import "testing"

func TestGetServerHealth(t *testing.T) {
	resp, err := GetServerHealth()
	if resp != "200 OK" {
		t.Errorf("Response != 200, resp == %s with error %s", resp, err)
	}
}
