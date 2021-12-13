// SPDX-License-Identifier: Apache-2.0

// Copyright 2021 PANTHEON.tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package descriptor

import (
	"strings"

	"github.com/go-errors/errors"
	"go.ligato.io/cn-infra/v2/logging"
	"go.ligato.io/vpp-agent/v3/pkg/models"
	"go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	kvs "go.ligato.io/vpp-agent/v3/plugins/kvscheduler/api"
	ifdescriptor "go.ligato.io/vpp-agent/v3/plugins/vpp/ifplugin/descriptor"
	interfaces "go.ligato.io/vpp-agent/v3/proto/ligato/vpp/interfaces"
	"go.pantheon.tech/stonework/plugins/isisx/descriptor/adapter"
	"go.pantheon.tech/stonework/plugins/isisx/vppcalls"
	isisx "go.pantheon.tech/stonework/proto/isisx"
)

const (
	// ISISXDescriptorName is descriptor name
	ISISXDescriptorName = "vpp-isisx"

	// dependency labels
	isisxInputInterfaceDep  = "input-interface-exists"
	isisxOutputInterfaceDep = "output-interface-exists"
)

var (
	// errEmptyInputInterface is returned when input interface is empty or blank-spaced name.
	errEmptyInputInterface = errors.New("Input interface name must be defined")
	// errEmptyOutputInterface is returned when output interface is empty or blank-spaced name.
	errEmptyOutputInterface = errors.New("Output interface name must be defined")
)

// ISISXDescriptor is descriptor for ISISXConnect
type ISISXDescriptor struct {
	log          logging.Logger
	isisxHandler vppcalls.ISISXVppAPI
}

// NewISIXDescriptor is constructor for ISISX descriptor and returns descriptor
// suitable for registration (via adapter) with the KVScheduler.
func NewISIXDescriptor(isisxHandler vppcalls.ISISXVppAPI, logger logging.PluginLogger) *api.KVDescriptor {
	ctx := &ISISXDescriptor{
		log:          logger.NewLogger("isisx-descriptor"),
		isisxHandler: isisxHandler,
	}
	typedDescr := &adapter.ISISXDescriptor{
		Name:                 ISISXDescriptorName,
		NBKeyPrefix:          isisx.ModelISISX.KeyPrefix(),
		ValueTypeName:        isisx.ModelISISX.ProtoName(),
		KeySelector:          isisx.ModelISISX.IsKeyValid,
		KeyLabel:             isisx.ModelISISX.StripKeyPrefix,
		WithMetadata:         false,
		ValueComparator:      ctx.EquivalentISIXConnections,
		Validate:             ctx.Validate,
		Create:               ctx.Create,
		Delete:               ctx.Delete,
		Retrieve:             ctx.Retrieve,
		Dependencies:         ctx.Dependencies,
		RetrieveDependencies: []string{ifdescriptor.InterfaceDescriptorName},
	}
	return adapter.NewISISXDescriptor(typedDescr)
}

// EquivalentISIXConnections compares input and output interface pairs of new and old isisx.ISISXConnection
func (d *ISISXDescriptor) EquivalentISIXConnections(key string, oldISISX, newISISX *isisx.ISISXConnection) bool {
	return oldISISX.InputInterface == newISISX.InputInterface && oldISISX.OutputInterface == newISISX.OutputInterface
}

// Validate validates VPP isisx configuration.
func (d *ISISXDescriptor) Validate(key string, isisXConnection *isisx.ISISXConnection) error {
	if strings.TrimSpace(isisXConnection.GetInputInterface()) == "" {
		return kvs.NewInvalidValueError(errEmptyInputInterface, "inputInterface")
	}
	if strings.TrimSpace(isisXConnection.GetOutputInterface()) == "" {
		return kvs.NewInvalidValueError(errEmptyOutputInterface, "outputInterface")
	}
	return nil
}

// Create creates isisXConnection using vppcalls
func (d *ISISXDescriptor) Create(key string, isisXConnection *isisx.ISISXConnection) (metadata interface{}, err error) {
	if err := d.isisxHandler.AddISISXConnection(isisXConnection.GetInputInterface(), isisXConnection.GetOutputInterface()); err != nil {
		d.log.Error(err)
		return nil, err
	}
	return nil, nil
}

// Delete removes isisXConnection using vppcalls
func (d *ISISXDescriptor) Delete(key string, isisXConnection *isisx.ISISXConnection, metadata interface{}) error {
	if err := d.isisxHandler.DeleteISISXConnection(isisXConnection.GetInputInterface(), isisXConnection.GetOutputInterface()); err != nil {
		d.log.Error(err)
		return err
	}
	return nil
}

// Retrieve returns ISISX configuration from the VPP.
func (d *ISISXDescriptor) Retrieve(correlate []adapter.ISISXKVWithMetadata) (isisXConnects []adapter.ISISXKVWithMetadata, err error) {
	// Retrieve VPP configuration.
	connections, err := d.isisxHandler.DumpISISXConnections()
	if err != nil {
		return nil, errors.Errorf("failed to dump ISISX configuration due to: %v", err)
	}

	// convert it to appropriate form
	for _, connection := range connections {
		isisXConnects = append(isisXConnects, adapter.ISISXKVWithMetadata{
			Key:    models.Key(connection),
			Value:  connection,
			Origin: api.FromNB,
		})
	}

	return isisXConnects, nil
}

// Dependencies provide list of dependencies for applying operations with ISIX
func (d *ISISXDescriptor) Dependencies(key string, isisXConnection *isisx.ISISXConnection) (dependencies []api.Dependency) {
	dependencies = append(dependencies, api.Dependency{
		Label: isisxInputInterfaceDep,
		Key:   interfaces.InterfaceKey(isisXConnection.InputInterface),
	})
	dependencies = append(dependencies, api.Dependency{
		Label: isisxOutputInterfaceDep,
		Key:   interfaces.InterfaceKey(isisXConnection.OutputInterface),
	})

	return dependencies
}
