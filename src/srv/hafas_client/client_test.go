package hafasclient

import (
	"fmt"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient()
	deps, err := client.GetDepartures("1109")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	fmt.Println(deps)
}

func TestClientNotExist(t *testing.T) {
	client := NewClient()
	deps, err := client.GetDepartures("12900929292")
	if err != nil {
		fmt.Println(deps)
		fmt.Println(err)
	} else {
		fmt.Println(deps)
		t.Errorf("Expected Error")
	}
}
