package main

import (
    "sync"
    "fmt"
	"io"
	"net/http"
	"github.com/nareix/joy4/format"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/rtmp"
	"github.com/nareix/joy4/format/flv"
)

// The channel type for multiple channel usage ¯\_(ツ)_/¯
type Channel struct {
    queue *pubsub.Queue
}

func main() {

    // Setup some variables
    server      := &rtmp.Server{}
    rwmutex     := &sync.RWMutex{}
    channels    := map[string] *Channel{}

    // Handles publishing when something comes into the stream ¯\_(ツ)_/¯
    // TODO: Read the docs on this more...
    server.HandlePublish = func(conn *rtmp.Conn) {
        streams, _ := conn.Streams()

        // Lock it to prevent further writes...
        rwmutex.Lock()

        // The current channel for this stream
        channel := channels[conn.URL.Path]
        if channel == nil {
            channel = &Channel{}
            channel.queue = pubsub.NewQueue()
            channel.queue.WriteHeader(streams)
            channels[conn.URL.Path] = channel
        } else {
            channel = nil
        }

        // Unlock the mutex
        rwmutex.Unlock()

        if channel == nil {
            return
        }

        // Copies the packets that are currently being published.
        avutil.CopyPackets(channel.queue, conn)
    }

}
