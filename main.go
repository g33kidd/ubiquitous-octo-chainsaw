package main

import (
	"./models"
	"./streaming"

	"io"
	"log"
	"net/http"

	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func init() {
	format.RegisterAll()
}

// TODO: Move these to its own package
type writeFlusher struct {
	httpflusher http.Flusher
	io.Writer
}

func (wf writeFlusher) Flush() error {
	wf.httpflusher.Flush()
	return nil
}

// Stream handles the PubSub queue for I assume what is transmitting packets or
// at least holding chunks?
// TODO: Read the docs on PubSub Queue or look at the code.
type Stream struct {
	que *pubsub.Queue
}

// Env is our application environment
type Env struct {
	db     *gorm.DB
	stream *streaming.Server
}

func main() {

	// Open the database connection.
	db, err := models.NewDB("sqlite3", "./database/development.db")
	if err != nil {
		log.Panic(err)
	}

	// Setup the streaming server
	server, err := streaming.NewStreamingServer()
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db, server}
	models.InitTables(env.db)

	env.stream.Start()

	// env.server.HandlePublish = streaming.HandlePublish
	// env.server.HandleConn = streaming.HandleConn
	// env.server.HandlePlay = streaming.HandlePlay

	// Just for testing... instead of migrations and all that stuff..

	// Creates the chat Hub
	// hub := NewHub()
	// go hub.run()

	// TODO: figure out where to move this stuff. Stays here or ??? ¯\_(ツ)_/¯
	// Setup some variables
	// server := &rtmp.Server{}
	// l := &sync.RWMutex{}
	// channels := map[string]*Stream{}
	//
	// // Handles publishing when something comes into the stream ¯\_(ツ)_/¯
	// // TODO: Read the docs on this more...
	// server.HandlePublish = func(conn *rtmp.Conn) {
	// 	streams, _ := conn.Streams()
	// 	streamKey := strings.Replace(conn.URL.Path, "/", "", -1)
	//
	// 	log.Println("handling publish for", streamKey)
	//
	// 	// Get the channel we're publishing to..
	// 	// If we can't find the channel, close the connection immediately.
	// 	var channel models.Channel
	// 	if db.Where("stream_key = ?", streamKey).First(&channel).RecordNotFound() {
	// 		log.Println("Stream key invalid for", streamKey, "closing connection.")
	// 		conn.Close()
	// 		return
	// 	}
	//
	// 	log.Println("Stream key is valid. Continuing.")
	//
	// 	l.Lock()
	// 	log.Println("locked and loaded. Creating the stream.")
	// 	stream := channels[channel.Username]
	// 	if stream == nil {
	// 		stream = &Stream{}
	// 		stream.que = pubsub.NewQueue()
	// 		stream.que.WriteHeader(streams)
	// 		channels[channel.Username] = stream
	// 		log.Println("Created stream", channel.Username)
	// 	} else {
	// 		stream = nil
	// 	}
	// 	l.Unlock()
	//
	// 	if stream == nil {
	// 		log.Println("Couldn't find stream... It is nil")
	// 		return
	// 	}
	//
	// 	avutil.CopyPackets(stream.que, conn)
	//
	// 	l.Lock()
	// 	delete(channels, channel.Username)
	// 	log.Println("Stopping stream", channel.Username)
	// 	l.Unlock()
	//
	// 	stream.que.Close()
	// }
	//
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	username := strings.Replace(r.URL.Path, "/", "", -1)
	// 	log.Println(r)
	// 	log.Println(username)
	//
	// 	var channel models.Channel
	// 	if db.Where("username = ?", username).First(&channel).RecordNotFound() {
	// 		http.NotFound(w, r)
	// 	}
	//
	// 	l.RLock()
	// 	stream := channels[channel.Username]
	// 	l.RUnlock()
	//
	// 	fmt.Printf("Handling http request\n")
	// 	fmt.Printf("%+v\n", stream)
	//
	// 	if stream != nil {
	// 		w.Header().Set("Content-Type", "video/x-flv")
	// 		w.Header().Set("Transfer-Encoding", "chunked")
	// 		w.Header().Set("Access-Control-Allow-Origin", "*")
	// 		w.WriteHeader(200)
	//
	// 		flusher := w.(http.Flusher)
	// 		flusher.Flush()
	//
	// 		muxer := flv.NewMuxerWriteFlusher(writeFlusher{httpflusher: flusher, Writer: w})
	// 		cursor := stream.que.Latest()
	//
	// 		avutil.CopyFile(muxer, cursor)
	// 	} else {
	// 		w.Header().Set("Access-Control-Allow-Origin", "*")
	// 		w.WriteHeader(200)
	// 	}
	// })
	//
	// go server.ListenAndServe()
	// http.ListenAndServe(":8089", nil)
}
