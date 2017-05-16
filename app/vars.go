package app

// ServiceName is an identifier-like name used anywhere this app needs to be identified.
//
// It identifies the service itself, the actual instance needs to be identified via environment
// and other details.
const ServiceName = "boilerplate.service"

// FriendlyServiceName is the visible name of the service.
const FriendlyServiceName = "Boilerplate service"

// FluentdTag is usually a static value across all instances of the same application
// as such it is set here as a constant value.
//
// Falls back to the ServiceName.
const FluentdTag string = ""
