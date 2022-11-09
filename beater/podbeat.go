package beater

/*
TODO:
	- add params to build url
	- flesh out checkpoint logic
		- if first run / checkpoint file does not exist ?
			- start with Now() and offset
*/

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/dukeofdisaster/podbeat/config"
	"github.com/gorilla/websocket"
)

func checkpointExists(c config.Config) (*bool, error) {
	var f = false
	var t = true
	_, err := os.Stat(c.CheckPoint.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &f, nil
		}
		return nil, err
	}
	return &t, nil
}

// pod docs only specify 0400 -> 0800 but let's assume the api is sane and accepts others? if not.. womp womp
func isValidTimezone(s string) bool {
	switch s {
	case "0000":
		return true
	case "0100":
		return true
	case "0200":
		return true
	case "0300":
		return true
	case "0400":
		return true
	case "0500":
		return true
	case "0600":
		return true
	case "0700":
		return true
	case "0800":
		return true
	}
	return false
}

// cap this at something sane
func isValidOffset(ago int64) bool {
	// cap at 2 days worth of minutes... I can't remember if there was a max sinceTime in POD api
	return (ago >= 0) && (ago < 2880)
}

// podbeat configuration.
type podbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of podbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &podbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts podbeat.
func (bt *podbeat) Run(b *beat.Beat) error {
	//logp.Info("podbeat is running! Hit CTRL-C to stop it.")
	logp.Info("podbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	// validate checkpoint path; if exists, try load time, report any errors, else, try create new, report errors
	checkpoint_ok, err := checkpointExists(bt.config)
	if err != nil {
		logp.Error(err)
		return err
	}
	if !*checkpoint_ok {
		logp.Info("checkpoint path does not exist, first checkpoint will be now epoch")
	}
	// validate timezione
	timezone_ok := isValidTimezone(bt.config.Timezone)
	if !timezone_ok {
		logp.Warn("got invalid timezone - may see unexpected results")
	}
	// validate offset
	offset_ok := isValidOffset(bt.config.CheckPoint.Offset)
	if !offset_ok {
		logp.Warn("expect the offset to be >= 0 and <= 2880")
		return fmt.Errorf("expect the offset to be >=0 and <= 2880")
	}

	conn, _, err := websocket.DefaultDialer.Dial(bt.config.Endpoint, nil)
	//ticker := time.NewTicker(1 * time.Second)
	if err != nil {
		logp.Error(err)
		return err
	}
	defer conn.Close()
	for {
		// jank issue? when this select statement is present we have no event reads
		/*
			select {
			case <-bt.done:
				return nil
			case <-ticker.C:
			}
		*/
		messageType, message, err := conn.ReadMessage()
		//_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			//logp.Error(err)
			return err
		}
		log.Println("message type: ", messageType)
		event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"message": string(message),
			},
		}
		bt.client.Publish(event)
	}
}

// Stop stops podbeat.
func (bt *podbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
