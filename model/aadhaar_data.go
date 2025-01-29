package model

import "time"

type Address struct {
	Entity      string `bson:"entity,omitempty" json:"@entity,omitempty"`
	Country     string `bson:"country,omitempty" json:"country,omitempty"`
	District    string `bson:"district,omitempty" json:"district,omitempty"`
	House       string `bson:"house,omitempty" json:"house,omitempty"`
	Landmark    string `bson:"landmark,omitempty" json:"landmark,omitempty"`
	Pincode     int    `bson:"pincode,omitempty" json:"pincode,omitempty"`
	PostOffice  string `bson:"post_office,omitempty" json:"post_office,omitempty"`
	State       string `bson:"state,omitempty" json:"state,omitempty"`
	Street      string `bson:"street,omitempty" json:"street,omitempty"`
	Subdistrict string `bson:"subdistrict,omitempty" json:"subdistrict,omitempty"`
	Vtc         string `bson:"vtc,omitempty" json:"vtc,omitempty"`
}
type Aadhaar_Data struct {
	ID          string   `bson:"_id,omitempty" json:"id,omitempty"`
	AadhaarNo   string   `bson:"aadhaar_no,unique" json:"aadhaar_number,omitempty"`
	Name        *string  `bson:"name" json:"name,omitempty"`
	Gender      *string  `bson:"gender" json:"gender,omitempty"`
	Address     *Address `bson:"address" json:"address,omitempty"`
	FullAddress *string  `bson:"full_address" json:"full_address,omitempty"`
	CareOf      *string  `bson:"care_of" json:"care_of,omitempty"`
	DOB         *string  `bson:"date_of_birth" json:"date_of_birth,omitempty"`
	ShareCode   *string  `bson:"share_code" json:"share_code,omitempty"`
	Status      *string  `bson:"status" json:"status,omitempty"`
	YOB         *int     `bson:"year_of_birth" json:"year_of_birth,omitempty"`
	// Data      interface{} `bson:"data" json:"data,omitempty"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at,omitempty"`
}
