package streaming

import (
	"sync"

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
		HandleConn:    server.HandleConn,
	}
	l := &sync.RWMutex{}
	streams := map[string]*Stream{}

	server.rtmpServer = rtmpServer
	server.l = l
	server.streams = streams

	return server, nil
}

func (server *Server) start() {
	go server.rtmpServer.ListenAndServe()
}

// HandlePublish handles incoming RTMP connections
func (server *Server) HandlePublish(conn *rtmp.Conn) {

}

// HandleConn handles incoming RTMP connections
func (server *Server) HandleConn(conn *rtmp.Conn) {

}
