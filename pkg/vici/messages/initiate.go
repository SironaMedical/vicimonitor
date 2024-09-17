package messages

type Initiate struct {
	// child = <CHILD_SA configuration name to initiate>
	Child string `vici:"child"`
	// ike = <IKE_SA configuration name to initiate or to find child under>
	Ike string `vici:"ike"`
	// timeout = <timeout in ms before returning>
	Timeout int64 `vici:"timeout"`
	// init-limits = <whether limits may prevent initiating the CHILD_SA>
	InitLimits bool `vici:"init-limits"`
	// loglevel = <loglevel to issue "control-log" events for>
	LogLevel int64 `vici:"loglevel"`
}
