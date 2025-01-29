package aadhaar

import (
	"adhar-verification/database"
	"adhar-verification/helper"
	"adhar-verification/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GenerateOTPRequest struct {
	Entity        string `json:"@entity"`
	AadhaarNumber string `json:"aadhaar_number"`
	Consent       string `json:"consent"`
	Reason        string `json:"reason"`
}
type APIResponse struct {
	Timestamp     int64       `json:"timestamp"`
	TransactionID string      `json:"transaction_id"`
	Data          interface{} `json:"data"`
	Code          int         `json:"code"`
}

func GenerateOtp(c *gin.Context) {
	otpRequest := GenerateOTPRequest{
		Entity:        "in.co.sandbox.kyc.aadhaar.okyc.otp.request",
		Reason:        os.Getenv("REASON"),
		Consent:       "Y",
		AadhaarNumber: "",
	}

	if err := c.Bind(&otpRequest); err != nil {
		c.Error(helper.New(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	if len(otpRequest.AadhaarNumber) != 12 {
		c.Error(helper.New(http.StatusBadRequest, "Invalid Aadhaar number", fmt.Errorf("invalid length: %d", len(otpRequest.AadhaarNumber))))
		return
	}

	if len(otpRequest.Reason) < 20 {
		c.Error(helper.New(http.StatusBadRequest, "Reason must be at least 20 characters", fmt.Errorf("short reason: %s", otpRequest.Reason)))
		return
	}
	_, err := findOne(otpRequest.AadhaarNumber)
	if err == nil {
		createAadhaar(otpRequest.AadhaarNumber)
	}
	payload, err := json.Marshal(otpRequest)
	if err != nil {
		c.Error(helper.New(http.StatusInternalServerError, "Failed to marshal payload", err))
		return
	}

	url := "https://api.sandbox.co.in/kyc/aadhaar/okyc/otp"
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

	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		c.Error(helper.New(http.StatusInternalServerError, "Failed to unmarshal API response", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "OTP Sent Successfully",
		"response": apiResponse,
	})
}

func createAadhaar(adhar_no string) *mongo.InsertOneResult {

	newAadhaar := model.Aadhaar_Data{AadhaarNo: adhar_no}
	data, _ := database.DB.Collection("aadhaar_data").InsertOne(context.TODO(), newAadhaar)

	return data

}

func findOne(adhar_no string) (model.Aadhaar_Data, error) {
	filter := bson.M{"aadhaar_no": adhar_no}
	data := model.Aadhaar_Data{}

	err := database.DB.Collection("aadhaar_data").FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return data, nil
		}
		return data, err
	}
	return data, nil
}
