package api

import "github.com/uoul/go-common/async"

type ISevenIoApi interface {
	SendSms(from string, to string, text string) chan async.ActionResult[any]
}
