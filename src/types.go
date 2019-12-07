package main

// Stat storing stat data
type Stat struct {
	OSRelease string  `json:"os_release"` //
	Uptime    uint64  `json:"uptime"`     //
	Load      string  `json:"load"`       //
	CPURate   float64 `json:"cpu_rate"`   //
	MemRate   float64 `json:"mem_rate"`   //
	NetRead   uint64  `json:"net_read"`   //
	NetWrite  uint64  `json:"net_write"`  //
	Version   string  `json:"version"`
	Ip        string  `json:"ip"`
	Id        int     `json:"id"`
}
type MasterNodeHeight struct {
	Result struct {
		Height   int `json:"height:"`
		Producer int `json:"producer:"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    int         `json:"id"`
}
type MasterNode struct {
	TrxHash        string
	NodeStatus     string
	Id             string
	BlockNum       string
	ExpiryProducer string
	IsselfProblock string
}
type Produce struct {
	Height   int `json:"height:"`
	Producer int `json:"producer:"`
}
type Version struct {
	Code    int    `json:"code"`
	Version string `json:"version"`
}
type StatusInfo struct {
	Vin             string `json:"vin"`
	Service         string `json:"service"`
	Publickey       string `json:"publickey"`
	Payee           string `json:"payee"`
	LicenseVersion  int    `json:"license version"`
	LicensePeriod   string `json:"license period"`
	LicenseData     string `json:"license data"`
	LicenseStatus   string `json:"license status"`
	Validkey        string `json:"validkey"`
	Masternodeindex int    `json:"masternodeindex"`
	Status          string `json:"status"`
}
