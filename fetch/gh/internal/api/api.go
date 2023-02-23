package api

import (
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/cli/go-gh"
)

func Run(path string) []byte {
	args := []string{"api", path, "--paginate"}
	log.Printf("Running: gh %s\n", strings.Join(args, " "))
	oBuf := runWithRetry(&retryManager{}, args)
	data, err := ioutil.ReadAll(&oBuf)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func runWithRetry(m *retryManager, args []string) bytes.Buffer {
	oBuf, _, err := gh.Exec(args...)
	if err == nil {
		return oBuf
	}
	retry, backoff := m.retryAfter(err)
	if !retry {
		log.Fatalf("FAILED: %s. No retry supported with this error.", err)
	}
	log.Printf("WARNING: %s. Retrying after %s", err, backoff)
	time.Sleep(backoff)
	return runWithRetry(m, args)
}
