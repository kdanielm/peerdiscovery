package peerdiscovery

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiscovery(t *testing.T) {
	for _, version := range []IPVersion{IPv4, IPv6} {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

		// should not be able to "discover" itself
		discoveries, err := Discover(ctx, Settings{
			Delay: 500 * time.Millisecond,
		})
		assert.Nil(t, err)
		assert.Zero(t, len(discoveries))
		cancel()

		ctx, cancel = context.WithTimeout(context.TODO(), 1*time.Second)
		// should be able to "discover" itself
		discoveries, err = Discover(ctx, Settings{
			Limit:     -1,
			AllowSelf: true,
			Payload:   []byte("payload"),
			Delay:     500 * time.Millisecond,
			IPVersion: version,
		})
		fmt.Println(discoveries)
		assert.Nil(t, err)
		assert.NotZero(t, len(discoveries))
		cancel()
	}
}

func TestDiscoverySelf(t *testing.T) {
	for _, version := range []IPVersion{IPv4, IPv6} {
		// broadcast self to self
		go func() {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			_, err := Discover(
				ctx,
				Settings{
					Limit:     -1,
					Payload:   []byte("payload"),
					Delay:     10 * time.Millisecond,
					IPVersion: version,
				},
			)
			assert.Nil(t, err)
			cancel()
		}()

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		discoveries, err := Discover(
			ctx,
			Settings{
				Limit:            1,
				Payload:          []byte("payload"),
				Delay:            500 * time.Millisecond,
				DisableBroadcast: true,
				AllowSelf:        true,
			},
		)
		assert.Nil(t, err)
		assert.NotZero(t, len(discoveries))
		cancel()
	}
}

func TestTimeout(t *testing.T) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

	discoveries, err := Discover(
		ctx,
		Settings{
			Limit:            -1,
			Payload:          []byte("payload"),
			Delay:            500 * time.Millisecond,
			DisableBroadcast: true,
		},
	)

	assert.Nil(t, err)
	assert.Zero(t, len(discoveries))
	assert.Greater(t, time.Since(start).Seconds(), 5.0)

	cancel()
}

func TestCancel(t *testing.T) {
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)

	go func() {
		time.AfterFunc(2*time.Second, cancel)
	}()

	discoveries, err := Discover(
		ctx,
		Settings{
			Limit:            -1,
			Payload:          []byte("payload"),
			Delay:            500 * time.Millisecond,
			DisableBroadcast: true,
		},
	)

	assert.Nil(t, err)
	assert.Zero(t, len(discoveries))
	assert.Greater(t, time.Since(start).Seconds(), 2.0)
	assert.Less(t, time.Since(start).Seconds(), 2.5)
}
