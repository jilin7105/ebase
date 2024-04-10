package ebase

import (
	"github.com/jilin7105/ebase/util/LinkTracking"
)

func (eb *Eb) initLinkTracking() {
	Producer := GetKafka(eb.Config.LinkTrack.KafkaProducerName)
	err := LinkTracking.InitLinkTracking(eb.Config, Producer)
	if err != nil {
		return
	}
}
