package golph

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestProjects_ListProjects(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/project.query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"results": [{"name":"Test Project 1"}]}`)
	})

	projects, _, err := client.Projects.List(nil)
	if err != nil {
		t.Errorf("Projects.List returned error: %v", err)
	}

	expected := []Project{
		{Name: "Test Project 1"},
	}
	if !reflect.DeepEqual(projects, expected) {
		t.Errorf("Projects.List returned %+v, expected %+v", projects, expected)
	}
}

func TestProjects_ListProjectsMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/project.query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"results": [{"region":{"slug":"nyc3"},"droplet":{"id":1},"ip":"192.168.0.1"},{"region":{"slug":"nyc3"},"droplet":{"id":2},"ip":"192.168.0.2"}], "links":{"pages":{"next":"http://example.com/v2/floating_ips/?page=2"}}}`)
	})

	_, resp, err := client.Projects.List(nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 1)
}

func TestProjects_RetrievePageByNumber(t *testing.T) {
	setup()
	defer teardown()

	jBlob := `
	{
		"floating_ips": [{"region":{"slug":"nyc3"},"droplet":{"id":1},"ip":"192.168.0.1"},{"region":{"slug":"nyc3"},"droplet":{"id":2},"ip":"192.168.0.2"}],
		"links":{
			"pages":{
				"next":"http://example.com/v2/floating_ips/?page=3",
				"prev":"http://example.com/v2/floating_ips/?page=1",
				"last":"http://example.com/v2/floating_ips/?page=3",
				"first":"http://example.com/v2/floating_ips/?page=1"
			}
		}
	}`

	mux.HandleFunc("/api/project.query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{Page: 2}
	_, resp, err := client.Projects.List(opt)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestProjects_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/project.query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{"floating_ip":{"region":{"slug":"nyc3"},"droplet":{"id":1},"ip":"192.168.0.1"}}`)
	})

	project, _, err := client.Projects.Get("Test Project 1")
	if err != nil {
		t.Errorf("Projects.Get returned error: %v", err)
	}

	expected := &Project{Name: "Test Project 1"}
	if !reflect.DeepEqual(project, expected) {
		t.Errorf("Projects.Get returned %+v, expected %+v", project, expected)
	}
}

func TestProjects_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &ProjectCreateRequest{
		Name: "Test Project 1",
	}

	mux.HandleFunc("/api/project.create", func(w http.ResponseWriter, r *http.Request) {
		v := new(ProjectCreateRequest)
		err := json.NewDecoder(r.Body).Decode(v)
		if err != nil {
			t.Fatal(err)
		}

		testMethod(t, r, "POST")
		if !reflect.DeepEqual(v, createRequest) {
			t.Errorf("Request body = %+v, expected %+v", v, createRequest)
		}

		fmt.Fprint(w, `{"result":{"PHID": "PHID-1234567","name":"Test Project 1"}}`)
	})

	project, _, err := client.Projects.Create(createRequest)
	if err != nil {
		t.Errorf("Project.Create returned error: %v", err)
	}

	expected := &Project{PHID: "PHID-1234567", Name: "Test Project 1"}
	if !reflect.DeepEqual(project, expected) {
		t.Errorf("Projects.Create returned %+v, expected %+v", project, expected)
	}
}
