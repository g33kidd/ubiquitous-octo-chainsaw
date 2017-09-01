package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/format/flv"
	"github.com/nareix/joy4/format/rtmp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	format.RegisterAll()
}

type writeFlusher struct {
	httpflusher http.Flusher
	io.Writer
}

func (wf writeFlusher) Flush() error {
	wf.httpflusher.Flush()
	return nil
}

// Stream something something
type Stream struct {
	que *pubsub.Queue
}

// Channel for the user account and for the stream to be hosted on
type Channel struct {
	gorm.Model

	Username  string
	StreamKey string
}

func main() {

	db, err := gorm.Open("sqlite3", "./database/development.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	db.DropTable(&Channel{})
	db.CreateTable(&Channel{})
	db.Create(&Channel{Username: "g33kidd", StreamKey: "stream_key_12345"})

	// Creates the chat Hub
	hub := NewHub()
	go hub.run()

	// Setup some variables
	server := &rtmp.Server{}
	l := &sync.RWMutex{}
	channels := map[string]*Stream{}

	// Handles publishing when something comes into the stream ¯\_(ツ)_/¯
	// TODO: Read the docs on this more...
	server.HandlePublish = func(conn *rtmp.Conn) {
		streams, _ := conn.Streams()
		streamKey := strings.Replace(conn.URL.Path, "/", "", -1)

		log.Println("handling publish for", streamKey)

		// Get the channel we're publishing to..
		// If we can't find the channel, close the connection immediately.
		var channel Channel
		if db.Where("stream_key = ?", streamKey).First(&channel).RecordNotFound() {
			log.Println("Stream key invalid for", streamKey, "closing connection.")
			conn.Close()
			return
		}

		log.Println("Stream key is valid. Continuing.")

		l.Lock()
		log.Println("locked and loaded. Creating the stream.")
		stream := channels[channel.Username]
		if stream == nil {
			stream = &Stream{}
			stream.que = pubsub.NewQueue()
			stream.que.WriteHeader(streams)
			channels[channel.Username] = stream
			log.Println("Created stream", channel.Username)
		} else {
			stream = nil
		}
		l.Unlock()

		if stream == nil {
			log.Println("Couldn't find stream... It is nil")
			return
		}

		avutil.CopyPackets(stream.que, conn)

		l.Lock()
		delete(channels, channel.Username)
		log.Println("Stopping stream", channel.Username)
		l.Unlock()

		stream.que.Close()
	}

	// NOTE: We probably really don't need this.
	// TODO: Just read the stream as rtmp://
	// NOTE: Keeps getting invalid data ¯\_(ツ)_/¯
	// TODO: Look at updating the joy4 package if it needs updating
	// HTTP Handler for clients and plays for the server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		username := strings.Replace(r.URL.Path, "/", "", -1)
		log.Println(r)
		log.Println(username)

		var channel Channel
		if db.Where("username = ?", username).First(&channel).RecordNotFound() {
			http.NotFound(w, r)
		}

		l.RLock()
		stream := channels[channel.Username]
		l.RUnlock()

		fmt.Printf("Handling http request\n")
		fmt.Printf("%+v\n", stream)

		if stream != nil {
			w.Header().Set("Content-Type", "video/x-flv")
			w.Header().Set("Transfer-Encoding", "chunked")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.WriteHeader(200)

			flusher := w.(http.Flusher)
			flusher.Flush()

			muxer := flv.NewMuxerWriteFlusher(writeFlusher{httpflusher: flusher, Writer: w})
			cursor := stream.que.Latest()

			avutil.CopyFile(muxer, cursor)
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	go server.ListenAndServe()
	http.ListenAndServe(":8089", nil)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }
}
