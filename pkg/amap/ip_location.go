package amap

import (
	"encoding/json"
	"net/url"
)

// IPLocationResponse ip location response
type IPLocationResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	AdCode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

// IPLocation 获取ip定位
// see https://lbs.amap.com/api/webservice/guide/api/ipconfig
func (sf Amap) IPLocation(ip string) (*IPLocationResponse, error) {
	query := make(url.Values)
	query.Add("ip", ip)
	query.Add("key", sf.Key)

	resp, err := sf.httpc.Get("https://restapi.amap.com/v3/ip?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rsp := &IPLocationResponse{}
	if err = json.NewDecoder(resp.Body).Decode(rsp); err != nil {
		return nil, err
	}
	return rsp, nil
}
