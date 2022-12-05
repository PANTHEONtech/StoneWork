// Code generated by GoVPP's binapi-generator. DO NOT EDIT.

package isisx

import (
	"context"
	"fmt"
	"io"

	api "go.fd.io/govpp/api"
	memclnt "go.pantheon.tech/stonework/plugins/binapi/vpp2210/memclnt"
)

// RPCService defines RPC service isisx.
type RPCService interface {
	IsisxConnectionAddDel(ctx context.Context, in *IsisxConnectionAddDel) (*IsisxConnectionAddDelReply, error)
	IsisxConnectionDump(ctx context.Context, in *IsisxConnectionDump) (RPCService_IsisxConnectionDumpClient, error)
	IsisxPluginGetVersion(ctx context.Context, in *IsisxPluginGetVersion) (*IsisxPluginGetVersionReply, error)
}

type serviceClient struct {
	conn api.Connection
}

func NewServiceClient(conn api.Connection) RPCService {
	return &serviceClient{conn}
}

func (c *serviceClient) IsisxConnectionAddDel(ctx context.Context, in *IsisxConnectionAddDel) (*IsisxConnectionAddDelReply, error) {
	out := new(IsisxConnectionAddDelReply)
	err := c.conn.Invoke(ctx, in, out)
	if err != nil {
		return nil, err
	}
	return out, api.RetvalToVPPApiError(out.Retval)
}

func (c *serviceClient) IsisxConnectionDump(ctx context.Context, in *IsisxConnectionDump) (RPCService_IsisxConnectionDumpClient, error) {
	stream, err := c.conn.NewStream(ctx)
	if err != nil {
		return nil, err
	}
	x := &serviceClient_IsisxConnectionDumpClient{stream}
	if err := x.Stream.SendMsg(in); err != nil {
		return nil, err
	}
	if err = x.Stream.SendMsg(&memclnt.ControlPing{}); err != nil {
		return nil, err
	}
	return x, nil
}

type RPCService_IsisxConnectionDumpClient interface {
	Recv() (*IsisxConnectionDetails, error)
	api.Stream
}

type serviceClient_IsisxConnectionDumpClient struct {
	api.Stream
}

func (c *serviceClient_IsisxConnectionDumpClient) Recv() (*IsisxConnectionDetails, error) {
	msg, err := c.Stream.RecvMsg()
	if err != nil {
		return nil, err
	}
	switch m := msg.(type) {
	case *IsisxConnectionDetails:
		return m, nil
	case *memclnt.ControlPingReply:
		err = c.Stream.Close()
		if err != nil {
			return nil, err
		}
		return nil, io.EOF
	default:
		return nil, fmt.Errorf("unexpected message: %T %v", m, m)
	}
}

func (c *serviceClient) IsisxPluginGetVersion(ctx context.Context, in *IsisxPluginGetVersion) (*IsisxPluginGetVersionReply, error) {
	out := new(IsisxPluginGetVersionReply)
	err := c.conn.Invoke(ctx, in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
