package elog

const (
	ClientAddress = "client.address"
	ClientBytes   = "client.bytes"
	ClientDomain  = "client.domain"
	ClientIP      = "client.ip"

	ContainerID        = "container.id"
	ContainerImageName = "container.image.name"
	ContainerName      = "container.name"

	DestinationAddress = "destination.address"
	DestinationBytes   = "destination.bytes"
	DestinationDomain  = "destination.domain"
	DestinationIP      = "destination.ip"

	DeviceID              = "device.id"
	DeviceModelIdentifier = "device.model.identifier"
	DeviceModelName       = "device.model.name"

	ECSVersion = "ecs.version"

	DeviceEmailFromAddress = "email.from.address"
	DeviceEmailToAddress   = "email.to.address"

	ErrorCode       = "error.code"
	ErrorData       = "error.data"
	ErrorMessage    = "error.message"
	ErrorStackTrace = "error.stack_trace"

	EventAction   = "event.action"
	EventCode     = "event.code"
	EventCreated  = "event.created"
	EventDuration = "event.duration"
	EventEnd      = "event.end"
	EventID       = "event.id"
	EventOriginal = "event.original"
	EventStart    = "event.start"
	EventType     = "event.type"

	GeoCityName       = "geo.city_name"
	GeoCountryISOCode = "geo.country_iso_code"
	GeoCountryName    = "geo.country_name"
	GeoLocation       = "geo.location"
	GeoName           = "geo.name"
	GeoTimezone       = "geo.timezone"

	HTTPRequestBodyContent  = "http.request.body.content"
	HTTPRequestBytes        = "http.request.bytes"
	HTTPRequestID           = "http.request.id"
	HTTPRequestMethod       = "http.request.method"
	HTTPResponseBodyContent = "http.response.body.content"
	HTTPResponseBytes       = "http.response.bytes"
	HTTPResponseStatusCode  = "http.response.status_code"
	HTTPVersion             = "http.version"

	LevelField = "log.level"

	ServerAddress = "server.address"
	ServerDomain  = "server.domain"
	ServerIP      = "server.ip"

	ServiceEnvironment = "service.environment"
	ServiceName        = "service.name"
	ServiceTargetName  = "service.target.name"
	ServiceVersion     = "service.version"

	SpanID        = "span.id"
	TraceID       = "trace.id"
	TransactionID = "transaction.id"

	URLOrigin = "url.origin"
	URLPath   = "url.path"

	UserID    = "user.id"
	UserEmail = "user.email"
	UserName  = "user.name"
	UserRoles = "user.roles"

	UserAgentName     = "user_agent.name"
	UserAgentOriginal = "user_agent.original"
	UserAgentVersion  = "user_agent.version"
)
