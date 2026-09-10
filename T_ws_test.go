package EasyOnebot

/*
func TestXNetWebsocket(t *testing.T) {
	// c1, err1 := websocket.Dial("ws://127.0.0.1:8081/event", "", "http://127.0.0.1")
	c1, err1 := websocket.Dial("ws://192.168.1.103:8081/event", "", "http://127.0.0.1")
	// c2, err2 := websocket.Dial("ws://127.0.0.1:8081/api", "", "http://127.0.0.1")
	c2, err2 := websocket.Dial("ws://192.168.1.103:8081/api", "", "http://127.0.0.1")
	if err1 != nil {
		t.Fatal(err1)
	}
	if err2 != nil {
		t.Fatal(err2)
	}

	data := make([]byte, 0)
	err := websocket.Message.Receive(c1, &data)
	if err != nil {
		t.Logf("err: %v", err)
	} else {
		t.Logf("first data packet: %s", data)
	}

	go func() {
		for {
			data := make([]byte, 0)
			err := websocket.Message.Receive(c1, &data)
			if err != nil {
				t.Logf("err: %v", err)
				break
			}
			t.Logf("event data: %s", data)
		}
	}()
	go func() {
		for {
			data := make([]byte, 0)
			err := websocket.Message.Receive(c2, &data)
			if err != nil {
				t.Logf("err: %v", err)
				break
			}
			t.Logf("api data: %s", data)
		}
	}()

	<-time.After(10 * time.Second)
	c1.Close()
	c2.Close()
	t.Log("closed")
}
*/
