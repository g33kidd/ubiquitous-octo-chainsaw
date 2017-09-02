package streaming

import (
	"io"
	"log"
	"net/http"
	"strings"
	"sync"

	"../models"

	"github.com/jinzhu/gorm"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/flv"
	"github.com/nareix/joy4/format/rtmp"
)

// Stream has information about the stream
type Stream struct {
	que *pubsub.Queue
}

// Server handles RTMP requests.
// It also stores the current streams that are running.
type Server struct {
	rtmpServer *rtmp.Server
	l          *sync.RWMutex
	streams    map[string]*Stream
	db         *gorm.DB
}

type writeFlusher struct {
	httpflusher http.Flusher
	io.Writer
}

func (wf writeFlusher) Flush() error {
	wf.httpflusher.Flush()
	return nil
}

// NewStreamingServer creates a new streaming server
// TODO: Handle errors Errors
func NewStreamingServer(db *gorm.DB) (*Server, error) {
	server := &Server{}

	rtmpServer := &rtmp.Server{HandlePublish: server.HandlePublish}
	l := &sync.RWMutex{}
	streams := map[string]*Stream{}

	server.db = db
	server.rtmpServer = rtmpServer
	server.l = l
	server.streams = streams

	return server, nil
}

// Start : starts the RTMP server
func (server *Server) Start() {
	server.rtmpServer.ListenAndServe()
}

// HandlePublish handles incoming RTMP connections
func (server *Server) HandlePublish(conn *rtmp.Conn) {
	streamKey := ""
	params := strings.Split(conn.URL.Path, "/")
	streams, _ := conn.Streams()

	// Get the StreamKey from URL Path.
	// It should be the only path parameter.
	// Close the connection if we don't have a StreamKey
	if len(params) > 1 {
		streamKey = params[1]
	} else if len(params) == 0 {
		log.Println("No streamKey found")
		conn.Close()
		return
	}

	channel, err := models.FindChannelByStreamKey(server.db, streamKey)
	if err != nil {
		log.Println(err)
	}

	if channel == nil {
		conn.Close()
		return
	}

	// Initialize the stream. Send stream data over the Queue
	server.l.Lock()
	stream := server.streams[channel.Username]
	if stream == nil {
		stream = &Stream{}
		stream.que = pubsub.NewQueue()
		stream.que.WriteHeader(streams)
		server.streams[channel.Username] = stream
		log.Println("Created stream", channel.Username)
	} else {
		stream = nil
	}
	server.l.Unlock()

	if stream == nil {
		return
	}

	// Transmit the data...
	avutil.CopyPackets(stream.que, conn)

	// Stop the stream
	server.l.Lock()
	delete(server.streams, channel.Username)
	server.l.Unlock()

	// Close the PubSub Queue, we are done with it...
	log.Println("Stopping stream", channel.Username)
	stream.que.Close()
}

// HandleHTTP : Handles HTTP requests to a given stream.
// TODO: comment this stuff...
func (server *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {

	username := ""
	params := strings.Split(r.URL.Path, "/")

	// Get the StreamKey from URL Path.
	// It should be the only path parameter.
	// Close the connection if we don't have a StreamKey
	if len(params) > 1 {
		username = params[1]
	} else if len(params) == 0 {
		http.NotFound(w, r)
		return
	}

	// Sets the current stream from the list of streams.
	server.l.RLock()
	stream := server.streams[username]
	server.l.RUnlock()

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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
	}

}

// HandleConn handles incoming RTMP connections
func (server *Server) HandleConn(conn *rtmp.Conn) {

}
