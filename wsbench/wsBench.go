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
	var compareError uint64 = 0
	var readError uint64 = 0
	var writeError uint64 = 0
	var connectionError uint64 = 0
	var writeBytes uint64 = 0
	var readBytes uint64 = 0

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

			if err != nil {
				log.Println("dial:", err)
				connectionError++
			} else {
				defer co.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				for {
					if time.Since(start) > (time.Duration(duration) * time.Second) {
						break
					}
					writeTime := time.Now()
					err = co.WriteMessage(websocket.TextMessage, []byte(message))
					if err != nil {
						log.Println("write:", err)
						//writeError++
						atomic.AddUint64(&writeError, 1)
					} else {
						//writeBytes += len([]byte(message))
						atomic.AddUint64(&writeBytes, uint64(len([]byte(message))))
						//writeCounter++
						atomic.AddUint64(&writeCounter, 1)
						//writeChan <- 1

						_, readMessage, err := co.ReadMessage()
						if err != nil {
							log.Println("read:", err)
							readError++
						} else if string(readMessage) != message {
							log.Printf("received message is not the same! recv: %s", readMessage)
							//compareError++
							atomic.AddUint64(&compareError, 1)
							//readCounter++
							atomic.AddUint64(&readCounter, 1)
							//	readChan <- 1
						} else {
							//readBytes += len(readMessage)
							atomic.AddUint64(&readBytes, uint64(len([]byte(readMessage))))
							//readCounter++
							//	readChan <- 1
							atomic.AddUint64(&readCounter, 1)
						}
					}
					dur := time.Since(writeTime)
					log.Println(dur)
					durr += dur
					//readCounter++
					//ch <- counter
					time.Sleep(time.Duration(interval) * time.Second)
				}
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

	//log.Println("here")
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
	writeBytesF := atomic.LoadUint64(&writeBytes)
	readBytesF := atomic.LoadUint64(&readBytes)
	connectionErrorF := atomic.LoadUint64(&connectionError)
	writeErrorF := atomic.LoadUint64(&writeError)
	readErrorF := atomic.LoadUint64(&readError)
	compareErrorF := atomic.LoadUint64(&compareError)

	//	for {
	//		<-ch
	//	}
	//log.Println("wtf:", writeC, readC)
	var averageRtt time.Duration
	if readCounterF == 0 {
		averageRtt = time.Duration(0)
	} else {
		averageRtt = durr / time.Duration(readCounterF)
	}
	log.Println("Total Sent:", writeCounterF, ", Total Received:", readCounterF, ", Bytes Sent", writeBytesF, ", Bytes Received:", readBytesF, ", Average RTT:", averageRtt, ", Connection Error:", connectionErrorF, ", Write Error:", writeErrorF, ", Read Error:", readErrorF, ", Message Mismatch:", compareErrorF)
}
