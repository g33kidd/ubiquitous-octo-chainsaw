package streaming

import (
	"log"
	"strings"
	"sync"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
)

// Stream has information about the stream
type Stream struct {
	que *pubsub.Queue
}

// Server handles RTMP requests
type Server struct {
	rtmpServer *rtmp.Server
	l          *sync.RWMutex
	streams    map[string]*Stream
}

// NewStreamingServer creates a new streaming server
// TODO: Handle errors Errors
func NewStreamingServer() (*Server, error) {
	server := &Server{}

	rtmpServer := &rtmp.Server{
		HandlePublish: server.HandlePublish,
		// HandleConn:    server.HandleConn,
	}
	l := &sync.RWMutex{}
	streams := map[string]*Stream{}

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
	var streamKey string

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

	// Initialize the stream. Send stream data over the Queue
	server.l.Lock()
	stream := server.streams[streamKey]
	if stream == nil {
		stream = &Stream{}
		stream.que = pubsub.NewQueue()
		stream.que.WriteHeader(streams)
		server.streams[streamKey] = stream
		log.Println("Created stream", stream)
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
	delete(server.streams, streamKey)
	server.l.Unlock()

	// Close the PubSub Queue, we are done with it...
	stream.que.Close()
}

// HandleConn handles incoming RTMP connections
func (server *Server) HandleConn(conn *rtmp.Conn) {

}
