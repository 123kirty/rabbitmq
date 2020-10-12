package eatask

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type (
	Offer struct {
		Hotel    Hotel    `json:"hotel"`
		Room     Room     `json:"room"`
		RatePlan RatePlan `json:"rate_plan"`
	}

	PubModel struct {
		Offers []Offer `json:"offers"`
	}
)

func (eaServer *EaServer) ProcessInputFromFile(filePath string) (err error) {

	var (
		contB []byte
		inp   PubModel
	)

	if contB, err = ioutil.ReadFile(filePath); err != nil {
		return
	}

	if err = json.Unmarshal(contB, &inp); err != nil {
		return
	}

	for _, offer := range inp.Offers {

		log.Println("Processing for hotel with ID", offer.Hotel.HotelId)

		if err = eaServer.ProcessInput(offer); err != nil {
			return
		}
	}

	return
}

func (eaServer *EaServer) ProcessInput(inp Offer) (err error) {

	var (
		hotel    Hotel
		room     Room
		ratePlan RatePlan

		result *gorm.DB
	)

	eaServer.Db.Where("hotel_id = ?", inp.Hotel.HotelId).First(&hotel)

	if hotel.ID == 0 {
		if result = eaServer.Db.Create(&inp.Hotel); result.Error != nil {
			return
		}
	}

	eaServer.Db.Where("hotel_id = ?", inp.Hotel.HotelId).First(&hotel)

	eaServer.Db.Where("room_id = ?", inp.Room.RoomId).First(&room)

	if room.ID == 0 {
		inp.Room.HotelId = hotel.ID
		if result = eaServer.Db.Create(&inp.Room); result.Error != nil {
			return
		}
	}

	eaServer.Db.Where("rate_plan_id = ?", inp.RatePlan.RatePlanId).First(&ratePlan)

	if ratePlan.ID == 0 {
		inp.RatePlan.HotelId = hotel.ID
		if result = eaServer.Db.Create(&inp.RatePlan); result.Error != nil {
			return
		}
	}

	return
}

func (eaServer *EaServer) startServer() (err error) {

	var (
		conn  *amqp.Connection
		ch    *amqp.Channel
		msgCh <-chan amqp.Delivery
	)

	if conn, err = amqp.Dial(eaServer.GetQueueUrl()); err != nil {
		return
	}

	if ch, err = conn.Channel(); err != nil {
		return
	}
	defer ch.Close()

	if msgCh, err = ch.Consume(EaTaskQueueName, "", true, false,
		false, false, nil); err != nil {

		return
	}

	go func() {
		for msg := range msgCh {
			log.Println("Received Message:", msg.Body)
			//eaServer.ProcessInput(msg.Body)
		}
	}()

	return
}
