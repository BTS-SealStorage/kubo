package corehttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/ipfs/boxo/routing/http/server"
	"github.com/ipfs/boxo/routing/http/types"
	"github.com/ipfs/boxo/routing/http/types/iter"
	cid "github.com/ipfs/go-cid"
	core "github.com/ipfs/kubo/core"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/multiformats/go-multiaddr"
)

func RoutingOption() ServeOption {
	return func(n *core.IpfsNode, _ net.Listener, mux *http.ServeMux) (*http.ServeMux, error) {
		handler := server.Handler(&contentRouter{n})
		mux.Handle("/routing/v1/", handler)
		return mux, nil
	}
}

type contentRouter struct {
	n *core.IpfsNode
}

func (r *contentRouter) FindProviders(ctx context.Context, key cid.Cid, limit int) (iter.ResultIter[types.ProviderResponse], error) {
	ctx, cancel := context.WithCancel(ctx)
	ch := r.n.Routing.FindProvidersAsync(ctx, key, limit)
	return iter.ToResultIter[types.ProviderResponse](&peerChanIter{
		ch:     ch,
		cancel: cancel,
	}), nil
}

func (r *contentRouter) Provide(ctx context.Context, req *server.WriteProvideRequest) (types.ProviderResponse, error) {
	// Kubo /routing/v1 endpoint does not support write operations.
	return nil, routing.ErrNotSupported
}

func (r *contentRouter) ProvideBitswap(ctx context.Context, req *server.BitswapWriteProvideRequest) (time.Duration, error) {
	// Kubo /routing/v1 endpoint does not support write operations.
	return 0, routing.ErrNotSupported
}

type peerChanIter struct {
	ch     <-chan peer.AddrInfo
	cancel context.CancelFunc
	next   *peer.AddrInfo
}

func (it *peerChanIter) Next() bool {
	addr, ok := <-it.ch
	if ok {
		it.next = &addr
		return true
	} else {
		it.next = nil
		return false
	}
}

func (it *peerChanIter) Val() types.ProviderResponse {
	if it.next == nil {
		return nil
	}

	// We don't know what type of protocol this peer provides. It is likely Bitswap
	// but it might not be. Therefore, we set an unknown protocol with an unknown schema.
	rec := &providerRecord{
		Protocol: "transport-unknown",
		Schema:   "unknown",
		ID:       it.next.ID,
		Addrs:    it.next.Addrs,
	}

	return rec
}

func (it *peerChanIter) Close() error {
	it.cancel()
	return nil
}

type providerRecord struct {
	Protocol string
	Schema   string
	ID       peer.ID
	Addrs    []multiaddr.Multiaddr
}

func (pr *providerRecord) GetProtocol() string {
	return pr.Protocol
}

func (pr *providerRecord) GetSchema() string {
	return pr.Schema
}
