package rest

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mtdx/assert"
	"github.com/mtdx/keyc/common"
	"github.com/mtdx/keyc/config"
	"github.com/mtdx/keyc/labels"
	"github.com/mtdx/keyc/steam"
	"github.com/mtdx/keyc/validator"
)

func tradeoffersCheck(t *testing.T) {
	tradeofferreq1 := &steam.TradeoffersRequest{
		TradeofferID:    testTradeOfferID1,
		SteamID:         testSteamID,
		Type:            labels.BUYFROMUS,
		MerchantSteamID: testSteamID,
		Amount:          2,
		AppID:           labels.CSGO_APP_ID,
	}
	tradeofferreq2 := &steam.TradeoffersRequest{
		TradeofferID:    testTradeOfferID2,
		SteamID:         testSteamID,
		Type:            labels.BUYFROMUS,
		MerchantSteamID: testSteamID,
		Amount:          4,
		AppID:           labels.CSGO_APP_ID,
	}
	jsonreq, _ = json.Marshal(tradeofferreq1)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/tradeoffers?key="+config.SteamBotsAPIKey(), bytes.NewReader(jsonreq), jwt)
	jsonreq, _ = json.Marshal(tradeofferreq2)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/tradeoffers?key="+config.SteamBotsAPIKey(), bytes.NewReader(jsonreq), jwt)

	var successResp common.SuccessResponse
	if err := json.Unmarshal([]byte(body), &successResp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}
	assert.Equal(t, successResp.StatusText, "Tradeoffer has been created", successResp.StatusText+", "+successResp.SuccessText)

	jsonreq, _ = json.Marshal(tradeofferreq2)
	_, body = callEndpoint(t, ts, "POST", "/api/v1/tradeoffers", bytes.NewReader(jsonreq), jwt)
	var errResp common.ErrResponse
	if err := json.Unmarshal([]byte(body), &errResp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}
	assert.Equal(t, errResp.StatusText, "Invalid request", errResp.StatusText)

	_, body = callEndpoint(t, ts, "GET", "/api/v1/tradeoffers", nil, jwt)
	tradeoffersresp := make([]steam.TradeoffersResponse, 2)
	if err := json.Unmarshal([]byte(body), &tradeoffersresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, tradeoffer := range tradeoffersresp {
		if err := validator.Validate(tradeoffer); err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
	}

	if assert.Equal(t, 2, len(tradeoffersresp), body) {
		assert.Equal(t, labels.ACTIVE, int(tradeoffersresp[0].Status), body)
		assert.Equal(t, labels.BUYFROMUS, int(tradeoffersresp[1].Type), body)
		assert.Equal(t, 4, int(tradeoffersresp[0].Amount), body)
	}

	// validate update next
	tradeofferupdatereq1 := &steam.TradeoffersUpdateRequest{
		Status:         labels.ACCEPTED,
		FailureDetails: "",
	}
	tradeofferupdatereq2 := &steam.TradeoffersUpdateRequest{
		Status:         labels.DECLINED,
		FailureDetails: "Declined",
	}
	jsonreq, _ = json.Marshal(tradeofferupdatereq1)
	_, body = callEndpoint(t, ts, "PUT", "/api/v1/tradeoffers/"+testTradeOfferID1+"?key="+config.SteamBotsAPIKey(), bytes.NewReader(jsonreq), jwt)
	jsonreq, _ = json.Marshal(tradeofferupdatereq2)
	_, body = callEndpoint(t, ts, "PUT", "/api/v1/tradeoffers/"+testTradeOfferID2+"?key="+config.SteamBotsAPIKey(), bytes.NewReader(jsonreq), jwt)

	var successUpdateResp common.SuccessResponse
	if err := json.Unmarshal([]byte(body), &successUpdateResp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}
	assert.Equal(t, successUpdateResp.StatusText, "Tradeoffer has been updated", successUpdateResp.StatusText+", "+successUpdateResp.SuccessText)

	_, body = callEndpoint(t, ts, "GET", "/api/v1/tradeoffers", nil, jwt)
	tradeoffersresp = make([]steam.TradeoffersResponse, 2)
	if err := json.Unmarshal([]byte(body), &tradeoffersresp); err != nil {
		t.Fatalf("Failed to Unmarshal, got: %s, error: %s", body, err.Error())
	}

	for _, tradeoffer := range tradeoffersresp {
		if err := validator.Validate(tradeoffer); err != nil {
			t.Fatalf("Error: %s", err.Error())
		}
	}

	if assert.Equal(t, 2, len(tradeoffersresp), body) {
		assert.Equal(t, labels.ACCEPTED, int(tradeoffersresp[1].Status), body)
		assert.Equal(t, labels.DECLINED, int(tradeoffersresp[0].Status), body)
		assert.Equal(t, "Declined", tradeoffersresp[0].FailureDetails.String, body)
	}
}
