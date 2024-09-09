package messages

type ControlLog struct {
	Group       string `vici:"group"`
	Level       string `vici:"level"`
	IkeSAName   string `vici:"ikesa-name"`
	IkeSAUiqued string `vici:"ikesa-uniqued"`
	Message     string `vici:"msg"`
}
