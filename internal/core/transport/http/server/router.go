package core_http_server

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
	ApiVersion3 = ApiVersion("v3")
)

type APIVersionRouter struct {
	apiVersion ApiVersion
	routes     []Route
}

func NewAPIVersionRouter(apiVersion ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		apiVersion: apiVersion,
		routes:     make([]Route, 0),
	}
}

func (v *APIVersionRouter) RegisterRoutes(routes ...Route) {
	v.routes = append(v.routes, routes...)
}
