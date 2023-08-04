package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"github.com/bsm/openrtb/v3"
)

var predefinedCampaigns = map[string]openrtb.Bid{
	"imp_id_1": {
		ID:         "predefined-bid-id-1",
		ImpID:      "imp_id_1",
		Price:      0.25,
		CreativeID: "creative-123",
		Width:      300,
		Height:     250,
		AdMarkup:   `<html><a href="//example.com"><img src="//example.com/ad.img"></a></html>`,
		AdID:       "ad-12345",
		AdvDomains: []string{"example.com"},
		Ext:        json.RawMessage(`{}`),
		Categories: []openrtb.ContentCategory{"IAB1-1", "IAB1-2"},
	},
	"imp_id_2": {
		ID:         "predefined-bid-id-2",
		ImpID:      "imp_id_2",
		Price:      0.15,
		CreativeID: "creative-456",
		Width:      728,
		Height:     90,
		AdMarkup:   `<html><a href="//example.com"><img src="//example.com/ad2.img"></a></html>`,
		AdID:       "ad-67890",
		AdvDomains: []string{"example.com"},
		Ext:        json.RawMessage(`{}`),
		Categories: []openrtb.ContentCategory{"IAB1-2", "IAB-4"},
	},
	"imp_id_3": {
		ID:         "predefined-bid-id-3",
		ImpID:      "imp_id_3",
		Price:      0.12,
		CreativeID: "creative-789",
		Width:      160,
		Height:     600,
		AdMarkup:   `<html><a href="//example.com"><img src="//example.com/ad3.img"></a></html>`,
		AdID:       "ad-45678",
		AdvDomains: []string{"example.com"},
		Ext:        json.RawMessage(`{}`),
	},
	"imp_id_4": {
		ID:         "predefined-bid-id-4",
		ImpID:      "imp_id_4",
		Price:      0.18,
		CreativeID: "creative-890",
		Width:      300,
		Height:     600,
		AdMarkup:   `<html><a href="//example.com"><img src="//example.com/ad4.img"></a></html>`,
		AdID:       "ad-23456",
		AdvDomains: []string{"example.com"},
		Ext:        json.RawMessage(`{}`),
	},
	"imp_id_5": {
		ID:         "predefined-bid-id-5",
		ImpID:      "imp_id_5",
		Price:      0.20,
		CreativeID: "creative-567",
		Width:      320,
		Height:     50,
		AdMarkup:   `<html><a href="//example.com"><img src="//example.com/ad5.img"></a></html>`,
		AdID:       "ad-78901",
		AdvDomains: []string{"example.com"},
		Ext:        json.RawMessage(`{}`),
	},
}

func BidRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method. Only POST requests are allowed.", http.StatusMethodNotAllowed)
	}

	var (
		ctx         context.Context
		bidResponse *openrtb.BidResponse
		bidRequest  *openrtb.BidRequest
	)

	bidRequest = new(openrtb.BidRequest)
	if err := json.NewDecoder(r.Body).Decode(bidRequest); err != nil {
		log.Printf("Json decode error: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Json decode error"))
		return
	}
	ctx = context.Background()

	if len(bidRequest.Impressions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Imp field is empty"))
		return
	}

	imp1 := bidRequest.Impressions[0]

	ext1, _ := json.Marshal(map[string]interface{}{})
	var bid1 openrtb.Bid
	if predefinedBid, found := predefinedCampaigns[imp1.ID]; found {
		bid1 = predefinedBid
	} else {
		bid1 = openrtb.Bid{
			ID:         "test-bid-id-1",
			ImpID:      imp1.ID,
			Price:      0.1,
			CreativeID: "test-creative-id-1",
			Width:      720,
			Height:     80,
			AdMarkup:   `<html><a href="//localhost"><img src="//localhost/ad.img"></a></html>`,
			AdID:       "test-ad-id-12345",
			AdvDomains: []string{"example.com"},
			Ext:        ext1,
		}
	}

	switch {
	case imp1.Banner != nil:
		if formats := imp1.Banner.Formats; len(formats) > 0 {
			bid1.Width = formats[0].Width
			bid1.Height = formats[0].Height
		}

	case imp1.Video != nil:
		if bid1.Width, bid1.Height = imp1.Video.Width, imp1.Video.Height; bid1.Width == 0 || bid1.Height == 0 {
			bid1.Width = bidRequest.Device.Width
			bid1.Height = bidRequest.Device.Height
		}
		bid1.AdMarkup = `<html><a href="//localhost"><img src="//localhost/ad.img"></a></html>`

	case imp1.Audio != nil:
	case imp1.Native != nil:
	}

	generatedBids := []openrtb.Bid{
		bid1,
	}

	var elibibleBids []openrtb.Bid

	for _, generatedBid := range generatedBids {
		for _, blockedCategory := range bidRequest.BlockedCategories {
			if containsCategory(generatedBid.Categories, blockedCategory) {
				continue
			}
		}
		elibibleBids = append(elibibleBids, generatedBid)
	}

	bidResponse = &openrtb.BidResponse{
		ID: bidRequest.ID,
		SeatBids: []openrtb.SeatBid{
			{
				Seat:  "Bidder",
				Group: 0,
				Bids:  elibibleBids,
			},
		},
		BidID:      "TEST_BID_ID",
		Currency:   "USD",
		CustomData: "",
		NBR:        0,
		Ext:        json.RawMessage(`{}`),
	}

	if content, err := json.Marshal(bidResponse); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(content)
	} else {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	ctx.Done()

}

func containsCategory(categories []openrtb.ContentCategory, blockedCategory openrtb.ContentCategory) bool {
	for _, category := range categories {
		if category == blockedCategory {
			return true
		}
	}
	return false
}
