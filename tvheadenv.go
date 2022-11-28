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
	"log"
	"net/http"
)

type apiTagsResponse struct {
	Entries []apiTag `json:"entries"`
}

type apiTag struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type apiChannelsResponse struct {
	Entries []apiChannel `json:"entries"`
}

type apiChannel struct {
	UUID   string   `json:"uuid"`
	Name   string   `json:"name"`
	Number int      `json:"number"`
	Tags   []string `json:"tags"`
	URL    string   `json:"url"`
}

func (p *plexHeadend) tvhGetTags() map[string]string {
	resp, err := http.Get(fmt.Sprintf("%s/api/channeltag/list", p.tvhBaseURL))
	if err != nil {
		return nil
	}

	apiResp := apiTagsResponse{}
	json.NewDecoder(resp.Body).Decode(&apiResp)
	tags := make(map[string]string)
	for _, tag := range apiResp.Entries {
		tags[tag.Key] = tag.Val
	}

	return tags
}

func (p *plexHeadend) tvhGetChannels() []apiChannel {
	resp, err := http.Get(fmt.Sprintf("%s/api/channel/grid?start=0&limit=100000", p.tvhBaseURL))
	if err != nil {
		return nil
	}

	apiResp := apiChannelsResponse{}
	json.NewDecoder(resp.Body).Decode(&apiResp)
	channels := make([]apiChannel, 0)
	tags := p.tvhGetTags()
	log.Printf("Got %d channels from tvheadend", len(apiResp.Entries))
	for _, channel := range apiResp.Entries {
		channelTags := make([]string, 0)
		for _, tag := range channel.Tags {
			tagName := tags[tag]
			if len(tagName) > 0 {
				channelTags = append(channelTags, tags[tag])
			}
		}
		channelURL := fmt.Sprintf("%s/stream/channel/%s", p.tvhBaseURL, channel.UUID)
		channels = append(channels,
			apiChannel{channel.UUID,
				channel.Name,
				channel.Number,
				channelTags,
				channelURL,
			})
	}

	return channels
}
