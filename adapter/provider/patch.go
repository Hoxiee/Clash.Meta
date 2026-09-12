package provider

import "sync"

func (pp *proxySetProvider) GetSubscriptionInfo() *SubscriptionInfo {
	return pp.subscriptionInfo
}

var (
	screenMu     sync.Mutex
	screenIsOff  bool
	screenWakeCh = make(chan struct{})
)

// Periodic probes are the tunnel's main idle-time radio wake-up, and a group
// nothing routed through while the screen was off has no reason to be probed.
// The tick is deferred rather than dropped: it runs when the screen comes back,
// so the delays a user sees are never the ones measured in their pocket.
func SetScreenOff(off bool) {
	screenMu.Lock()
	defer screenMu.Unlock()
	if screenIsOff == off {
		return
	}
	screenIsOff = off
	if off {
		return
	}
	close(screenWakeCh)
	screenWakeCh = make(chan struct{})
}

func screenOff() bool {
	screenMu.Lock()
	defer screenMu.Unlock()
	return screenIsOff
}

func screenWake() <-chan struct{} {
	screenMu.Lock()
	defer screenMu.Unlock()
	return screenWakeCh
}
