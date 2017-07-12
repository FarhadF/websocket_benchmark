package wsbench

import (
	_ "fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	_ "os"
	_ "os/signal"
	"sync"
	"sync/atomic"
	"time"
)

func WsBench(address string, path string, sockets int, interval int, message string, duration int) {
	u := url.URL{Scheme: "ws", Host: address, Path: path}
	log.Printf("connecting to %s", u.String())
	start := time.Now()
	counter := 0
	var readCounter uint64 = 0
	var writeCounter uint64 = 0
	compareError := 0
	readError := 0
	writeError := 0
	connectionError := 0
	writeBytes := 0
	readBytes := 0

	//readChan := make(chan int)
	//writeChan := make(chan int)
	//controlChan := make(chan int)
	var durr time.Duration
	var wg sync.WaitGroup
	for {
		counter++
		wg.Add(1)
		go func() {
			co, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			defer co.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
					//writeCounter++
					atomic.AddUint64(&writeCounter, 1)
					//writeChan <- 1
				}
				_, readMessage, err := co.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					readError++
				} else if string(readMessage) != message {
					log.Printf("received message is not the same! recv: %s", readMessage)
					compareError++
					//readCounter++
					atomic.AddUint64(&readCounter, 1)
					//	readChan <- 1
				} else {
					readBytes += len(readMessage)
					//readCounter++
					//	readChan <- 1
					atomic.AddUint64(&readCounter, 1)
				}
				dur := time.Since(writeTime)
				log.Println(dur)
				durr += dur
				//readCounter++
				//ch <- counter
				time.Sleep(time.Duration(interval) * time.Second)
			}
			defer wg.Done()
			//controlChan <- 1
		}()

		//time.Sleep(1 * time.Millisecond)
		if counter >= sockets {
			break
		}
	}
	/*	for i := 0; i < counter; i++ {
		<-controlChan
	}*/

	log.Println("here")
	/*	var readC int
		var writeC int
		for i := 0; i < len(readChan); i++ {
			temp := <-readChan
			readC += temp
			log.Println("wtf")

		}
		for i := 0; i < len(writeChan); i++ {
			temp := <-writeChan
			writeC += temp
		}
	*/
	wg.Wait()
	readCounterF := atomic.LoadUint64(&readCounter)
	writeCounterF := atomic.LoadUint64(&writeCounter)
	//	for {
	//		<-ch
	//	}
	//log.Println("wtf:", writeC, readC)
	log.Println("Total Sent:", writeCounterF, "Total Received:", readCounterF, "Bytes Sent", writeBytes, "Bytes Received:", readBytes, "Average RTT:", (durr / time.Duration(readCounter)), "Connection Error:", connectionError, "Write Error:", writeError, "Read Error:", readError, "Message Mismatch:", compareError)
}
