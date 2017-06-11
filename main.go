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
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var version string

type plexHeadend struct {
	tvhBaseURL string
	pxyBaseURL string
	listenAddr string
	name       string
	deviceID   string
	tuners     int
	tag        string
}

func (p *plexHeadend) listen() {
	http.HandleFunc("/discover.json", p.discoverHandler)
	http.HandleFunc("/lineup.json", p.lineupHandler)
	http.HandleFunc("/lineup_status.json", p.lineupStatusHandler)
	http.HandleFunc("/lineup.post", p.lineupPostHandler)

	http.ListenAndServe(p.listenAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		http.DefaultServeMux.ServeHTTP(w, r)
	}))
}

func init() {
	pflag.CommandLine.BoolP("version", "", false, "Display version and exit")
	pflag.CommandLine.StringP("tvh-user", "u", "plex", "TVHeadend Username")
	pflag.CommandLine.StringP("tvh-pass", "P", "plex", "TVHeadend Password")
	pflag.CommandLine.StringP("tvh-host", "h", "localhost", "TVHeadend Host")
	pflag.CommandLine.StringP("tvh-port", "p", "9981", "TVHeadend Port")
	pflag.CommandLine.StringP("proxy-bind", "b", "", "Bind address (default all)")
	pflag.CommandLine.StringP("proxy-listen", "l", "80", "Listen port")
	pflag.CommandLine.StringP("proxy-hostname", "H", "localhost", "Hostname reported to Plex")
	pflag.CommandLine.StringP("name", "n", "plexHeadend", "Friendly name reported to Plex")
	pflag.CommandLine.StringP("device-id", "i", "1", "Device ID reported to Plex")
	pflag.CommandLine.IntP("tuners", "t", 1, "Number of Tuners reported to Plex")
	pflag.CommandLine.StringP("tag", "f", "", "TVHeadend tag to filter reported channels (default none)")

	viper.BindPFlags(pflag.CommandLine)
	viper.SetEnvPrefix("PLEXHEADEND")
	viper.BindEnv("tvh-user", "PLEXHEADEND_TVH_USER")
	viper.BindEnv("tvh-pass", "PLEXHEADEND_TVH_PASS")
	viper.BindEnv("tvh-host", "PLEXHEADEND_TVH_HOST")
	viper.BindEnv("tvh-port", "PLEXHEADEND_TVH_PORT")
	viper.BindEnv("proxy-bind", "PLEXHEADEND_PROXY_BIND")
	viper.BindEnv("proxy-listen", "PLEXHEADEND_PROXY_LISTEN")
	viper.BindEnv("proxy-hostname", "PLEXHEADEND_PROXY_HOSTNAME")
	viper.BindEnv("name") // PLEXHEADEND_NAME
	viper.BindEnv("device-id", "PLEXHEADEND_DEVICE_ID")
	viper.BindEnv("tag")    // PLEXHEADEND_TAG
	viper.BindEnv("tuners") // PLEXHEADEND_TUNERS
}

func main() {
	pflag.Parse()

	if viper.GetBool("version") {
		fmt.Printf("plexheadend v%s\n", version)
		return
	}

	p := plexHeadend{
		fmt.Sprintf("http://%s:%s@%s:%s",
			viper.GetString("tvh-user"), viper.GetString("tvh-pass"),
			viper.GetString("tvh-host"), viper.GetString("tvh-port")),
		fmt.Sprintf("http://%s:%s",
			viper.GetString("proxy-bind"), viper.GetString("proxy-listen")),
		fmt.Sprintf("%s:%s",
			viper.GetString("proxy-bind"), viper.GetString("proxy-listen")),
		viper.GetString("name"),
		viper.GetString("device-id"),
		viper.GetInt("tuners"),
		viper.GetString("tag"),
	}
	if hostname := viper.GetString("proxy-hostname"); len(hostname) > 0 {
		p.pxyBaseURL = fmt.Sprintf("http://%s:%s", hostname, viper.GetString("proxy-listen"))
	}

	safeBaseURL := fmt.Sprintf("http://%s:****@%s:%s",
		viper.GetString("tvh-user"), viper.GetString("tvh-host"), viper.GetString("tvh-port"))
	log.Printf("TvhURL:\t%s", safeBaseURL)
	log.Printf("ProxyURL:\t%s", p.pxyBaseURL)
	log.Printf("ListenAddr:\t%s", p.listenAddr)
	log.Printf("DeviceName:\t%s", p.name)
	log.Printf("DeviceID:\t%s", p.deviceID)
	log.Printf("TunerCount:\t%d", p.tuners)
	if p.tag != "" {
		log.Printf("ChannelTag:\t%s", p.tag)
	}

	p.listen()
}
