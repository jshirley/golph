package golph

import (
	"errors"
	"fmt"
	"net/url"
)

const projectsQueryPath = "api/project.query"
const projectsCreatePath = "api/project.create"

// ProjectsService is an interface for interfacing with Projects
// See: https://secure.phabricator.com/conduit/ (and search for projects)
type ProjectsService interface {
	List(*ListOptions) ([]Project, *Response, error)
	Get(string) (*Project, *Response, error)
	Create(*ProjectCreateRequest) (*Project, *Response, error)
	Update(*ProjectUpdateRequest) (*Response, error)
	Delete(*Project) (*Response, error)
}

// ProjectsServiceOp handles communication with the conduit methods
type ProjectsServiceOp struct {
	client *Client
}

var _ ProjectsService = &ProjectsServiceOp{}

// Project represents a Phabricator project.
type Project struct {
	PHID    string   `json:"phid"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags"`
	Members []string `json:"members"`
	Icon    string   `json:"icon"`
	Color   string   `json:"color"`
}

func (f Project) String() string {
	return Stringify(f)
}

// What is returned by a base project.query request
type projectsRoot struct {
	Projects []Project `json:"results"`
	Links    *Links    `json:"links"`
}

// Single request to a project root (Phabricator doesn't understand this)
type projectRoot struct {
	Project *Project `json:"results"`
	Links   *Links   `json:"links,omitempty"`
}

// ProjectCreateRequest represents a request to create a Project.
type ProjectCreateRequest struct {
	Name    string   `json:"name"`
	Tags    []string `json:"tags,omitempty"`
	Members []string `json:"members,omitempty"`
	Icon    string   `json:"icon,omitempty"`
	Color   string   `json:"color,omitempty"`
}

// ProjectUpdateRequest represents a request to update a Project.
// Note: Phabricator doesn't support this, I'm just optimistic.
type ProjectUpdateRequest struct {
	PHID    string   `json:"phid"`
	Name    string   `json:"name"`
	Tags    []string `json:"tags,omitempty"`
	Members []string `json:"members,omitempty"`
	Icon    string   `json:"icon,omitempty"`
	Color   string   `json:"color,omitempty"`
}

// List all projects.
func (f *ProjectsServiceOp) List(opt *ListOptions) ([]Project, *Response, error) {
	path := projectsQueryPath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := f.client.NewRequest("POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectsRoot)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Projects, resp, err
}

// Get an individual project.
func (f *ProjectsServiceOp) Get(name string) (*Project, *Response, error) {
	form := url.Values{}
	form.Add("names", fmt.Sprintf("[\"%s\"]", name))

	req, err := f.client.NewRequest("POST", projectsQueryPath, form)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectsRoot)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	if len(root.Projects) < 1 {
		return nil, resp, nil
	}

	return &root.Projects[0], resp, err
}

// Create a project
func (f *ProjectsServiceOp) Create(createRequest *ProjectCreateRequest) (*Project, *Response, error) {
	path := projectsCreatePath

	req, err := f.client.NewRequest("POST", path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(projectRoot)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}
	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Project, resp, err
}

// Update a Project (Not Available Yet!)
func (f *ProjectsServiceOp) Update(project *ProjectUpdateRequest) (*Response, error) {
	return nil, errors.New("Update is not available for Projects")
}

// Delete a Project (Not Available Yet!)
func (f *ProjectsServiceOp) Delete(project *Project) (*Response, error) {
	return nil, errors.New("Delete is not available for Projects")
}
