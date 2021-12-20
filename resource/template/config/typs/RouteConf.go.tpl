package typs

type RouteConf struct {
	Path        string                 `json:"path"`
	Method      string                 `json:"method"`
	Handler     string                 `json:"handler"`
	Middlewares *RouteConf_Middlewares `json:"middlewares"`
	Routes      []*RouteConf           `json:"routes"`
}

type RouteConf_Middlewares struct {
	Fronts []string `json:"fronts"`
	Posts  []string `json:"posts"`
}

func (i *RouteConf) ConfigName() string {
	return "RouteConf"
}