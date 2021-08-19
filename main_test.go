package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
)

func Test_main(t *testing.T) {
	subcmd := "./subcmd"
	buildHelper := exec.Command("go", "build", "-o", subcmd, "./test/subcmd/main.go")
	buildHelper.Stdout = os.Stdout
	buildHelper.Stderr = os.Stderr
	if err := buildHelper.Run(); err != nil {
		t.Errorf("unable to build test helper: %v", err)
	}
	defer os.Remove(subcmd)

	tests := []struct {
		args     string
		calls    int
		subCalls int
	}{
		//one parallel call - will result one call of sub command
		{"./subcmd --duration=1s ./test1", 1, 1},

		//two parallel calls - will result one call immediately and one call deferred
		{"./subcmd --duration=1s ./test2", 2, 2},

		//two parallel calls - will result one call immediately and one call deferred, third will be skipped as there is deferred one already
		{"./subcmd --duration=1s ./test3", 3, 2},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s:%d:%d", tt.args, tt.calls, tt.subCalls), func(t *testing.T) {
			args := strings.Split(tt.args, " ")
			testFile := args[len(args)-1]
			_ = os.Remove(testFile)
			defer os.Remove(testFile)
			os.Args = append([]string{""}, args...)
			wg := &sync.WaitGroup{}
			for i := 0; i < tt.calls; i++ {
				wg.Add(1)
				go func() {
					main()
					wg.Done()
				}()
			}
			wg.Wait()
			bytes, err := os.ReadFile(args[len(args)-1])
			if err != nil {
				t.Fatal(err)
			}

			realCalls, _ := strconv.Atoi(string(bytes))
			if realCalls != tt.subCalls {
				t.Fatalf("subcalls want %d, got %d", tt.subCalls, realCalls)
			}
		})
	}
}
