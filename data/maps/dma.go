package maps

import (
	"context"
	"errors"
	"github.com/febytanzil/dockerapp/framework/config"
	"googlemaps.github.io/maps"
	"log"
	"net/http"
	"time"
)

type DMA interface {
	GetShortestPath(origin maps.LatLng, destinations string) (*Destinations, error)
}

type RouteResponse struct {
	Pos  maps.LatLng
	Dist int
	Dur  float64
}

type Destinations struct {
	Path          []maps.LatLng `json:"path"`
	TotalDistance int           `json:"total_distance"`
	TotalTime     float64       `json:"total_time"`
}

type dmaImpl struct{}

var (
	h = &http.Client{
		Timeout: time.Second * 5,
	}
	dma DMA
	key string
)

func Init(d DMA) {
	if nil == d {
		dma = &dmaImpl{}
		key = config.Get().Maps.Key
	} else {
		dma = d
	}
}

func GetInstance() DMA {
	if nil == dma {
		log.Fatal("db is not initialized")
	}
	return dma
}

func (i *dmaImpl) GetShortestPath(origin maps.LatLng, destinations string) (*Destinations, error) {
	c, err := maps.NewClient(maps.WithAPIKey(key))
	if nil != err {
		log.Println("error get gmaps client", err)
		return nil, err
	}
	destArr, _ := maps.DecodePolyline(destinations)
	wpArr := make([]string, len(destArr))
	for _, v := range destArr {
		wpArr = append(wpArr, v.String())
	}

	routes, _, err := c.Directions(context.Background(), &maps.DirectionsRequest{
		Origin:      origin.String(),
		Destination: origin.String(),
		Waypoints:   wpArr,
		Optimize:    true,
	})
	if nil != err {
		log.Println("error calculate distance", err)
		return nil, err
	}
	if 0 == len(routes) {
		return nil, errors.New("cannot find route")
	}

	route := routes[0]
	d := new(Destinations)
	dests, _ := maps.DecodePolyline(destinations)
	d.Path = append(d.Path, origin)
	for i, v := range route.WaypointOrder {
		d.Path = append(d.Path, dests[v])
		d.TotalDistance += route.Legs[i].Meters
		d.TotalTime += route.Legs[i].Duration.Seconds()
	}

	return d, nil
}
