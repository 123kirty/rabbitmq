package eatask

import (
	"encoding/json"

	"gorm.io/gorm"
)

type (
	Hotel struct {
		gorm.Model

		HotelId string `json:"hotel_id"`
		Name    string `json:"name"`

		Country string `json:"country"`
		Address string `json:"address"`

		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`

		Telephone string `json:"telephone"`

		Amenities   []string `json:"amenities" gorm:"-"`
		AmenitiesDb string   `json:"-"`

		Description string `json:"description"`

		RoomCount int    `json:"room_count"`
		Currency  string `json:"currency"`
	}

	Room struct {
		gorm.Model

		HotelId    uint   `json:"-"`
		HotelIdReq string `json:"hotel_id" gorm:"-"`
		RoomId     string `json:"room_id"`

		Description string `json:"description"`
		Name        string `json:"name"`

		Capacity struct {
			MaxAdults     int `json:"max_adults"`
			ExtraChildren int `json:"extra_children"`
		} `json:"capacity" gorm:"-"`

		CapacityDb string `json:"-"`
	}

	RatePlan struct {
		gorm.Model

		HotelId    uint   `json:"-"`
		HotelIdReq string `json:"hotel_id" gorm:"-"`
		RatePlanId string `json:"rate_plan_id"`

		CancellationPolicy []struct {
			Type              string `json:"type"`
			ExpiresDaysBefore int    `json:"expires_days_before"`
		} `json:"cancellation_policy" gorm:"-"`
		CancellationPolicyDb string `json:"-"`

		Name string `json:"name"`

		OtherConditions   []string `json:"other_conditions" gorm:"-"`
		OtherConditionsDb string   `json:"-"`
	}
)

func (hotel *Hotel) BeforeCreate(tx *gorm.DB) (err error) {

	var (
		amenB []byte
	)

	if amenB, err = json.Marshal(hotel.Amenities); err != nil {
		return
	}

	hotel.AmenitiesDb = string(amenB)

	return
}

func (room *Room) BeforeCreate(tx *gorm.DB) (err error) {

	var (
		capB []byte
	)

	if capB, err = json.Marshal(room.Capacity); err != nil {
		return
	}

	room.CapacityDb = string(capB)

	return
}

func (ratePlan *RatePlan) BeforeCreate(tx *gorm.DB) (err error) {

	var (
		cpB  []byte
		otcB []byte
	)

	if cpB, err = json.Marshal(ratePlan.CancellationPolicy); err != nil {
		return
	}

	if otcB, err = json.Marshal(ratePlan.OtherConditions); err != nil {
		return
	}

	ratePlan.CancellationPolicyDb = string(cpB)
	ratePlan.OtherConditionsDb = string(otcB)

	return
}
