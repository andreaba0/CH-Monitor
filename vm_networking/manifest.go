package vmnetworking

type Manifest struct {
	TapPrefix     string `json:"tap_prefix" yaml:"tap_prefix"`
	BridgePrefix  string `json:"bridge_prefix" yaml:"bridge_prefix"`
	TapCounter    uint32 `json:"tap_counter" yaml:"tap_counter"`
	BridgeCounter uint32 `json:"bridge_counter" yaml:"bridge_counter"`
}
