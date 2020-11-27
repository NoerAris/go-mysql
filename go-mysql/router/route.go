package router

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"log"
)

const ACTION_METHOD_TYPE_GET = "GET"
const ACTION_METHOD_TYPE_POST = "POST"
const ACTION_METHOD_TYPE_PUT = "PUT"
const ACTION_METHOD_TYPE_DELETE = "DELETE"
const ACTION_METHOD_TYPE_ANY = "ANY"

type Controller struct {
	Name        string    `json:"name"`
	BasePath    string    `json:"base_path"`
	Actions     []*Action `json:"actions"`
	PreRequest  func(ctx context.Context)
	PostRequest func(ctx context.Context)
}

type Action struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	AccessLevel string
	Process     func(ctx context.Context)
}

func (c *Controller) Get(path string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_GET, Process: process})
}
func (c *Controller) Post(path string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_POST, Process: process})
}
func (c *Controller) Put(path string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_PUT, Process: process})
}
func (c *Controller) Delete(path string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_DELETE, Process: process})
}
func (c *Controller) Any(path string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_ANY, Process: process})
}
func (c *Controller) GetWithRole(path string, accessLevel string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_GET, Process: process, AccessLevel: accessLevel})
}
func (c *Controller) PostWithRole(path string, accessLevel string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_POST, Process: process, AccessLevel: accessLevel})
}
func (c *Controller) PutWithRole(path string, accessLevel string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_PUT, Process: process, AccessLevel: accessLevel})
}
func (c *Controller) DeleteWithRole(path string, accessLevel string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_DELETE, Process: process, AccessLevel: accessLevel})
}
func (c *Controller) AnyWithRole(path string, accessLevel string, process func(ctx context.Context)) {
	c.Actions = append(c.Actions, &Action{Path: path, Method: ACTION_METHOD_TYPE_ANY, Process: process, AccessLevel: accessLevel})
}

var controllers []*Controller

func ProcessController(app *iris.Application) {
	if controllers == nil {
		log.Println("there aren't any controllers initialized yet..")
	} else {
		for _, c := range controllers {
			log.Println("Registering controller " + c.Name + " at " + c.BasePath)
			for _, a := range c.Actions {
				log.Println("----> Action " + a.Method + ":" + a.Path)
				switch a.Method {
				case ACTION_METHOD_TYPE_GET:
					app.Get(c.BasePath+a.Path, c.PreRequest, a.Process)
					break
				case ACTION_METHOD_TYPE_POST:
					app.Post(c.BasePath+a.Path, c.PreRequest, a.Process)
					break
				case ACTION_METHOD_TYPE_PUT:
					app.Put(c.BasePath+a.Path, c.PreRequest, a.Process)
					break
				case ACTION_METHOD_TYPE_DELETE:
					app.Delete(c.BasePath+a.Path, c.PreRequest, a.Process)
					break
				default:
					app.Any(c.BasePath+a.Path, c.PreRequest, a.Process)
					break
				}
			}
		}
	}
}
func ProcessControllerWithJwt(app *iris.Application) {
	if controllers == nil {
		log.Println("there aren't any controllers initialized yet..")
	} else {
		for _, c := range controllers {
			log.Println("Registering controller " + c.Name + " at " + c.BasePath)
			for _, a := range c.Actions {
				log.Println("----> Action " + a.Method + ":" + a.Path)
				var handlers []context.Handler
				handlers = append(handlers, c.PreRequest)
				switch a.AccessLevel {
				case RolePublic:
					handlers = append(handlers, CheckPublic, JwtValidator)
					break
				case RoleMemberRegistered:
					handlers = append(handlers, CheckRegistered, JwtValidator)
					break
				case RoleMemberApproved:
					handlers = append(handlers, CheckApproved, JwtValidator)
					break
				default:
					break
				}
				handlers = append(handlers, a.Process)
				switch a.Method {
				case ACTION_METHOD_TYPE_GET:
					app.Get(c.BasePath+a.Path, handlers...)
					break
				case ACTION_METHOD_TYPE_POST:
					app.Post(c.BasePath+a.Path, handlers...)
					break
				case ACTION_METHOD_TYPE_PUT:
					app.Put(c.BasePath+a.Path, handlers...)
					break
				case ACTION_METHOD_TYPE_DELETE:
					app.Delete(c.BasePath+a.Path, handlers...)
					break
				default:
					app.Any(c.BasePath+a.Path, handlers...)
					break
				}
			}
		}
	}
}

func CreateNewControllerInstance(name string, basePath string) *Controller {
	c := Controller{
		Name:        name,
		BasePath:    basePath,
		PreRequest:  func(ctx context.Context) { ctx.Next() },
		PostRequest: func(ctx context.Context) {},
		Actions:     make([]*Action, 0),
	}
	addController(&c)
	return &c
}

func addController(controller *Controller) {
	if controller == nil {
		controllers = make([]*Controller, 0)
	}
	controllers = append(controllers, controller)
}
