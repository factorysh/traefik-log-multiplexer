package output

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFluent(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	assert.NoError(t, err)
	var msg []byte
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		buff := make([]byte, 1024*1024)
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("accept", err)
				continue
			}
			n, err := conn.Read(buff)
			if err != nil {
				panic(err)
			}
			if n > 0 {
				msg = buff[:n]
				wg.Done()
			}
		}
	}()
	addr := strings.Split(l.Addr().String(), ":")
	fmt.Println(addr)
	port, err := strconv.ParseInt(addr[1], 10, 32)
	assert.NoError(t, err)
	f, err := FluentOutputFactory(map[string]interface{}{
		"labels": []string{"pim"},
		"config": map[string]interface{}{
			"fluentHost": addr[0],
			"fluentPort": int(port),
			"async":      false,
			"requestAck": false,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, f)
	defer f.Close()
	ff, ok := f.(*FluentOutput)
	assert.True(t, ok)
	assert.Equal(t, int(port), ff.config.Config.FluentPort)
	err = f.Write(context.TODO(), time.Now(), `{"beuha": "aussi", "time":"2021-05-04T16:54:20Z"}`, map[string]interface{}{"age": 42})
	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)
	wg.Wait()
	fmt.Println(msg)
	l.Close()
}
