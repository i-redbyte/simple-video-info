package model

import "time"

type Info struct {
	Video Video `json:"video"`
	Audio Audio `json:"audio"`
}

type Video struct {
	Name     string        `json:"name"`
	Width    int           `json:"width"`
	Height   int           `json:"height"`
	BitRate  int64         `json:"bitRate"`
	Duration time.Duration `json:"duration"`
}

type Audio struct {
	Name     string        `json:"name"`
	BitRate  int64         `json:"bitRate"`
	Duration time.Duration `json:"duration"`
}
