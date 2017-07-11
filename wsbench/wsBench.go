package wsbench

import (
	_ "fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	_ "os"
	_ "os/signal"
	"sync"
	"time"
)

func WsBench(address string, path string, sockets int, interval int, message string, duration int) {
	u := url.URL{Scheme: "ws", Host: address, Path: path}
	log.Printf("connecting to %s", u.String())
	start := time.Now()
	counter := 0
	readCounter := 0
	writeCounter := 0
	compareError := 0
	readError := 0
	writeError := 0
	connectionError := 0
	writeBytes := 0
	readBytes := 0
	var durr time.Duration
	var wg sync.WaitGroup
	for {
		counter++
		wg.Add(1)
		go func() {
			co, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
				connectionError++
			}
			for {
				if time.Since(start) > (time.Duration(duration) * time.Second) {
					break
				}
				writeTime := time.Now()
				err = co.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					log.Println("write:", err)
					writeError++
				} else {
					writeBytes += len([]byte(message))
					writeCounter++
				}
				_, readMessage, err := co.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					readError++
				} else if string(readMessage) != message {
					log.Printf("received message is not the same! recv: %s", readMessage)
					compareError++
				} else {
					readBytes += len(readMessage)
				}
				dur := time.Since(writeTime)
				log.Println(dur)
				durr += dur
				readCounter++
				//ch <- counter
				time.Sleep(time.Duration(interval) * time.Second)
			}
			defer wg.Done()

		}()

		//time.Sleep(1 * time.Millisecond)
		if counter >= sockets {
			break
		}
	}
	wg.Wait()
	//	for {
	//		<-ch
	//	}
	log.Println("Total Sent:", writeCounter, "Total Received:", readCounter, "Bytes Sent", writeBytes, "Bytes Received:", readBytes, "Average RTT:", (durr / time.Duration(readCounter)), "Connection Error:", connectionError, "Write Error:", writeError, "Read Error:", readError, "Message Mismatch:", compareError)
}
