package api

import (
	"encoding/json"
	"fmt"
	"github.com/febytanzil/dockerapp/services/route"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/url"
)

type GetResponse struct {
	Status   string     `json:"status"`
	Error    string     `json:"error,omitempty"`
	Distance int        `json:"total_distance,omitempty"`
	Duration float64    `json:"total_time,omitempty"`
	Path     [][]string `json:"path,omitempty"`
}

func GetRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	raw, ok := vars["token"]
	w.Header().Set("Content-Type", "application/json")

	if !ok {
		res := &GetResponse{
			Error: "failed to read token",
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
		return
	}
	token, err := url.PathUnescape(raw)
	if nil != err {
		log.Println("error decode token", err)
		res := &GetResponse{
			Error: "failed to decode token",
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
		return
	}
	result, err := route.GetInstance().GetShortestRoute(token)
	if nil != err {
		res := &GetResponse{}
		switch err {
		case route.ErrProgress:
			res.Status = "in progress"
		case route.ErrCalculate:
			res.Status = "failure"
			res.Error = "error calculate path"
		default:
			res.Status = "failure"
			res.Error = err.Error()
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
	} else {
		res := &GetResponse{
			Status: "failure",
			Error:  "error data",
		}
		if nil == result {
			resBytes, _ := json.Marshal(res)
			w.WriteHeader(http.StatusOK)
			w.Write(resBytes)
			return
		}
		res = &GetResponse{
			Status:   "success",
			Distance: result.TotalDistance,
			Duration: result.TotalTime,
		}
		for _, v := range result.Path {
			res.Path = append(res.Path, []string{fmt.Sprintf("%.6f", v.Lat), fmt.Sprintf("%.6f", v.Lng)})
		}
		resBytes, _ := json.Marshal(res)
		w.WriteHeader(http.StatusOK)
		w.Write(resBytes)
	}
}
