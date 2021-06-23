package utils

import "time"

type GCEauth struct {
	IamUser string
	UID     string
	Key1    string
	Key2    string
}

type ImgRegister struct {
	Timestamp  time.Time
	NewImgVers string
}

type Config struct {
	GCEAuthPath string
}
