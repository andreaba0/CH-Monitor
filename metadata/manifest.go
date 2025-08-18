package metadata

type Manifest struct {
	TapCounter    int    `json:"tap_counter" yaml:"tap_counter"`
	TapPrefix     string `json:"tap_prefix" yaml:"tap_prefix"`
	BridgeCounter int    `json:"bridge_counter" yaml:"bridge_counter"`
	BridgePrefix  string `json:"bridge_prefix" yaml:"bridge_prefix"`
}
