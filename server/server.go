package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// ServerConfig는 각 포트와 해당 포트에서 응답할 텍스트를 정의하는 구조체입니다.
type ServerConfig struct {
	Port int
	Text string
}

// LogEntry는 각 클라이언트의 접속 기록을 나타내는 구조체입니다.
type LogEntry struct {
	Timestamp  time.Time
	ClientIP   string
	ClientPort string
	ServerPort int
}

// ReportEntry는 리포트 출력을 위한 구조체로, 최종 접속 시간, 시도 횟수, 접속 IP와 포트를 포함합니다.
type ReportEntry struct {
	ClientIP    string
	ClientPort  string
	LastAccess  time.Time
	TryCount    int
	ServerPorts map[int]bool
}

// printServerHelp는 CSV 파일 형식에 대한 설명과 사용법을 출력하는 함수입니다.
func printServerHelp() {
	fmt.Println("Usage:")
	fmt.Println("  --start | -s         Start the server with the specified configuration")
	fmt.Println("  --report | -r        Generate a connection report based on access logs")
	fmt.Println("  --config | -c <file> Path to server configuration CSV file (default: servers.csv)")
	fmt.Println("\nExpected CSV Format for --config file:")
	fmt.Println("  servers.csv:")
	fmt.Println("    Port,Text")
	fmt.Println("    8080,Hello from port 8080")
	fmt.Println("    8081,Welcome to port 8081")
}

// validateServerCSV는 입력된 서버 설정 파일의 형식을 검증합니다.
func validateServerCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open server config file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read server config file: %v", err)
	}

	for i, record := range records {
		if len(record) < 2 {
			return fmt.Errorf("invalid format at line %d: expected [Port, Text]", i+1)
		}
		if _, err := strconv.Atoi(record[0]); err != nil {
			return fmt.Errorf("invalid port at line %d: %v", i+1, err)
		}
	}
	return nil
}

// readCSV는 서버 설정 파일을 읽어 ServerConfig 구조체로 반환합니다.
func readCSV(filePath string) ([]ServerConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var configs []ServerConfig
	for _, record := range records {
		port, _ := strconv.Atoi(record[0])
		text := record[1]
		configs = append(configs, ServerConfig{Port: port, Text: text})
	}
	return configs, nil
}

// logConnection은 클라이언트 접속을 기록하는 함수입니다.
func logConnection(clientIP, clientPort string, serverPort int) {
	file, err := os.OpenFile("access_log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	record := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		clientIP,
		clientPort,
		strconv.Itoa(serverPort),
	}
	if err := writer.Write(record); err != nil {
		log.Printf("Error writing log record: %v", err)
	}
}

// startServer는 각 포트에 대한 서버를 설정하고 접속 기록을 남기는 함수입니다.
func startServer(config ServerConfig, wg *sync.WaitGroup) {
	defer wg.Done()
	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientAddr := r.RemoteAddr
		clientIP, clientPort, _ := net.SplitHostPort(clientAddr)
		logConnection(clientIP, clientPort, config.Port)
		fmt.Fprintf(w, config.Text)
	})

	address := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Starting server on port %d\n", config.Port)
	if err := http.ListenAndServe(address, server); err != nil {
		log.Fatalf("Failed to start server on port %d: %v", config.Port, err)
	}
}

// generateReport는 access_log.csv 파일을 기반으로 리포트를 생성하고 report.csv에 저장합니다.
func generateReport(logFilePath, reportFilePath string) error {
	logs, err := readLog(logFilePath)
	if err != nil {
		return fmt.Errorf("Error reading log file: %v", err)
	}

	report := buildReport(logs)
	if err := saveReportToCSV(report, reportFilePath); err != nil {
		return fmt.Errorf("Error saving report to CSV: %v", err)
	}

	fmt.Printf("Report saved to %s\n", reportFilePath)
	return nil
}

// readLog는 access_log.csv 파일을 읽고 LogEntry 리스트로 반환합니다.
func readLog(filePath string) ([]LogEntry, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var logs []LogEntry
	for _, record := range records {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", record[0])
		serverPort, _ := strconv.Atoi(record[3])

		logs = append(logs, LogEntry{
			Timestamp:  timestamp,
			ClientIP:   record[1],
			ClientPort: record[2],
			ServerPort: serverPort,
		})
	}
	return logs, nil
}

// buildReport는 LogEntry 리스트를 바탕으로 리포트를 생성합니다.
func buildReport(logs []LogEntry) []ReportEntry {
	reportMap := make(map[string]*ReportEntry)

	for _, log := range logs {
		key := log.ClientIP + ":" + log.ClientPort
		entry, exists := reportMap[key]
		if !exists {
			entry = &ReportEntry{
				ClientIP:    log.ClientIP,
				ClientPort:  log.ClientPort,
				LastAccess:  log.Timestamp,
				TryCount:    0,
				ServerPorts: make(map[int]bool),
			}
			reportMap[key] = entry
		}
		entry.TryCount++
		entry.ServerPorts[log.ServerPort] = true
		if log.Timestamp.After(entry.LastAccess) {
			entry.LastAccess = log.Timestamp
		}
	}

	var report []ReportEntry
	for _, entry := range reportMap {
		report = append(report, *entry)
	}
	return report
}

// saveReportToCSV는 ReportEntry 리스트를 report.csv 파일에 저장합니다.
func saveReportToCSV(report []ReportEntry, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Client IP", "Client Port", "Last Access", "Tries", "Server Ports"})

	for _, entry := range report {
		ports := make([]int, 0, len(entry.ServerPorts))
		for port := range entry.ServerPorts {
			ports = append(ports, port)
		}
		fmtPorts := fmt.Sprint(ports)
		writer.Write([]string{
			entry.ClientIP,
			entry.ClientPort,
			entry.LastAccess.Format("2006-01-02 15:04:05"),
			strconv.Itoa(entry.TryCount),
			fmtPorts,
		})
	}
	return nil
}

func main() {
	start := flag.Bool("start", false, "Start the server (short: -s)")
	report := flag.Bool("report", false, "Generate connection report (short: -r)")
	config := flag.String("config", "servers.csv", "Path to server configuration file (short: -c)")
	flag.BoolVar(start, "s", false, "Start the server")
	flag.BoolVar(report, "r", false, "Generate connection report")
	flag.StringVar(config, "c", "servers.csv", "Path to server configuration file")
	flag.Parse()

	if len(os.Args) < 2 {
		printServerHelp()
		return
	}

	if *start {
		if err := validateServerCSV(*config); err != nil {
			log.Fatalf("Configuration file validation failed: %v", err)
		}
		configs, err := readCSV(*config)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		var wg sync.WaitGroup
		for _, config := range configs {
			wg.Add(1)
			go startServer(config, &wg)
		}
		wg.Wait()
	} else if *report {
		if err := generateReport("access_log.csv", "report.csv"); err != nil {
			log.Fatalf("Error generating report: %v", err)
		}
	} else {
		printServerHelp()
	}
}
