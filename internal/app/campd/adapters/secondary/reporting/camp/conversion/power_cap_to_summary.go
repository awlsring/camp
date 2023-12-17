package conversion

import (
	"github.com/awlsring/camp/internal/pkg/domain/machine"
	local "github.com/awlsring/camp/pkg/gen/local_grpc"
	grpcmodel "github.com/awlsring/camp/pkg/grpc_model"
)

func PowerCapabilitiesFromDomain(capabilities *machine.PowerCapabilities) *local.ReportedPowerCapabilitiesSummary {
	sum := &local.ReportedPowerCapabilitiesSummary{
		WakeOnLan: &local.MachinePowerCapabilityWakeOnLanSummary{
			Enabled: false,
		},
		PowerOff: &local.MachinePowerCapabilityPowerOffSummary{
			Enabled: false,
		},
		Reboot: &local.MachinePowerCapabilityRebootSummary{
			Enabled: false,
		},
	}
	if capabilities != nil {
		sum.PowerOff.Enabled = capabilities.PowerOff.Enabled
		sum.Reboot.Enabled = capabilities.Reboot.Enabled
		if capabilities.WakeOnLan.Enabled && capabilities.WakeOnLan.MacAddress != nil {
			sum.WakeOnLan.Enabled = capabilities.WakeOnLan.Enabled
			mac := capabilities.WakeOnLan.MacAddress.String()
			sum.WakeOnLan.MacAddress = grpcmodel.NewStringValue(&mac)
		}
	}

	return sum
}
