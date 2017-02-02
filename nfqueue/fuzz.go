// +build gofuzz

package nfqueue

import (
	"net"
	"sync"
)

// Fuzz uses go-fuzz https://github.com/dvyukov/go-fuzz
// we need to setup a NFQ at Queue ID 2 in the iptables rules
// as well as setup a listener on the assigned port
// socat -v tcp-l:6612,fork exec:'cat'
func Fuzz(data []byte) int {
	var err error
	nfq := NewNFQueue(2)
	receiveChan := make(<-chan *NFQPacket, 0)

	receiveChan, err = nfq.Open()
	if err != nil {
		panic(err)
	}
	defer nfq.Close()

	wg := sync.WaitGroup{}

	go func() {
		// receive packet
		packet := <-receiveChan
		packet.Accept()
		packet = <-receiveChan
		packet.Accept()
		packet = <-receiveChan
		packet.Accept()
		wg.Done()
	}()
	wg.Add(1)

	// send packet
	conn, err := net.Dial("tcp", "127.0.0.1:6612")
	if err != nil {
		// handle error
		panic(err)
	}
	_, err = conn.Write(data)
	if err != nil {
		panic(err)
	}

	wg.Wait()
	conn.Close()

	return 0
}
