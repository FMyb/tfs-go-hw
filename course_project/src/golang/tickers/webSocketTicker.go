package tickers

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/FMyb/tfs-go-hw/course_project/course_project/src/golang/domain"
	"github.com/gorilla/websocket"
)

type WebSocketTicker struct {
	interrupt chan interface{}
	conn      *websocket.Conn
	done      chan interface{}
	stoped    chan interface{}
}

func NewWebSocketTicker(interrupt chan interface{}) *WebSocketTicker {
	return &WebSocketTicker{interrupt: interrupt}
}

func (ws *WebSocketTicker) NextTicker() domain.ResponseStatus {
stopFor:
	for {
		select {
		case <-ws.done:
			break stopFor
		default:
			_, msg, err := ws.conn.ReadMessage()
			if err != nil {
				log.Errorf("Error in receive: %s", err)
				return nil
			}
			msgJson := domain.ResponseTicker{}
			err = json.Unmarshal(msg, &msgJson)
			if msgJson.Event != "" && msgJson.Message != "" {
				log.Errorf("Error in read ticker: %s", msgJson.Message)
				return nil
			}
			if msgJson.Time == 0 {
				continue
			}
			return domain.ToSuccessResponseTicker(&msgJson)
		}
	}
	return nil
}

func (ws *WebSocketTicker) Start(productIDs []string) {
	//ws.done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	//ws.interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	//signal.Notify(ws.interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	socketUrl := "wss://demo-futures.kraken.com/ws/v1"
	var err error
	ws.conn, _, err = websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	_, _, err = ws.conn.ReadMessage()
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	// Our main loop for the client
	// We send our relevant packets here
	go func() {
		defer ws.conn.Close()
		for {
			select {
			case <-time.After(time.Duration(1) * time.Millisecond * 100):
				// Send an echo packet every second
				err = ws.conn.WriteMessage(websocket.PingMessage, []byte(""))
				if err != nil {
					log.Errorf("Error during writing to websocket: %s", err)
					return
				}
			case <-ws.interrupt:
				log.Info("Stop web socket connection!")
				request := domain.UnsubscribeTickerRequestWS(productIDs)
				requestJson, err := json.Marshal(&request)
				if err != nil {
					log.Fatal("Error in marshal request:", err)
					return
				}
				err = ws.conn.WriteMessage(websocket.TextMessage, requestJson)
				if err != nil {
					log.Errorf("Error in unsubscribe tickers")
				}
				// Close our websocket connection
				err = ws.conn.WriteMessage(
					websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
				)
				if err != nil {
					log.Errorf("Error during closing websocket: %s", err)
					return
				}
				ws.stoped <- nil
				return
			}
		}
	}()
	request := domain.SubscribeTickerRequestWS(productIDs)
	requestJson, err := json.Marshal(&request)
	if err != nil {
		log.Fatal("Error in marshal request:", err)
		return
	}
	connected := false
	for i := 0; i < 10; i++ {
		err = ws.conn.WriteMessage(websocket.TextMessage, requestJson)
		if err != nil {
			log.Errorf("Can't subscribe to ticker: %s", requestJson)
			log.Errorf("Try again...")
			<-time.After(time.Duration(1) * time.Second)
			continue
		}
		connected = true
		break
	}
	if !connected {
		log.Fatalf("Can't subscribe to ticker: %s", requestJson)
	}

}

func (ws *WebSocketTicker) Stop() chan interface{} {
	ws.interrupt <- 0
	ws.done <- 0
	return ws.stoped
}
