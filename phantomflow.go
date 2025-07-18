package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Event struct {
	EventName        string  `json:"event_name"`
	Timestamp        string  `json:"timestamp"`
	ClientID         string  `json:"client_id"`
	UserID           *string `json:"user_id"`
	SessionID        string  `json:"session_id"`
	DeviceType       string  `json:"device_type"`
	TrafficMedium    string  `json:"traffic_medium"`
	TrafficSource    string  `json:"traffic_source"`
	TrafficCampaign  *string `json:"traffic_campaign"`
	PagePath         string  `json:"page_path"`
	PageTitle        string  `json:"page_title"`
	Referrer         string  `json:"referrer"`
	ClickX           *uint16 `json:"click_x"`
	ClickY           *uint16 `json:"click_y"`
	PageX            *uint16 `json:"page_x"`
	PageY            *uint16 `json:"page_y"`
	TargetText       string  `json:"target_text"`
	TargetTag        string  `json:"target_tag"`
	TargetClass      string  `json:"target_class"`
	TargetID         string  `json:"target_id"`
	TargetHref       string  `json:"target_href"`
	TargetType       string  `json:"target_type"`
	TargetValue      string  `json:"target_value"`
	IsButton         uint8   `json:"is_button"`
	IsLink           uint8   `json:"is_link"`
	IsInput          uint8   `json:"is_input"`
	IsTextarea       uint8   `json:"is_textarea"`
	IsSelect         uint8   `json:"is_select"`
	ElementTag       string  `json:"element_tag"`
	ElementID        string  `json:"element_id"`
	ElementClass     string  `json:"element_class"`
	ElementText      string  `json:"element_text"`
	ElementPath      string  `json:"element_path"`
	Country          string  `json:"country"`
	City             string  `json:"city"`
	Timezone         string  `json:"timezone"`
	DeviceOS         string  `json:"device_os"`
	Browser          string  `json:"browser"`
	Language         string  `json:"language"`
	UserAgent        string  `json:"user_agent"`
	ScreenResolution string  `json:"screen_resolution"`
	ViewportSize     string  `json:"viewport_size"`
	UtmParams        string  `json:"utm_params"`
	FormAction       string  `json:"form_action"`
	FormMethod       string  `json:"form_method"`
	FormID           string  `json:"form_id"`
	FormClass        string  `json:"form_class"`
	FormFields       string  `json:"form_fields"`
	FormFieldCount   uint8   `json:"form_field_count"`
	UserGender       *string `json:"user_gender"`
	UserAge          *uint8  `json:"user_age"`
	TimeOnPageSec    uint32  `json:"time_on_page_seconds"`
	SDKKey           string  `json:"sdk_key"`
}

var (
	successCount uint64
	failCount    uint64
	sbPool       = sync.Pool{
		New: func() interface{} {
			return &strings.Builder{}
		},
	}
)

func randomEvent(rng *rand.Rand) Event {
	currentTime := time.Now()
	nowStr := currentTime.Format("2006-01-02 15:04:05.000")

	sb := sbPool.Get().(*strings.Builder)
	defer sbPool.Put(sb)
	sb.Reset()

	sb.WriteString(strconv.Itoa(rng.Intn(1000000)))
	uid := sb.String()
	sb.Reset()

	fmt.Fprintf(sb, "%08x-%04x-%04x-%04x-%012x", rng.Uint32(), rng.Uint32()&0xffff, rng.Uint32()&0xffff, rng.Uint32()&0xffff, rng.Uint64()&0xffffffffffff)
	clientID := sb.String()
	sb.Reset()

	sb.WriteString("sess_")
	sb.WriteString(strconv.FormatInt(currentTime.UnixMilli(), 10))
	sb.WriteString("_")
	sb.WriteString(clientID[:6])
	sessionID := sb.String()
	sb.Reset()

	sb.WriteString("button ")
	sb.WriteString(strconv.Itoa(rng.Intn(10)))
	targetText := sb.String()
	sb.Reset()

	pagePaths := []string{"/", "/home", "/product", "/cart", "/checkout"}
	deviceTypes := []string{"desktop", "mobile", "tablet"}
	browsers := []string{"Chrome", "Safari", "Edge", "Firefox"}
	countries := []string{"KR", "US", "JP", "CN"}
	cities := []string{"Seoul", "Busan", "Tokyo", "Beijing"}
	langs := []string{"ko-KR", "en-US", "ja-JP"}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"Mozilla/5.0 (Linux; Android 10)",
	}

	var userAge *uint8
	if rng.Float32() < 0.7 {
		age := uint8(rng.Intn(60) + 18)
		userAge = &age
	}
	var userGender *string
	if rng.Float32() < 0.5 {
		g := []string{"M", "F"}[rng.Intn(2)]
		userGender = &g
	}

	return Event{
		EventName:        "auto_click",
		Timestamp:        nowStr,
		ClientID:         clientID,
		UserID:           &uid,
		SessionID:        sessionID,
		DeviceType:       deviceTypes[rng.Intn(len(deviceTypes))],
		TrafficMedium:    "direct",
		TrafficSource:    "cli_simulator",
		TrafficCampaign:  nil,
		PagePath:         pagePaths[rng.Intn(len(pagePaths))],
		PageTitle:        "CLI Simulate",
		Referrer:         "",
		ClickX:           nil,
		ClickY:           nil,
		PageX:            nil,
		PageY:            nil,
		TargetText:       targetText,
		TargetTag:        "button",
		TargetClass:      "",
		TargetID:         "",
		TargetHref:       "",
		TargetType:       "",
		TargetValue:      "",
		IsButton:         1,
		IsLink:           0,
		IsInput:          0,
		IsTextarea:       0,
		IsSelect:         0,
		ElementTag:       "button",
		ElementID:        "",
		ElementClass:     "",
		ElementText:      "",
		ElementPath:      "",
		Country:          countries[rng.Intn(len(countries))],
		City:             cities[rng.Intn(len(cities))],
		Timezone:         "Asia/Seoul",
		DeviceOS:         "Linux",
		Browser:          browsers[rng.Intn(len(browsers))],
		Language:         langs[rng.Intn(len(langs))],
		UserAgent:        userAgents[rng.Intn(len(userAgents))],
		ScreenResolution: "1920x1080",
		ViewportSize:     "1920x900",
		UtmParams:        "",
		FormAction:       "",
		FormMethod:       "",
		FormID:           "",
		FormClass:        "",
		FormFields:       "",
		FormFieldCount:   0,
		UserGender:       userGender,
		UserAge:          userAge,
		TimeOnPageSec:    uint32(rng.Intn(300)),
		SDKKey:           "test_sdk_key",
	}
}

func worker(client *http.Client, endpoint string, requests int, duration time.Duration, wg *sync.WaitGroup, seed int64) {
	defer wg.Done()
	rng := rand.New(rand.NewSource(seed))
	start := time.Now()
	var sent int
	for {
		if requests > 0 && sent >= requests {
			break
		}
		if duration > 0 && time.Since(start) > duration {
			break
		}
		event := randomEvent(rng)
		body, _ := json.Marshal(event)
		req, _ := http.NewRequest("POST", endpoint, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", event.UserAgent)
		resp, err := client.Do(req)
		if err != nil {
			atomic.AddUint64(&failCount, 1)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == 200 {
			atomic.AddUint64(&successCount, 1)
		} else {
			atomic.AddUint64(&failCount, 1)
		}
		sent++
	}
}

func main() {
	go func() {
		_ = http.ListenAndServe("localhost:6060", nil)
	}()

	var (
		endpoint string
		totalReq int
		duration int
		workers  int
	)
	flag.IntVar(&totalReq, "n", 0, "Total number of requests (0 for duration-based)")
	flag.IntVar(&duration, "d", 0, "Duration in seconds (0 for request-based)")
	flag.IntVar(&workers, "c", runtime.NumCPU()*2, "Number of concurrent workers (default: CPU cores * 2)")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: phantomflow [ENDPOINT] [-n total_requests] [-d duration_seconds] [-c concurrent_workers]")
		os.Exit(1)
	}
	endpoint = flag.Arg(0)
	fmt.Printf("Target: %s\nWorkers: %d\n", endpoint, workers)
	if totalReq > 0 {
		fmt.Printf("Total Requests: %d\n", totalReq)
	} else {
		fmt.Printf("Duration: %d seconds\n", duration)
	}

	kst, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		log.Fatalf("Could not load KST location: %v", err)
	}

	done := make(chan struct{})
	go func() {
		var lastSuccess, lastFail uint64
		var lastTime = time.Now()
		for {
			select {
			case <-done:
				return
			case <-time.After(1 * time.Second):
				now := time.Now()
				s := atomic.LoadUint64(&successCount)
				f := atomic.LoadUint64(&failCount)
				dt := now.Sub(lastTime).Seconds()
				if dt == 0 {
					continue
				}
				rps := float64((s+f)-(lastSuccess+lastFail)) / dt
				fmt.Printf("%s [Progress] Success: %d, Fail: %d, RPS: %.2f\n", now.In(kst).Format("2006-01-02 15:04:05"), s, f, rps)
				lastSuccess, lastFail = s, f
				lastTime = now
			}
		}
	}()

	tr := &http.Transport{
		MaxIdleConns:        4096,
		MaxIdleConnsPerHost: 4096,
		MaxConnsPerHost:     4096,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   3 * time.Second,
			KeepAlive: 90 * time.Second,
		}).DialContext,
		DisableCompression: false,
		ForceAttemptHTTP2:  false,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}
	var wg sync.WaitGroup
	start := time.Now()
	reqsPerWorker := 0
	if totalReq > 0 {
		reqsPerWorker = totalReq / workers
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(client, endpoint, reqsPerWorker, time.Duration(duration)*time.Second, &wg, int64(i)*1000+time.Now().UnixNano())
	}
	wg.Wait()
	close(done)

	elapsed := time.Since(start).Seconds()
	total := atomic.LoadUint64(&successCount) + atomic.LoadUint64(&failCount)
	fmt.Printf("\nTotal: %d, Success: %d, Fail: %d\n", total, successCount, failCount)
	if elapsed > 0 {
		fmt.Printf("Average RPS: %.2f\n", float64(total)/elapsed)
	}
}
