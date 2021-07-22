package memory

import (
	"ARPSpoofing/models"
)

var DataCh chan *models.HTTPPacket = make(chan *models.HTTPPacket, 2048)
