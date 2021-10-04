// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package plugin

import (
	"context"

	appproto "github.com/ava-labs/avalanchego/app/plugin/proto"
	"github.com/ava-labs/avalanchego/app/process"
)

// Server wraps a node so it can be served with the hashicorp plugin harness
type Server struct {
	appproto.UnimplementedNodeServer
	app *process.App
}

func NewServer(app *process.App) *Server {
	return &Server{
		app: app,
	}
}

// Blocks until the node returns
func (ns *Server) Start(_ context.Context, req *appproto.StartRequest) (*appproto.StartResponse, error) {
	exitCode := ns.app.Start()
	return &appproto.StartResponse{ExitCode: int32(exitCode)}, nil
}

// Blocks until the node is done shutting down
func (ns *Server) Stop(_ context.Context, req *appproto.StopRequest) (*appproto.StopResponse, error) {
	ns.app.Stop()
	return &appproto.StopResponse{}, nil
}
