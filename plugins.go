package main

type VpnTwitcherPlugin interface {
	OnActive()
	OnInactive()
	Init(options map[string]string)
}

type NotificationPlugin interface {
	SendMessage(msg string) Error
	Init(options map[string]string)
}
