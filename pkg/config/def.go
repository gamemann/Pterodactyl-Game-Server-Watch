package config

// Server struct used for each server config.
type Server struct {
	Name         string `json:"name"`
	Enable       bool   `json:"enable"`
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	UID          string `json:"uid"`
	ScanTime     int    `json:"scantime"`
	MaxFails     int    `json:"maxfails"`
	MaxRestarts  int    `json:"maxrestarts"`
	RestartInt   int    `json:"restartint"`
	ReportOnly   bool   `json:"reportonly"`
	A2STimeout   int    `json:"a2stimeout"`
	RconPassword string `json:"rconpassword"`
	Mentions     string `json:"mentions"`
	ViaAPI       bool
	Delete       bool
}

// Misc options.
type Misc struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Config struct used for the general config.
type Config struct {
	APIURL          string   `json:"apiurl"`
	Token           string   `json:"token"`
	AppToken        string   `json:"apptoken"`
	AddServers      bool     `json:"addservers"`
	DebugLevel      int      `json:"debug"`
	ReloadTime      int      `json:"reloadtime"`
	DefEnable       bool     `json:"defenable"`
	DefScanTime     int      `json:"defscantime"`
	DefMaxFails     int      `json:"defmaxfails"`
	DefMaxRestarts  int      `json:"defmaxrestarts"`
	DefRestartInt   int      `json:"defrestartint"`
	DefReportOnly   bool     `json:"defreportonly"`
	DefA2STimeout   int      `json:"defa2stimeout"`
	DefRconPassword string   `json:"defrconpassword"`
	DefMentions     string   `json:"defmentions"`
	Servers         []Server `json:"servers"`
	Misc            []Misc   `json:"misc"`
	ConfLoc         string
}
