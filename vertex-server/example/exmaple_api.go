package example

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"gitlab.doit9.com/backend/vertex"
	"gitlab.doit9.com/backend/vertex/middleware"
)

type UserHandler struct {
	Id     string `schema:"id" required:"true" doc:"The Id Of the user" in:"path"`
	Name   string `schema:"name" maxlen:"100" required:"true" doc:"The Name Of the user"`
	Banana Banana `schema:"banana" required:"true"`
}

func (h UserHandler) Handle(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	fmt.Printf("%#v\n", h)
	return User{Id: h.Id, Name: h.Name, Banana: h.Banana}, nil
}

func testUserHandler(t *vertex.TestContext) error {

	vals := url.Values{}
	vals.Set("name", "foofi")
	params := vertex.Params{"id": "foo"}

	req, err := t.NewRequest("POST", vals, params)
	if err != nil {
		return err
	}

	resp := map[string]interface{}{}
	if r, err := t.JsonRequest(req, resp); err != nil {
		b, _ := ioutil.ReadAll(r.Body)
		t.Log("Got response: %v", string(b))
		return err
	}

	return nil
}

type User struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Banana Banana `json:"banana"`
}

type Banana struct {
	Foo string
	Bar string
}

func (b Banana) UnmarshalRequestData(data string) interface{} {

	parts := strings.Split(data, ",")
	if len(parts) == 2 {
		return Banana{parts[0], parts[1]}
	}
	return Banana{}
}

func init() {
	root := "/testung/1.0"
	vertex.RegisterAPI(&vertex.API{
		Host:          "localhost:9947",
		Name:          "testung",
		Version:       "1.0",
		Root:          root,
		Doc:           "Our fancy testung API",
		Title:         "Testung API!",
		Middleware:    middleware.DefaultMiddleware,
		Renderer:      vertex.JSONRenderer{},
		AllowInsecure: true,
		Routes: vertex.RouteMap{
			"/user/byId/{id}": {
				Description: "Get User Info by id or name",
				Handler:     UserHandler{},
				Methods:     vertex.GET,
				Test:        vertex.WarningTest(testUserHandler),
				Returns:     User{},
			},
			//			"/user/byName/{id}": {
			//				Description: "Get User Info by id or name",
			//				Handler:     UserHandler{},
			//				Methods:     vertex.POST,
			//				Test:        vertex.WarningTest(testUserHandler),
			//				Returns:     User{},
			//			},
			//			"/test/foo": {
			//				Description: "just for testing",
			//				Handler:     vertex.VoidHandler,
			//				Methods:     vertex.GET | vertex.POST,
			//			},
			//			"/test/bar": {
			//				Description: "Static",
			//				Handler:     vertex.StaticHandler(path.Join(root, "static"), http.Dir("/tmp")),
			//				Methods:     vertex.GET,
			//				Test:        vertex.DummyTest,
			//			},
		},
	})

}