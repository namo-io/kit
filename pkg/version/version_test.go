package version

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	fmt.Println(String())
	fmt.Println(Info().String())
}
