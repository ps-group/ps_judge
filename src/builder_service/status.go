package main

type Status string

const (
	StatusPending   = "pending"
	StatusBuilding  = "building"
	StatusFailed    = "failed"
	StatusSucceed   = "succeed"
	StatusException = "exception"
)
