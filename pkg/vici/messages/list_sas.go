package messages

type ListSAS struct {
	ReKeyTime int64  `vici:"rekey_time"`
	State     string `vici:"state"`
	UniqueID  int64  `vici:"uniqueid"`
	Version   string `vici:"version"`

	Children map[string]ListSASChildSA `vici:"child-sas"`
}

type ListSASChildSA struct {
	Name        string `vici:"name"`
	State       string `vici:"state"`
	BytesIn     int64  `vici:"bytes-in"`
	BytesOut    int64  `vici:"bytes-out"`
	InstallTime int64  `vici:"install-time"`
	LifeTime    int64  `vici:"life-time"`
	ReKeyTime   int64  `vici:"rekey-time"`

	LocalTS  []string `vici:"local-ts"`
	RemoteTS []string `vici:"remote-ts"`
}
