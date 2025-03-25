package ciscoaciaim

import (
	"github.com/openstack-k8s-operators/lib-common/modules/storage"
)

const (
	// ServiceName -
	ServiceName = "ciscoaciaim"

	// HorizonPortName -
	CiscoAciAimPortName = "CiscoAciAim"

	// CiscoAciAimExtraVolTypeUndefined can be used to label an extraMount which is
	// not associated to anything in particular
	CiscoAciAimExtraVolTypeUndefined storage.ExtraVolType = "Undefined"
	// CiscoAciAim is the global ServiceType that refers to all the components deployed
	// by the horizon-operator
	CiscoAciAim storage.PropagationType = "CiscoAciAim"
)

// CisciAciAimPropagation is the  definition of the CiscoAciAim propagation service
var CiscoAciAimPropagation = []storage.PropagationType{CiscoAciAim}
