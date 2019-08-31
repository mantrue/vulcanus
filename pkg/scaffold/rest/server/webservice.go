package server

import (
	"bytes"
	"github.com/pkg/errors"
	"text/template"
)

type webServiceGenerator struct {
	cache *bytes.Buffer

	config *webServiceConfig
}

type webServiceConfig struct {
	Service *Service
	Model   *Model
}

func (g *webServiceGenerator) generate() error {

	if err := g.generateType(); err != nil {
		return errors.WithMessage(err, "generate type")
	}

	if err := g.generateWsFunc(); err != nil {
		return errors.WithMessage(err, "generate wsFunc")
	}

	if err := g.generateBusinessFunc(); err != nil {
		return errors.WithMessage(err, "generate wsFunc")
	}
	return nil

}

func (g *webServiceGenerator) generateBusinessFunc() error {

	return nil
}

func (g *webServiceGenerator) generateType() error {

	const tmplt = `

// alias the client & server communicate model
type {{.Model.Name}} = 


// {{.Service.Type}}Manager
// used to manage resource
type {{.Service.Type}} struct{

   ws *restful.WebService
   container *restful.Container
}
`

	t, err := template.New("types-tplt").Parse(tmplt)
	if err != nil {
		return errors.WithMessage(err, "parse template")
	}

	if err := t.Execute(g.cache, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}
	return nil
}

func (g *webServiceGenerator) generateWsFunc() error {

	const tmplt = `
func (s *{{.Service.Type}}) measureTime(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	now := time.Now()
	chain.ProcessFilter(req, resp)
	time.Now().Sub(now)
}

func (s *{{.Service.Type}}) intallWebService(){
	ws := new(restful.WebService)
	ws.
		Path("{{.Service.RootURLPrefix}}").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) 

	tags := []string{"{{.Service.Tag.Name}}"}

	ws.Route(ws.GET("/").To(s.getAll).
		// docs
		Doc("get all {{.Service.Kind}}").
		// spec a useful filter 
		Filter(s.measureTime).
		Param(ws.QueryParameter("", "").DataType("")).
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]{{.Model.Name}}{}).
		Returns(200, "OK", []{{.Model.Name}}{}))

	ws.Route(ws.GET("/{id}").To(s.get).
		// docs
		Doc("get a {{.Service.Kind}}").
		Filter(s.measureTime).
		Param(ws.PathParameter("id", "identifier of the {{.Service.Kind}}").DataType("string")).
		// set more rich query condition
		Param(ws.QueryParameter("", "").DataType("")).
		// set more rich header 
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes({{.Model.Name}}{}). // on the response
		Returns(200, "OK", {{.Model.Name}}{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.PUT("/{id}").To(s.update).
		// docs
		Doc("update a {{.Service.Kind}}").
		Filter(s.measureTime).
		Param(ws.PathParameter("id", "identifier of the {{.Service.Kind}}").DataType("string")).
		// set more rich query condition
		Param(ws.QueryParameter("", "").DataType("")).
		// set more rich header 
		Param(ws.HeaderParameter("", "").DataType("")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads({{.Model.Name}}{})) // from the request

	ws.Route(ws.POST("").To(s.create).
		// docs
		Doc("create a {{.Service.Kind}}").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads({{.Model.Name}}{})) // from the request

	ws.Route(ws.DELETE("/{id}").To(s.delete).
		// docs
		Doc("delete a {{.Service.Kind}}").
		Filter(s.measureTime).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("id", "identifier of the {{.Service.Kind}}").DataType("string")))

	s.ws = ws
	s.container.Add(ws)
}
`

	t, err := template.New("ws-func-tplt").Parse(tmplt)
	if err != nil {
		return errors.WithMessage(err, "parse template")
	}

	if err := t.Execute(g.cache, g.config); err != nil {
		return errors.WithMessage(err, "execute template")
	}
	return nil
}

func (g *webServiceGenerator) generateRegisterFunc() error {

	return nil
}
