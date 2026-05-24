package interface_enumerator

type EnumeratorManifest struct {
	TapPrefix     string `json:"tap_prefix" yaml:"tap_prefix"`
	BridgePrefix  string `json:"bridge_prefix" yaml:"bridge_prefix"`
	TapStorage    []bool `json:"tap_counter" yaml:"tap_counter"`
	BridgeStorage []bool `json:"bridge_counter" yaml:"bridge_counter"`
}
