package messages

type ListConn struct {
	// local_addrs = [
	//   <list of valid local IKE endpoint addresses>
	// ]
	LocalAddrs []string `vici:"local_addrs"`
	// remote_addrs = [
	//   <list of valid remote IKE endpoint addresses>
	// ]
	RemoteAddrs []string `vici:"remote_addrs"`
	// version = <IKE version as string, IKEv1|IKEv2 or 0 for any>
	Version string `vici:"version"`
	// reauth_time = <IKE_SA reauthentication interval in seconds>
	ReAuthTime int64 `vici:"reauth_time"`
	// rekey_time = <IKE_SA rekeying interval in seconds>
	ReKeyTime int64 `vici:"rekey_time"`
	// CHILD_SA config name>*
	Children map[string]ListConnChildSA `vici:"children"`

	// multiple local and remote auth sections
	// the keys are dynamic and will be local-* or remote-*
	LocalAuth  map[string]ListConnAuthSection
	RemoteAuth map[string]ListConnAuthSection
}

type ListConnAuthSection struct {
	// class = <authentication type>
	Class string `vici:"class"`
	// eap-type = <EAP type to authenticate if when using EAP>
	EAPEapType string `vici:"eap-type"`
	// eap-vendor = <EAP vendor for type, if any>
	EAPEapVendor string `vici:"eap-vendor"`
	// xauth = <xauth backend name>
	XAuth string `vici:"xauth"`
	// revocation = <revocation policy>
	Revocation string `vici:"revocation"`
	// id = <IKE identity>
	ID string `vici:"id"`
	// aaa_id = <AAA authentication backend identity>
	AAAID string `vici:"aaa_id"`
	// eap_id = <EAP identity for authentication>
	EAPID string `vici:"eap_id"`
	// xauth_id = <XAuth username for authentication>
	XAuthID string `vici:"xauth_id"`
	// groups = [
	//   <group membership required to use connection>
	// ]
	Groups []string `vici:"groups"`
	// certs = [
	//   <certificates allowed for authentication>
	// ]
	Certs []string `vici:"certs"`
	// cacerts = [
	//   <CA certificates allowed for authentication>
	// ]
	CACerts []string `vici:"cacerts"`
}

type ListConnChildSA struct {
	// mode = <IPsec mode>
	Mode string `vici:"mode"`
	// label = <hex encoded security label>
	Label string `vici:"label"`
	// rekey_time = <CHILD_SA rekeying interval in seconds>
	ReKeyTime int64 `vici:"rekey_time"`
	// rekey_bytes = <CHILD_SA rekeying interval in bytes>
	ReKeyBytes int64 `vici:"rekey_bytes"`
	// rekey_packets = <CHILD_SA rekeying interval in packets>
	ReKeyPackets int64 `vici:"rekey_packets"`
	// local-ts = [
	//   <list of local traffic selectors>
	// ]
	LocalTS []string `vici:"local-ts"`
	// remote-ts = [
	//   <list of remote traffic selectors>
	// ]
	RemoteTS []string `vici:"remote-ts"`
}
