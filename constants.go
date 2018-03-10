package tradfri

// Derived from https://github.com/eclipse/smarthome/blob/4204ce06bb28c28e5f711e720f87ef83beff2e27/extensions/binding/org.eclipse.smarthome.binding.tradfri/src/main/java/org/eclipse/smarthome/binding/tradfri/TradfriBindingConstants.java

const Remote = 0
const Remote2 = 1
const Lamp = 2
const DimMax = 254
const DimMin = 0
const MiredMin = 250
const MiredMax = 454
const ColorTempColdX = 24841
const ColorTempColdY = 24593
const ColorTempDayX = 29969
const ColorTempDayY = 26804
const ColorTempWarmX = 32977
const ColorTempWarmY = 27105
const ColorTempCold = "f5faf6"
const ColorTempDay = "f1e0b5"
const ColorTempWarm = "efd275"

const tradfriPort = 5684
const preauthIdentity = "Client_identity"
const (
	uriDevices             = "/15001"
	uriGroups              = "/15004"
	uriIdent               = "/15011/9063"
	uriGatewayInfo         = "/15011/15012"
	uriGatewayReboot       = "/15011/9030"
	uriGatewayFactoryReset = "/15011/9031"
)
