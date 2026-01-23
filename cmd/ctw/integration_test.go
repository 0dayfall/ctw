package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var ctwBinPath string

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration tests: failed to get working directory:", err)
		os.Exit(1)
	}

	repoRoot, err := findRepoRoot(wd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration tests:", err)
		os.Exit(1)
	}

	tempDir, err := os.MkdirTemp("", "ctw-integration-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "integration tests: failed to create temp dir:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempDir)

	ctwBinPath = filepath.Join(tempDir, "ctw")
	buildCmd := exec.Command("go", "build", "-o", ctwBinPath, "./cmd/ctw")
	buildCmd.Dir = repoRoot
	buildCmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Fprintln(os.Stderr, "integration tests: failed to build ctw:", err)
		fmt.Fprintln(os.Stderr, string(output))
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestSearchRecentGzipSuccess(t *testing.T) {
	errCh := make(chan error, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/2/tweets/search/recent" {
			recordError(errCh, fmt.Errorf("unexpected path: %s", r.URL.Path))
		}
		if r.URL.Query().Get("query") != "golang" {
			recordError(errCh, fmt.Errorf("unexpected query: %q", r.URL.Query().Get("query")))
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			recordError(errCh, fmt.Errorf("unexpected authorization header: %q", auth))
		}

		payload := `{"data":[{"id":"1","text":"hello"}],"meta":{"result_count":1}}`
		compressed := gzipPayload(t, payload)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compressed)
	}))
	defer server.Close()

	stdout, stderr, err := runCTW(t,
		"--base-url", server.URL,
		"--bearer-token", "test-token",
		"search", "recent",
		"--query", "golang",
	)
	if err != nil {
		t.Fatalf("expected success, got error: %v\nstderr: %s", err, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected no stderr output, got: %s", stderr)
	}

	var payload struct {
		Data []struct {
			Text string `json:"text"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("failed to decode stdout JSON: %v\nstdout: %s", err, stdout)
	}
	if len(payload.Data) != 1 || payload.Data[0].Text != "hello" {
		t.Fatalf("unexpected response payload: %+v", payload)
	}

	drainErrors(t, errCh)
}

func TestSearchRecentGzipError(t *testing.T) {
	errCh := make(chan error, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/2/tweets/search/recent" {
			recordError(errCh, fmt.Errorf("unexpected path: %s", r.URL.Path))
		}
		if r.URL.Query().Get("query") != "golang" {
			recordError(errCh, fmt.Errorf("unexpected query: %q", r.URL.Query().Get("query")))
		}

		payload := `{"errors":[{"code":89,"message":"Invalid or expired token"}]}`
		compressed := gzipPayload(t, payload)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write(compressed)
	}))
	defer server.Close()

	_, stderr, err := runCTW(t,
		"--base-url", server.URL,
		"--bearer-token", "test-token",
		"search", "recent",
		"--query", "golang",
	)
	if err == nil {
		t.Fatalf("expected error, got success")
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() == 0 {
		t.Fatalf("expected non-zero exit, got: %v", err)
	}
	if !strings.Contains(stderr, "twitter api error: status 401") {
		t.Fatalf("unexpected stderr output: %s", stderr)
	}
	if strings.Contains(stderr, "decode error response") {
		t.Fatalf("stderr contains gzip decode error: %s", stderr)
	}

	drainErrors(t, errCh)
}

func TestSearchRecentPipeJQ(t *testing.T) {
	if _, err := exec.LookPath("jq"); err != nil {
		t.Skip("jq not available")
	}

	errCh := make(chan error, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/2/tweets/search/recent" {
			recordError(errCh, fmt.Errorf("unexpected path: %s", r.URL.Path))
		}
		if r.URL.Query().Get("query") != "golang" {
			recordError(errCh, fmt.Errorf("unexpected query: %q", r.URL.Query().Get("query")))
		}

		payload := `{"data":[{"id":"1","text":"hello"},{"id":"2","text":"world"}],"meta":{"result_count":2}}`
		compressed := gzipPayload(t, payload)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compressed)
	}))
	defer server.Close()

	stdout, stderr, err := runCTW(t,
		"--base-url", server.URL,
		"--bearer-token", "test-token",
		"search", "recent",
		"--query", "golang",
	)
	if err != nil {
		t.Fatalf("expected success, got error: %v\nstderr: %s", err, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected no stderr output, got: %s", stderr)
	}

	jqOut, jqErr, err := runJQ(t, stdout, "-r", ".data[].text")
	if err != nil {
		t.Fatalf("jq failed: %v\nstderr: %s", err, jqErr)
	}
	if jqOut != "hello\nworld\n" {
		t.Fatalf("unexpected jq output: %q", jqOut)
	}

	drainErrors(t, errCh)
}

func TestUsersLookupGzipSuccess(t *testing.T) {
	errCh := make(chan error, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/2/users/by" {
			recordError(errCh, fmt.Errorf("unexpected path: %s", r.URL.Path))
		}
		if r.URL.Query().Get("usernames") != "twitter" {
			recordError(errCh, fmt.Errorf("unexpected usernames query: %q", r.URL.Query().Get("usernames")))
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			recordError(errCh, fmt.Errorf("unexpected authorization header: %q", auth))
		}

		payload := `{"data":[{"id":"42","username":"twitter","name":"Twitter"}]}`
		compressed := gzipPayload(t, payload)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compressed)
	}))
	defer server.Close()

	stdout, stderr, err := runCTW(t,
		"--base-url", server.URL,
		"--bearer-token", "test-token",
		"users", "lookup",
		"--usernames", "twitter",
	)
	if err != nil {
		t.Fatalf("expected success, got error: %v\nstderr: %s", err, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected no stderr output, got: %s", stderr)
	}

	var payload []struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("failed to decode stdout JSON: %v\nstdout: %s", err, stdout)
	}
	if len(payload) != 1 || payload[0].ID != "42" || payload[0].Username != "twitter" {
		t.Fatalf("unexpected response payload: %+v", payload)
	}

	drainErrors(t, errCh)
}

func TestTweetsGetGzipSuccess(t *testing.T) {
	errCh := make(chan error, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/2/tweets/123" {
			recordError(errCh, fmt.Errorf("unexpected path: %s", r.URL.Path))
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			recordError(errCh, fmt.Errorf("unexpected authorization header: %q", auth))
		}

		payload := `{"data":[{"id":"123","text":"hello"}]}`
		compressed := gzipPayload(t, payload)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compressed)
	}))
	defer server.Close()

	stdout, stderr, err := runCTW(t,
		"--base-url", server.URL,
		"--bearer-token", "test-token",
		"tweets", "get",
		"--id", "123",
	)
	if err != nil {
		t.Fatalf("expected success, got error: %v\nstderr: %s", err, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected no stderr output, got: %s", stderr)
	}

	var payload struct {
		Data []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("failed to decode stdout JSON: %v\nstdout: %s", err, stdout)
	}
	if len(payload.Data) != 1 || payload.Data[0].ID != "123" || payload.Data[0].Text != "hello" {
		t.Fatalf("unexpected response payload: %+v", payload)
	}

	drainErrors(t, errCh)
}

func TestTimelinesUserGzipSuccess(t *testing.T) {
	errCh := make(chan error, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/2/users/42/tweets" {
			recordError(errCh, fmt.Errorf("unexpected path: %s", r.URL.Path))
		}
		if r.URL.Query().Get("max_results") != "5" {
			recordError(errCh, fmt.Errorf("unexpected max_results query: %q", r.URL.Query().Get("max_results")))
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			recordError(errCh, fmt.Errorf("unexpected authorization header: %q", auth))
		}

		payload := `{"data":[{"id":"1","text":"hi"}],"meta":{"result_count":1}}`
		compressed := gzipPayload(t, payload)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(compressed)
	}))
	defer server.Close()

	stdout, stderr, err := runCTW(t,
		"--base-url", server.URL,
		"--bearer-token", "test-token",
		"timelines", "user",
		"--user-id", "42",
		"--param", "max_results=5",
	)
	if err != nil {
		t.Fatalf("expected success, got error: %v\nstderr: %s", err, stderr)
	}
	if strings.TrimSpace(stderr) != "" {
		t.Fatalf("expected no stderr output, got: %s", stderr)
	}

	var payload struct {
		Data []struct {
			ID   string `json:"id"`
			Text string `json:"text"`
		} `json:"data"`
		Meta struct {
			ResultCount int `json:"result_count"`
		} `json:"meta"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("failed to decode stdout JSON: %v\nstdout: %s", err, stdout)
	}
	if len(payload.Data) != 1 || payload.Data[0].Text != "hi" || payload.Meta.ResultCount != 1 {
		t.Fatalf("unexpected response payload: %+v", payload)
	}

	drainErrors(t, errCh)
}

func runCTW(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(ctwBinPath, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), "BEARER_TOKEN=")
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func runJQ(t *testing.T, input string, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command("jq", args...)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func gzipPayload(t *testing.T, payload string) []byte {
	t.Helper()
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	if _, err := gz.Write([]byte(payload)); err != nil {
		t.Fatalf("failed to gzip payload: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}
	return buffer.Bytes()
}

func recordError(errCh chan error, err error) {
	if err == nil {
		return
	}
	select {
	case errCh <- err:
	default:
	}
}

func drainErrors(t *testing.T, errCh chan error) {
	t.Helper()
	for {
		select {
		case err := <-errCh:
			if err != nil {
				t.Error(err)
			}
		default:
			return
		}
	}
}

func findRepoRoot(start string) (string, error) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", start)
		}
		dir = parent
	}
}
