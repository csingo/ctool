package vars

import "framework/config/typs"

var Route = &typs.RouteConf{Middlewares: &typs.RouteConf_Middlewares{Fronts: []string{"Test"}}, Routes: []*typs.RouteConf{
	{Path: "/ping", Method: "GET", Handler: "app/controller/HomeController::Ping"},
	{Path: "/hello", Method: "GET", Handler: "app/controller/HomeController::Hello"},
}}
