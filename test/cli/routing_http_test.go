package cli

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/ipfs/kubo/config"
	"github.com/ipfs/kubo/test/cli/harness"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/stretchr/testify/assert"
)

func TestRoutingV1(t *testing.T) {
	t.Parallel()
	nodes := harness.NewT(t).NewNodes(5).Init()
	nodes.ForEachPar(func(node *harness.Node) {
		node.UpdateConfig(func(cfg *config.Config) {
			cfg.Gateway.ExposeRoutingAPI = config.True
			cfg.Routing.Type = config.NewOptionalString("dht")
		})
	})
	nodes.StartDaemons().Connect()

	type record struct {
		Protocol string
		Schema   string
		ID       peer.ID
		Addrs    []string
	}

	type providers struct {
		Providers []record
	}

	t.Run("Non-streaming response with Accept: application/json", func(t *testing.T) {
		t.Parallel()

		cid := nodes[2].IPFSAddStr("hello world")
		_ = nodes[3].IPFSAddStr("hello world")

		resp := nodes[1].GatewayClient().Get("/routing/v1/providers/"+cid, func(r *http.Request) {
			r.Header.Set("Accept", "application/json")
		})
		assert.Equal(t, resp.Headers.Get("Content-Type"), "application/json")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var providers *providers
		err := json.Unmarshal([]byte(resp.Body), &providers)
		assert.NoError(t, err)

		var peers []peer.ID
		for _, prov := range providers.Providers {
			peers = append(peers, prov.ID)
		}
		assert.Contains(t, peers, nodes[2].PeerID())
		assert.Contains(t, peers, nodes[3].PeerID())
	})

	t.Run("Streaming response with Accept: application/x-ndjson", func(t *testing.T) {
		t.Parallel()

		cid := nodes[1].IPFSAddStr("hello world")
		_ = nodes[4].IPFSAddStr("hello world")

		resp := nodes[0].GatewayClient().Get("/routing/v1/providers/"+cid, func(r *http.Request) {
			r.Header.Set("Accept", "application/x-ndjson")
		})
		assert.Equal(t, resp.Headers.Get("Content-Type"), "application/x-ndjson")
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var peers []peer.ID
		dec := json.NewDecoder(strings.NewReader(resp.Body))

		for {
			var record *record
			err := dec.Decode(&record)
			if errors.Is(err, io.EOF) {
				break
			}

			assert.NoError(t, err)
			peers = append(peers, record.ID)
		}

		assert.Contains(t, peers, nodes[1].PeerID())
		assert.Contains(t, peers, nodes[4].PeerID())
	})
}
