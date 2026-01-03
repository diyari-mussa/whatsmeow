package whatsmeow

import (
	"context"
	"fmt"

	waBinary "go.mau.fi/whatsmeow/binary"
)

// SendNode sends a raw XMPP node to the server.
func (cli *Client) SendNode(ctx context.Context, node waBinary.Node) error {
	return cli.sendNode(ctx, node)
}

// SendIQ sends an IQ node and waits for the response.
func (cli *Client) SendIQ(ctx context.Context, node waBinary.Node) (*waBinary.Node, error) {
	if node.Tag != "iq" {
		return nil, fmt.Errorf("node tag must be iq")
	}
	id, ok := node.Attrs["id"].(string)
	if !ok {
		id = string(cli.GenerateMessageID())
		if node.Attrs == nil {
			node.Attrs = make(waBinary.Attrs)
		}
		node.Attrs["id"] = id
	}

	waiter := cli.waitResponse(id)

	err := cli.sendNode(ctx, node)
	if err != nil {
		cli.cancelResponse(id, waiter)
		return nil, err
	}

	select {
	case resp := <-waiter:
		return resp, nil
	case <-ctx.Done():
		cli.cancelResponse(id, waiter)
		return nil, ctx.Err()
	}
}
