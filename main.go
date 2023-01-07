package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func sendToHorn(text string) {
	m := map[string]interface{}{
		"text": text,
	}
	mJson, _ := json.Marshal(m)
	contentReader := bytes.NewReader(mJson)
	req, err := http.NewRequest("POST", os.Getenv("INTEGRAM_WEBHOOK_URI"), contentReader)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Println(err)

		fmt.Println(resp)
	}
}

// checkConnections checks the connection to a list of IP addresses and sends a POST
// request with the IP address of any failed connection. If all connections fail, it
// restarts the ipsec service.
func checkConnections(ips []string, hostname string, probe, sleep int) {
	failCount := 0

	for {
		time.Sleep(time.Duration(sleep) * time.Second)

		// Check the connection to each IP address
		for _, ip := range ips {
			// Test the connection to the IP address using the ping command
			if _, err := exec.Command("ping", "-c", "1", ip).Output(); err == nil {
				// The connection was successful, log a message and set the flag to false
				log.Printf("[IPSEC Checker]["+hostname+"] Successful connection to %s", ip)

				failCount = 0
				continue
			} else {
				log.Printf("[IPSEC Checker][ERROR] Failed connection to %s", ip)
				failCount++

				// The connection failed, send a POST request with the IP address
				sendToHorn(fmt.Sprintf("*IPSEC Checker "+strconv.Itoa(failCount)+"/"+strconv.Itoa(probe)+"* ```"+hostname+"``` Connection to %s failed ‼️", ip))

				if failCount >= probe {
					log.Printf("[IPSEC Checker] Restarting IPSec!")
					err := exec.Command("ipsec", "restart").Run()
					if err != nil {
						log.Printf("[IPSEC Checker][ERROR] Restarting IPSec!")
						sendToHorn(fmt.Sprintf("*IPSEC Checker* ```" + hostname + "``` Restart IPSec Failed! Try manually! ```sudo ipsec restart``` 💀"))
						log.Fatal(err)
					}
					sendToHorn(fmt.Sprintf("*IPSEC Checker* ```" + hostname + "``` Restarting IPSec! ⚠️"))
					time.Sleep(30 * time.Second)

					// Check the status of the ipsec service
					output, _ := exec.Command("ipsec", "status").Output()
					fmt.Println(string(output))
					failCount = 0
					break
				}
			}
		}
	}
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		// An error occurred, handle it
		fmt.Println(err)
		return
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	environmentPath := filepath.Join(dir, ".env")
	linuxEnvironmentPath := filepath.Join("/usr/local/etc/ipsec_checker", ".env")
	err = godotenv.Load(environmentPath)
	errLinuxConfigLoading := godotenv.Load(linuxEnvironmentPath)
	if err != nil && errLinuxConfigLoading != nil {
		log.Fatal("Error loading .env file \n Check .env in current directory or in /usr/local/etc/ipsec_checker")
	}

	appEnv := os.Getenv("APP_ENV")

	if appEnv == "production" {
		err := raven.SetDSN(os.Getenv("SENTRY_DSN"))
		if err != nil {
			log.Println(err)
		}
	}

	// Set the list of IP addresses to check
	// read ip from env file comma separated
	ips := strings.Split(os.Getenv("IP_LIST"), ",")

	if len(ips) == 0 {
		panic("IP_LIST is empty")
	}

	// convert probe to int
	probe, err := strconv.Atoi(os.Getenv("PROBE"))
	if err != nil {
		panic("[IPSEC Checker] ENV PROBE is not valid! Only integer allowed")
	}

	// convert probe to int
	sleep, err := strconv.Atoi(os.Getenv("SLEEP"))
	if err != nil {
		panic("[IPSEC Checker] ENV SLEEP is not valid! Only integer allowed")
	}

	checkConnections(ips, hostname, probe, sleep)
}
