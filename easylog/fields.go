package easylog

const (
	ClientAddress = "client.address"
	ClientIP      = "client.ip"
	ClientBytes   = "client.bytes"

	DestinationAddress = "destination.address"
	DestinationIP      = "destination.ip"
	DestinationBytes   = "destination.bytes"

	DeviceID    = "device.id"
	DeviceModel = "device.model"

	ErrorCode       = "error.code"
	ErrorID         = "error.id"
	ErrorMessage    = "error.message"
	ErrorStackTrace = "error.stack_trace"

	EventCreated  = "event.created"
	EventDuration = "event.duration"
	EventID       = "event.id"
	EventOriginal = "event.original"
	EventType     = "event.type"

	GeoCityName       = "geo.city_name"
	GeoCountryISOCode = "geo.country_iso_code"
	GeoCountryName    = "geo.country_name"
	GeoLocation       = "geo.location"
	GeoTimezone       = "geo.timezone"

	HTTPRequestBodyContent  = "http.request.body.content"
	HTTPRequestBytes        = "http.request.bytes"
	HTTPRequestID           = "http.request.id"
	HTTPRequestMethod       = "http.request.method"
	HTTPResponseBodyContent = "http.response.body.content"
	HTTPResponseBytes       = "http.response.bytes"
	HTTPResponseStatusCode  = "http.response.status_code"

	LevelField = "log.level"

	ServiceEnvironment = "service.environment"
	ServiceName        = "service.name"
	ServiceVersion     = "service.version"

	TraceID = "trace.id"

	URLOrigin = "url.origin"
	URLPath   = "url.path"

	UserID = "user.id"

	UserAgentName     = "user_agent.name"
	UserAgentOriginal = "user_agent.original"
	UserAgentVersion  = "user_agent.version"
)
