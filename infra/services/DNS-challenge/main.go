package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const DefaultPort = "8080"

var (
	BegetLogin    = os.Getenv("BEGET_LOGIN")
	BegetPassword = os.Getenv("BEGET_PASSWORD")
	HTTPTimeout   = 30 * time.Second
	HTTPClient    = &http.Client{Timeout: HTTPTimeout}
	Logger        = log.New(os.Stdout, "", 0)
)

const (
	Reset = "\033[0m"
	Red   = "\033[31m"
	Green = "\033[32m"
	Cyan  = "\033[36m"
)

var location = func() *time.Location {
	if tz := os.Getenv("TZ"); tz != "" {
		if loc, err := time.LoadLocation(tz); err == nil {
			return loc
		}
	}
	return time.UTC
}()

type DNSRequest struct {
	FQDN  string `json:"fqdn"`
	Value string `json:"value"`
}

func main() {
	if BegetLogin == "" || BegetPassword == "" {
		logFatalf("BEGET_LOGIN and BEGET_PASSWORD required")
	}

	http.HandleFunc("/healthz", handleHealthz)
	http.HandleFunc("/present", handlePresent)
	http.HandleFunc("/cleanup", handleCleanup)

	logInfof("Server started on port=%s", DefaultPort)
	if err := http.ListenAndServe(":"+DefaultPort, nil); err != nil {
		logFatalf("Server failed error=%v", err)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func handlePresent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DNSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logWarnf("Invalid JSON body")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.FQDN == "" || req.Value == "" {
		logWarnf("Missing required fields fqdn=%s value=%s", req.FQDN, req.Value)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	target := strings.TrimSuffix(req.FQDN, ".")

	logInfof("Setting TXT record target=%s value=%s", target, req.Value)
	if err := setTXTRecord(r.Context(), target, req.Value); err != nil {
		logErrorf("Failed to set TXT record error=%s target=%s", err.Error(), target)
		http.Error(w, "Set failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleCleanup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DNSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logWarnf("Invalid JSON body")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.FQDN == "" {
		logWarnf("Missing required field fqdn=%s", req.FQDN)
		http.Error(w, "Missing required field", http.StatusBadRequest)
		return
	}

	target := strings.TrimSuffix(req.FQDN, ".")

	logInfof("Clearing TXT record target=%s", target)
	if err := clearTXTRecord(r.Context(), target); err != nil {
		logErrorf("Failed to clear TXT record error=%s target=%s", err.Error(), target)
		http.Error(w, "Clear failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func setTXTRecord(ctx context.Context, domain, value string) error {
	params := map[string]interface{}{
		"fqdn": domain,
		"records": map[string]interface{}{
			"TXT": []interface{}{map[string]interface{}{"value": value, "priority": 0}},
		},
	}
	_, err := callBegetAPI(ctx, "changeRecords", params)
	return err
}

func clearTXTRecord(ctx context.Context, domain string) error {
	params := map[string]interface{}{
		"fqdn": domain,
		"records": map[string]interface{}{
			"TXT": []interface{}{},
		},
	}
	_, err := callBegetAPI(ctx, "changeRecords", params)
	return err
}

func callBegetAPI(ctx context.Context, method string, params interface{}) (interface{}, error) {
	jsonData, _ := json.Marshal(params)
	values := url.Values{}
	values.Set("login", BegetLogin)
	values.Set("passwd", BegetPassword)
	values.Set("input_format", "json")
	values.Set("output_format", "json")
	values.Set("input_data", string(jsonData))

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		"https://api.beget.com/api/dns/"+method+"?"+values.Encode(), nil)
	req.Header.Set("User-Agent", "Beget-DNS-ACME-Hook/1.0")

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	var response interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	if m, ok := response.(map[string]interface{}); ok {
		if s, ok := m["status"].(string); ok && s == "error" {
			return nil, fmt.Errorf("API error")
		}
	}
	return response, nil
}

func nowFormatted() string {
	return time.Now().In(location).Format("2006-01-02T15:04:05-07:00")
}

func logInfof(format string, args ...interface{}) {
	msg := colorizeKeys(fmt.Sprintf(format, args...))
	Logger.Printf("%s %s%s%s %s", nowFormatted(), Green, "INF", Reset, msg)
}

func logWarnf(format string, args ...interface{}) {
	msg := colorizeKeys(fmt.Sprintf(format, args...))
	Logger.Printf("%s %s%s%s %s", nowFormatted(), "\033[33m", "WRN", Reset, msg)
}

func logErrorf(format string, args ...interface{}) {
	msg := colorizeKeys(fmt.Sprintf(format, args...))
	Logger.Printf("%s %s%s%s %s", nowFormatted(), Red, "ERR", Reset, msg)
}

func logFatalf(format string, args ...interface{}) {
	msg := colorizeKeys(fmt.Sprintf(format, args...))
	Logger.Printf("%s %s%s%s %s", nowFormatted(), Red, "FATAL", Reset, msg)
	os.Exit(1)
}

func colorizeKeys(msg string) string {
	parts := strings.Fields(msg)
	var result []string

	for _, part := range parts {
		if i := strings.Index(part, "="); i > 0 {
			key := part[:i+1]
			val := part[i+1:]
			result = append(result, Cyan+key+Reset+val)
		} else {
			result = append(result, part)
		}
	}

	return strings.Join(result, " ")
}
