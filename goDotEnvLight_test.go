package goDotEnvLight

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestExport(t *testing.T) {

	envs := make(map[string]string)
	fmt.Println("\n== Prev Envs. ==")
	prevRawEnvs := os.Environ()
	for _, prevRawEnvLine := range prevRawEnvs {
		prevRawEnvLineSplit := strings.Split(prevRawEnvLine, "=")
		envs[prevRawEnvLineSplit[0]] = prevRawEnvLineSplit[1]
		fmt.Println(prevRawEnvLine)
	}

	fmt.Println()

	succeeds, fails, err := Export(false, "./sample01.env")
	if err != nil {
		fmt.Printf("\nERROR: %s", err)
		return
	}
	for k, v := range succeeds {
		fmt.Printf("\nSUCCEED : %s=%s", k, v)
	}
	for k, v := range fails {
		fmt.Printf("\nFAIL : %s=%s", k, v)
	}

	fmt.Println("\n\n\n== Updated Envs. ==")
	updatedRawEnvs := os.Environ()
	for _, updatedRawEnvLine := range updatedRawEnvs {
		updatedRawEnvLineSplit := strings.Split(updatedRawEnvLine, "=")

		if val, ok := envs[updatedRawEnvLineSplit[0]]; ok {
			if val != updatedRawEnvLineSplit[1] {
				fmt.Print("*Overwrited* ")
			}
		} else {
			fmt.Print("*New* ")
		}
		fmt.Println(updatedRawEnvLine)
	}
}
