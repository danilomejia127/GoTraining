package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	mcpEndpoint      = "https://o11y-mcp.melioffice.com/mcp"
	application      = "payment-methods-read-v2"
	outputFile       = "CannotInferPM/data/clean_request.txt"
	testingBaseURL   = "https://testing-payment-methods.melioffice.com"
	apiCallLimitPerScope = 50
)

var availableScopes = []string{
	"production-reader",
	"production-reader-mla",
	"production-reader-mco",
	"production-reader-mlb",
	"production-reader-mlm",
}

// --- MCP types ---

type mcpRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type initParams struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ClientInfo      map[string]any `json:"clientInfo"`
}

type toolCallParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type mcpResponse struct {
	Result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"result"`
}

type logsResult struct {
	Status string `json:"status"`
	Data   struct {
		Logs []struct {
			Body string `json:"body"`
		} `json:"logs"`
		LogsCount int `json:"logs_count"`
		Returned  int `json:"returned"`
	} `json:"data"`
}

// --- Interactive input ---

func selectScopes(reader *bufio.Reader) []string {
	fmt.Println("\nScopes disponibles:")
	for i, s := range availableScopes {
		fmt.Printf("  [%d] %s\n", i+1, s)
	}
	fmt.Print("\nIngresá los números separados por coma (ej: 1,3,4) o 'all' para todos: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.ToLower(input) == "all" {
		return availableScopes
	}

	var selected []string
	for _, part := range strings.Split(input, ",") {
		part = strings.TrimSpace(part)
		idx, err := strconv.Atoi(part)
		if err != nil || idx < 1 || idx > len(availableScopes) {
			fmt.Printf("  Opción inválida ignorada: %s\n", part)
			continue
		}
		selected = append(selected, availableScopes[idx-1])
	}
	return selected
}

func selectTimeWindow(reader *bufio.Reader) (time.Time, time.Time) {
	fmt.Println("\nVentana de tiempo:")
	fmt.Println("  [1] Última 1 hora")
	fmt.Println("  [2] Últimas 4 horas (máximo permitido por la API)")
	fmt.Println("  [3] Personalizada (ingresá fecha/hora manualmente)")
	fmt.Print("Elegí una opción: ")

	opt, _ := reader.ReadString('\n')
	opt = strings.TrimSpace(opt)

	now := time.Now().UTC()
	switch opt {
	case "1":
		return now.Add(-1 * time.Hour), now
	case "2":
		return now.Add(-4 * time.Hour), now
	case "3":
		fmt.Print("  Fecha inicio (formato 2006-01-02T15:04:05Z): ")
		fromStr, _ := reader.ReadString('\n')
		fromStr = strings.TrimSpace(fromStr)
		fmt.Print("  Fecha fin   (formato 2006-01-02T15:04:05Z): ")
		toStr, _ := reader.ReadString('\n')
		toStr = strings.TrimSpace(toStr)

		from, err1 := time.Parse(time.RFC3339, fromStr)
		to, err2 := time.Parse(time.RFC3339, toStr)
		if err1 != nil || err2 != nil {
			fmt.Println("  Formato inválido, usando última 1 hora.")
			return now.Add(-1 * time.Hour), now
		}
		// La API admite máximo 4 horas por ventana
		if to.Sub(from) > 4*time.Hour {
			fmt.Println("  La ventana supera 4 horas (límite de la API). Se ajusta al máximo.")
			from = to.Add(-4 * time.Hour)
		}
		return from, to
	default:
		fmt.Println("  Opción inválida, usando última 1 hora.")
		return now.Add(-1 * time.Hour), now
	}
}

// --- Auth ---

func getFuryToken() (string, error) {
	cmd := exec.Command(
		"/Users/damejia/.local/pipx/venvs/mcp-remote-proxy/bin/python3",
		"-c",
		"from mcp_remote_proxy.furyauth import get_fury_auth_token; print(get_fury_auth_token())",
	)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error obteniendo token fury: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// --- MCP session ---

func initMCPSession(token string) (string, error) {
	body, _ := json.Marshal(mcpRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: initParams{
			ProtocolVersion: "2024-11-05",
			Capabilities:    map[string]any{},
			ClientInfo:      map[string]any{"name": "go-client", "version": "1.0"},
		},
	})

	req, _ := http.NewRequest("POST", mcpEndpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	sessionID := resp.Header.Get("mcp-session-id")
	if sessionID == "" {
		return "", fmt.Errorf("no se recibió mcp-session-id")
	}
	return sessionID, nil
}

func queryLogs(token, sessionID, scope, dateFrom, dateTo string) (*logsResult, error) {
	body, _ := json.Marshal(mcpRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/call",
		Params: toolCallParams{
			Name: "query_logs",
			Arguments: map[string]any{
				"application": application,
				"scope":       scope,
				"date_from":   dateFrom,
				"date_to":     dateTo,
				"query_name":  "all_logs",
				"size":        10000,
			},
		},
	})

	req, _ := http.NewRequest("POST", mcpEndpoint, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Set("mcp-session-id", sessionID)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var mcpResp mcpResponse
	if err := json.Unmarshal(respBody, &mcpResp); err != nil {
		return nil, fmt.Errorf("error parseando respuesta MCP: %v", err)
	}
	if len(mcpResp.Result.Content) == 0 {
		return nil, fmt.Errorf("respuesta MCP sin contenido")
	}

	var result logsResult
	if err := json.Unmarshal([]byte(mcpResp.Result.Content[0].Text), &result); err != nil {
		return nil, fmt.Errorf("error parseando logs: %v", err)
	}
	return &result, nil
}

// --- Extraction ---

func extractRequests(logs []struct{ Body string `json:"body"` }, seen map[string]bool) ([]string, map[string]int) {
	urlRe := regexp.MustCompile(`\[url:(/v1/payment_methods/search[^\]]+)\]`)
	apiRe := regexp.MustCompile(`\[api-name:([^\]]+)\]`)
	tagRe := regexp.MustCompile(`\[tag_cannot_infer_pm:([^\]]+)\]`)

	tagCounts := map[string]int{}
	var lines []string

	for _, l := range logs {
		if !strings.Contains(l.Body, "status_response_total:zero") {
			continue
		}
		if !strings.Contains(l.Body, "tag_cannot_infer_pm:excludes_by_rule") &&
			!strings.Contains(l.Body, "tag_cannot_infer_pm:not_result_by_params") {
			continue
		}

		tagM := tagRe.FindStringSubmatch(l.Body)
		if len(tagM) < 2 {
			continue
		}
		tag := tagM[1]

		if tagCounts[tag] >= apiCallLimitPerScope {
			continue
		}

		urlM := urlRe.FindStringSubmatch(l.Body)
		if len(urlM) < 2 {
			continue
		}
		url := urlM[1]

		if seen[url] {
			continue
		}
		seen[url] = true
		tagCounts[tag]++

		api := ""
		if apiM := apiRe.FindStringSubmatch(l.Body); len(apiM) >= 2 {
			api = apiM[1]
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\t%s", url, api, tag))
	}
	return lines, tagCounts
}

// --- API call ---

type pagingResponse struct {
	Paging struct {
		Total int `json:"total"`
	} `json:"paging"`
}

type requestRecord struct {
	Path        string
	Application string
	CannotInfer string
}

func parseRecords(lines []string) []requestRecord {
	var records []requestRecord
	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) < 3 {
			continue
		}
		records = append(records, requestRecord{Path: parts[0], Application: parts[1], CannotInfer: parts[2]})
	}
	return records
}

func siteFromPath(path string) string {
	for _, param := range strings.Split(strings.SplitN(path, "?", 2)[1], "&") {
		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 && kv[0] == "site_id" {
			return kv[1]
		}
	}
	return "unknown"
}

func callAPI(records []requestRecord) {
	client := &http.Client{Timeout: 30 * time.Second}
	fmt.Printf("\n--- Iniciando llamadas a la API (%d requests) ---\n\n", len(records))

	perSite := map[string]int{}

	for _, r := range records {
		fullURL := testingBaseURL + r.Path

		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			fmt.Printf("Error creando request: %v\n", err)
			continue
		}
		req.Header.Set("x-cannot-infer-response", "true")
		req.Header.Set("X-Api-Client-Application", r.Application)
		req.Header.Set("x-core-flow-type", "pm_offer")

		start := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			fmt.Printf("Error en request %s — %v\n", r.Path, err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%s — HTTP %d — Tiempo: %d ms\n", r.Path, resp.StatusCode, elapsed)
			continue
		}

		var result pagingResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("%s — Error parseando JSON — Tiempo: %d ms\n", r.Path, elapsed)
			continue
		}

		site := siteFromPath(r.Path)
		perSite[site]++
		fmt.Printf("[%s] %s  Total=%d  Tiempo: %d ms\n", testingBaseURL, r.Path, result.Paging.Total, elapsed)
	}

	fmt.Println("\n--- Resumen de requests por site ---")
	for site, count := range perSite {
		fmt.Printf("  %s: %d requests\n", site, count)
	}
	fmt.Printf("  Total: %d requests\n", len(records))
}

// --- Main ---

func main() {
	consoleReader := bufio.NewReader(os.Stdin)

	scopes := selectScopes(consoleReader)
	if len(scopes) == 0 {
		log.Fatal("No se seleccionó ningún scope.")
	}

	from, to := selectTimeWindow(consoleReader)
	dateFrom := from.Format(time.RFC3339)
	dateTo := to.Format(time.RFC3339)

	fmt.Printf("\nObteniendo token Fury...\n")
	token, err := getFuryToken()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		log.Fatalf("Error creando directorio: %v", err)
	}
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creando archivo: %v", err)
	}
	defer f.Close()
	writer := bufio.NewWriter(f)
	defer writer.Flush()

	var allRecords []requestRecord
	globalSeen := map[string]bool{}

	for _, scope := range scopes {
		fmt.Printf("\n[%s] Iniciando sesión MCP...\n", scope)
		sessionID, err := initMCPSession(token)
		if err != nil {
			fmt.Printf("[%s] Error sesión: %v\n", scope, err)
			continue
		}

		fmt.Printf("[%s] Consultando logs (%s → %s)...\n", scope, dateFrom, dateTo)
		result, err := queryLogs(token, sessionID, scope, dateFrom, dateTo)
		if err != nil {
			fmt.Printf("[%s] Error consultando: %v\n", scope, err)
			continue
		}

		lines, tagCounts := extractRequests(result.Data.Logs, globalSeen)
		records := parseRecords(lines)
		for _, line := range lines {
			fmt.Fprintln(writer, line)
		}
		fmt.Printf("[%s] Total logs: %d | excludes_by_rule: %d | not_result_by_params: %d | Total extraídos: %d\n",
			scope, result.Data.LogsCount, tagCounts["excludes_by_rule"], tagCounts["not_result_by_params"], len(records))
		allRecords = append(allRecords, records...)
	}

	fmt.Printf("\nArchivo guardado en: %s\n", outputFile)
	fmt.Printf("Total de requests únicos guardados: %d\n", len(allRecords))

	if len(allRecords) > 0 {
		callAPI(allRecords)
	}
}
