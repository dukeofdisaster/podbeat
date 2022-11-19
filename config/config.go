// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Ago    AgoType `config:"ago"`
	ApiKey string  `config:"apikey"`
	// this string will also appear in the Bearer token generated for POD access
	CustomerID string `config:"customerid"`
	// path to a sqlite db for tracking seen guids to ensure no duplicate events
	// TODO:
	//	- choose config opts for max # of guids before we roll the db or table; i.e. new_events | old_events
	Database struct {
		Path string `config:"path"`
	}
	// Write a periodic check point to disk so the utility can ensure overlap; i.e. no missed events
	CheckPoint CheckpointType `config:"checkpoint"`
	Endpoint   string         `config:"endpoint"`
	// can be one of message|maillog
	MessageType string `config:"messagetype"`
	/*
		Log      struct {
			Path string `config:"path"`
		}
	*/
	// this will probably fail to be compat with libbeat? thought output was array... maybe not
	/*
		Output struct {
			File struct {
				Path     string `config:"path"`
				Filename string `config:"filename"`
			}
		}
	*/
	// 2020 POD docs define 0400 -> 0700 as 'correct' timezones, don't recall if utc was ever used?
	Timezone string `config:"timezone"`
}
type AgoType struct {
	// one of m|h
	Units string `config:"units"`
	Value int    `config:"value"`
}
type CheckpointType struct {
	// a writeable path
	Path string `config:"path"`
	// the period in minutes at which periodic checkpoints are written to disk
	Interval int64 `config:"interval"`
	// the offset in minutes  to start the event stream from... i.e. if last ran at 13:30 and Offset is 10, then sinceTime
	// supplied to the Proofpoint API will be 13:20
	Offset int64 `config:"offset"`
}

var DefaultConfig = Config{
	ApiKey:   "aGVsbG93b3JsZAo",
	Endpoint: "ws://localhost:8080",
	CheckPoint: CheckpointType{
		Path:     "/tmp/podbeat.checkpoint",
		Interval: 15,
		Offset:   15,
	},
	MessageType: "message",
}
