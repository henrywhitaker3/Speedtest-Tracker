package results

import (
	"time"
)

type Ping struct {
	Jitter  float64 `json:"jitter"`
	Latency float64 `json:"latency"`
	Low     float64 `json:"low"`
	High    float64 `json:"high"`
}

type Latency struct {
	Iqm    float64 `json:"iqm"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Jitter float64 `json:"jitter"`
}

type Measurement struct {
	Bandwidth int     `json:"bandwidth"`
	Bytes     int     `json:"bytes"`
	Elapsed   int     `json:"elapsed"`
	Latency   Latency `json:"latency"`
}

type Interface struct {
	InternalIP string `json:"internalIp"`
	Name       string `json:"name"`
	MacAddr    string `json:"macAddr"`
	IsVpn      bool   `json:"isVpn"`
	ExternalIP string `json:"externalIp"`
}

type Server struct {
	ID       int    `json:"id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Country  string `json:"country"`
	IP       string `json:"ip"`
}

type OoklaResult struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Persisted bool   `json:"persisted"`
}

// SpeedtestOutput is the raw output from the speedtest-cli
type SpeedtestOutput struct {
	Type       string       `json:"type"`
	Timestamp  time.Time    `json:"timestamp"`
	Ping       Ping         `json:"ping"`
	Download   Measurement  `json:"download"`
	Upload     Measurement  `json:"upload"`
	PacketLoss float64      `json:"packetLoss"`
	Isp        *string      `json:"isp,omitempty"`
	Interface  *Interface   `json:"interface,omitempty"`
	Server     *Server      `json:"server,omitempty"`
	Result     *OoklaResult `json:"result,omitempty"`
}

func (s *SpeedtestOutput) Sanitise() {
	s.Isp = nil
	s.Result = nil
	s.Server = nil
	s.Interface = nil
}
