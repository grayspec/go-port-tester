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

// TestResult 구조체는 각 서버와 포트 테스트 결과를 저장합니다.
type TestResult struct {
	IP     string
	Port   int
	Status string
}

// ServerConfig 구조체는 CSV 파일에서 읽어들인 각 서버와 포트 정보를 저장합니다.
type ServerConfig struct {
	IP   string
	Port int
}

// printClientHelp는 CSV 파일 형식에 대한 설명과 사용법을 출력합니다.
func printClientHelp() {
	fmt.Println("Usage:")
	fmt.Println("  --servers | -s <file>      Path to servers and ports CSV file (default: servers.csv)")
	fmt.Println("  --output | -o <file>       Path to output result CSV file (default: result.csv)")
	fmt.Println("  --timeout | -t <seconds>   Timeout for each request (default: 2 seconds)")
	fmt.Println("  --concurrency | -c <int>   Number of concurrent requests (default: 5)")
	fmt.Println("\nExpected CSV Format for servers and ports file:")
	fmt.Println("  servers.csv:")
	fmt.Println("    IP,Port")
	fmt.Println("    127.0.0.1,8080")
	fmt.Println("    192.168.1.10,8081")
}

// validateServersCSV는 서버와 포트 정보가 포함된 CSV 파일의 형식을 검증합니다.
func validateServersCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open servers file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read servers file: %v", err)
	}

	for i, record := range records {
		if len(record) < 2 {
			return fmt.Errorf("invalid format at line %d: expected [IP, Port]", i+1)
		}
		if net.ParseIP(record[0]) == nil {
			return fmt.Errorf("invalid IP address at line %d: %s", i+1, record[0])
		}
		if _, err := strconv.Atoi(record[1]); err != nil {
			return fmt.Errorf("invalid port at line %d: %v", i+1, err)
		}
	}
	return nil
}

// readServersCSV는 서버와 포트 정보를 CSV 파일에서 읽어 ServerConfig 리스트로 반환합니다.
func readServersCSV(filePath string) ([]ServerConfig, error) {
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

	var servers []ServerConfig
	for _, record := range records {
		port, _ := strconv.Atoi(record[1])
		servers = append(servers, ServerConfig{
			IP:   record[0],
			Port: port,
		})
	}
	return servers, nil
}

// checkPort는 지정된 IP와 포트에 대해 HTTP 요청을 보내 포트가 열려 있는지 확인합니다.
func checkPort(ip string, port int, timeout time.Duration, results chan<- TestResult) {
	url := fmt.Sprintf("http://%s:%d", ip, port)
	client := http.Client{Timeout: timeout}
	resp, err := client.Get(url)

	status := "closed"
	if err == nil && resp.StatusCode == http.StatusOK {
		status = "open"
	}

	results <- TestResult{IP: ip, Port: port, Status: status}
	if err == nil {
		resp.Body.Close()
	}
}

// writeResults는 포트 테스트 결과를 CSV 파일로 저장합니다.
func writeResults(results []TestResult, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"IP", "Port", "Status"})
	for _, result := range results {
		writer.Write([]string{result.IP, strconv.Itoa(result.Port), result.Status})
	}
	return nil
}

// runClient는 클라이언트의 메인 로직으로, 모든 IP와 포트에 대해 테스트를 수행하고 결과를 출력합니다.
func runClient(serversFile, outputFile string, timeout time.Duration, concurrency int) {
	servers, err := readServersCSV(serversFile)
	if err != nil {
		log.Fatalf("Error reading servers file: %v", err)
	}

	var wg sync.WaitGroup
	results := make(chan TestResult, len(servers))
	semaphore := make(chan struct{}, concurrency) // 동시 실행 제한

	totalTests := len(servers)
	successCount := 0
	failureCount := 0

	// 각 서버와 포트에 대해 상태 확인
	for _, server := range servers {
		wg.Add(1)
		semaphore <- struct{}{} // 동시 실행 개수 제한

		go func(ip string, port int) {
			defer wg.Done()
			checkPort(ip, port, timeout, results)
			<-semaphore
		}(server.IP, server.Port)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var finalResults []TestResult
	for result := range results {
		finalResults = append(finalResults, result)
		if result.Status == "open" {
			successCount++
		} else {
			failureCount++
		}
	}

	// 결과 CSV 파일 작성
	if err := writeResults(finalResults, outputFile); err != nil {
		log.Fatalf("Error writing results to CSV: %v", err)
	}

	// 통계 로그 출력
	fmt.Printf("Total Tests: %d, Success: %d, Failure: %d\n", totalTests, successCount, failureCount)
}

func main() {
	// CLI 플래그 설정
	serversFile := flag.String("servers", "servers.csv", "Path to servers and ports file (short: -s)")
	outputFile := flag.String("output", "result.csv", "Path to output result file (short: -o)")
	timeout := flag.Int("timeout", 2, "Timeout in seconds (short: -t)")
	concurrency := flag.Int("concurrency", 5, "Number of concurrent requests (short: -c)")

	// Short 옵션 설정
	flag.StringVar(serversFile, "s", "servers.csv", "Path to servers and ports file")
	flag.StringVar(outputFile, "o", "result.csv", "Path to output result file")
	flag.IntVar(timeout, "t", 2, "Timeout in seconds")
	flag.IntVar(concurrency, "c", 5, "Number of concurrent requests")
	flag.Parse()

	if len(os.Args) < 2 {
		printClientHelp()
		return
	}

	// 서버와 포트 정보를 포함한 CSV 파일 검증
	if err := validateServersCSV(*serversFile); err != nil {
		log.Fatalf("Servers file validation failed: %v", err)
	}

	timeoutDuration := time.Duration(*timeout) * time.Second
	fmt.Printf("Starting port test with %d second timeout and %d concurrent requests\n", *timeout, *concurrency)
	runClient(*serversFile, *outputFile, timeoutDuration, *concurrency)
}
