package aadhaar

import (
	"bytes"
	"context"
	"time"

	// "crypto/sha256"
	// "encoding/base64"
	// "encoding/hex"
	"adhar-verification/database"
	"adhar-verification/helper"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VerifyOTPRequest struct {
	Entity      string `json:"@entity"`
	ReferenceID string `json:"reference_id"`
	OTP         string `json:"otp"`
}

type Address struct {
	Entity      string `json:"@entity"`
	Country     string `json:"country"`
	District    string `json:"district"`
	House       string `json:"house"`
	Landmark    string `json:"landmark"`
	Pincode     int    `json:"pincode"`
	PostOffice  string `json:"post_office"`
	State       string `json:"state"`
	Street      string `json:"street"`
	Subdistrict string `json:"subdistrict"`
	Vtc         string `json:"vtc"`
}

type Kyc struct {
	Entity      string  `json:"@entity"`
	Address     Address `json:"address"`
	CareOf      string  `json:"care_of"`
	DateOfBirth string  `json:"date_of_birth"`
	YOB         int     `json:"year_of_birth"`
	EmailHash   string  `json:"email_hash"`
	FullAddress string  `json:"full_address"`
	Gender      string  `json:"gender"`
	Message     string  `json:"message"`
	MobileHash  string  `json:"mobile_hash"`
	Name        string  `json:"name"`
	Photo       string  `json:"photo"`
	Status      string  `json:"status"`
	ShareCode   string  `json:"share_code"`
}

type APIResponseVerify struct {
	Timestamp     int64  `json:"timestamp"`
	TransactionID string `json:"transaction_id"`
	Data          Kyc    `json:"data"`
	Code          int    `json:"code"`
}

func Verify(c *gin.Context) {
	verifyRequest := VerifyOTPRequest{}
	verifyRequest.Entity = "in.co.sandbox.kyc.aadhaar.okyc.request"
	aadhaarNo := c.Param("aadhaar_no")
	if err := c.Bind(&verifyRequest); err != nil {
		c.Error(helper.New(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	if len(verifyRequest.OTP) == 0 {
		c.Error(helper.New(http.StatusBadRequest, "OTP is required", fmt.Errorf("missing OTP")))
		return
	}
	if len(verifyRequest.ReferenceID) == 0 {
		c.Error(helper.New(http.StatusBadRequest, "Reference ID is required", fmt.Errorf("missing reference ID")))
		return
	}

	payload, err := json.Marshal(verifyRequest)
	if err != nil {
		c.Error(helper.New(http.StatusInternalServerError, "Failed to marshal payload", err))
		return
	}
	url := "https://api.sandbox.co.in/kyc/aadhaar/okyc/otp/verify"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		c.Error(helper.New(http.StatusInternalServerError, "Failed to create request", err))
		return
	}

	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		c.Error(helper.New(http.StatusUnauthorized, "Authorization header is missing", nil))
		return
	}
	token := strings.TrimPrefix(authToken, "Bearer ")
	if token == "" {
		c.Error(helper.New(http.StatusUnauthorized, "Invalid Authorization token format", nil))
		return
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-api-version", "2.0")
	req.Header.Add("content-type", "application/json")
	req.Header.Set("Authorization", token)
	req.Header.Add("x-api-key", os.Getenv("SANDBOX_API_KEY"))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.Error(helper.New(http.StatusInternalServerError, "Failed to call sandbox API", err))
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.Error(helper.New(http.StatusInternalServerError, "Failed to read response", err))
		return
	}
	apiResponse := APIResponseVerify{}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		c.Error(helper.New(http.StatusInternalServerError, "Failed to unmarshal API response", err))
		return
	}
	update(aadhaarNo, apiResponse)

	c.JSON(res.StatusCode, gin.H{
		"message":  "OTP Verification Response",
		"response": apiResponse,
	})
}

func update(adhar_no string, data APIResponseVerify) error {
	updateData := bson.M{
		"$set": bson.M{
			"name":          data.Data.Name,
			"gender":        data.Data.Gender,
			"address":       data.Data.Address,
			"full_address":  data.Data.FullAddress,
			"care_of":       data.Data.CareOf,
			"date_of_birth": data.Data.DateOfBirth,
			"share_code":    data.Data.ShareCode,
			"status":        data.Data.Status,
			"year_of_birth": data.Data.YOB,
			"updated_at":    time.Now(),
		},
	}

	updateResult, err := database.DB.Collection("aadhaar_data").UpdateOne(
		context.TODO(),
		bson.M{"aadhaar_no": adhar_no},
		updateData,
	)

	if err != nil {
		return fmt.Errorf("failed to update Aadhaar data: %v", err)
	}

	if updateResult.MatchedCount == 0 {
		return fmt.Errorf("no document found with AadhaarNo: %s", adhar_no)
	}

	return nil
}
