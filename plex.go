/*
Copyright 2017 Will Boyce

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type discoverResponse struct {
	FriendlyName    string
	ModelNumber     string
	FirmwareName    string
	TunerCount      uint
	FirmwareVersion string
	DeviceID        string
	DeviceAuth      string
	BaseURL         string
	LineupURL       string
}

type lineupResponse struct {
	GuideNumber string
	GuideName   string
	URL         string
}

type lineupStatusResponse struct {
	ScanInProgress uint
	ScanPossible   uint
	Source         string
	SourceList     []string
}

func (p *plexHeadend) discoverHandler(w http.ResponseWriter, r *http.Request) {
	data := discoverResponse{
		p.name,
		"HDHR4-2DT",
		"hdhomerun4_dvbt",
		uint(p.tuners),
		"20160629",
		p.deviceID,
		"",
		p.pxyBaseURL,
		fmt.Sprintf("%s/lineup.json", p.pxyBaseURL),
	}
	json.NewEncoder(w).Encode(data)
}

func (p *plexHeadend) lineupHandler(w http.ResponseWriter, r *http.Request) {
	sliceContains := func(a []string, b string) bool {
		for _, s := range a {
			if b == s {
				return true
			}
		}
		return false
	}

	data := make([]lineupResponse, 0)
	for _, channel := range p.tvhGetChannels() {
		if p.tag == "" || sliceContains(channel.Tags, p.tag) {
			data = append(data,
				lineupResponse{fmt.Sprintf("%d", channel.Number), channel.Name, channel.URL})
		}
	}
	json.NewEncoder(w).Encode(data)
}

func (p *plexHeadend) lineupStatusHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(lineupStatusResponse{
		0, 1, "Cable", []string{"Cable"}})
}

func (p *plexHeadend) lineupPostHandler(w http.ResponseWriter, r *http.Request) {}
