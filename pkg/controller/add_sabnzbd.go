package controller

import (
	"github.com/parflesh/sabnzbd-operator/pkg/controller/sabnzbd"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sabnzbd.Add)
}
