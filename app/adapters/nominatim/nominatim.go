package nominatim

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

type LookupResponse struct {
	PlaceId     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmId       int      `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	PlaceRank   int      `json:"place_rank"`
	Importance  float64  `json:"importance"`
	Addresstype string   `json:"addresstype"`
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Address     Address  `json:"address"`
	Boundingbox []string `json:"boundingbox"`
}

type Address struct {
	Amenity      string `json:"amenity"`
	HouseNumber  string `json:"house_number"`
	Road         string `json:"road"`
	Suburb       string `json:"suburb"`
	Borough      string `json:"borough"`
	City         string `json:"city"`
	ISO31662Lvl4 string `json:"ISO3166-2-lvl4"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
	CountryCode  string `json:"country_code"`
}

type Client struct {
	log       *zap.Logger
	client    *http.Client
	rateLimit *rate.Limiter

	lookupCache map[string]LookupResponse
}

func New() *Client {
	client := &http.Client{
		Transport: &http.Transport{},
	}

	return &Client{
		log:    zap.L().Named("nominatim"),
		client: client,
		// 1 call per second: https://operations.osmfoundation.org/policies/nominatim/
		rateLimit:   rate.NewLimiter(1, 1),
		lookupCache: make(map[string]LookupResponse),
	}
}

func (n *Client) get(ctx context.Context, url string) (resp *http.Response, err error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("User-Agent", "raumzeitalpaka")
	r = r.WithContext(ctx)

	err = n.rateLimit.Wait(ctx)
	if err != nil {
		return nil, err
	}
	n.log.Info("sending request to nominatim api", zap.String("url", r.URL.String()))
	return n.client.Do(r)
}

func (n *Client) Lookup(ctx context.Context, osmId string) (response LookupResponse, err error) {
	if osmId == "" {
		return response, fmt.Errorf("invalid osmId")
	}

	response, exists := n.lookupCache[osmId]
	if exists {
		n.log.Debug("using cached lookup response", zap.String("osmId", osmId))
		return response, nil
	}

	resp, err := n.get(ctx, fmt.Sprintf(`https://nominatim.openstreetmap.org/lookup?osm_ids=%v&format=json`, osmId))
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	var allResponses []LookupResponse
	err = json.NewDecoder(resp.Body).Decode(&allResponses)
	response = allResponses[0]
	n.lookupCache[osmId] = response
	return response, err
}
