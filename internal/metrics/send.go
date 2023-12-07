package metrics

import (
	"github.com/dkmelnik/metrics/internal/models"
	"log"
	"time"
)

func MakeRequest() {
	
}

func Send(md *models.Metrics) {
	for {
		log.Printf("%#v\n", *md)

		time.Sleep(time.Second * 10)
	}
}
