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
		fmt.Fprint(w, `{"result":{"data":{"PHID-PROJ-1":{"id":"181","phid":"PHID-PROJ-1","name":"Project 1","profileImagePHID":"PHID-FILE-1","icon":"flag-checkered","color":"disabled","members":["PHID-USER-1","PHID-USER-2"],"slugs":["project_1"],"dateCreated":"1445305386","dateModified":"1446586132"},"PHID-PROJ-2":{"id":"2","phid":"PHID-PROJ-2","name":"Project 2","profileImagePHID":"PHID-FILE-2","icon":"umbrella","color":"disabled","members":["PHID-USER-1"],"slugs":["project_2"],"dateCreated":"1447804194","dateModified":"1448327625"}},"slugMap":[],"cursor":{"limit":2,"after":"35","before":null}},"error_code":null,"error_info":null}`)
	})

	projects, _, err := client.Projects.List(nil)
	if err != nil {
		t.Errorf("Projects.List returned error: %v", err)
	}

	expected := []Project{
		{PHID: "PHID-PROJ-1", Name: "Project 1", Icon: "flag-checkered", Color: "disabled", Members: []string{"PHID-USER-1", "PHID-USER-2"}, Tags: []string{"project_1"}},
		{PHID: "PHID-PROJ-2", Name: "Project 2", Icon: "umbrella", Color: "disabled", Members: []string{"PHID-USER-1"}, Tags: []string{"project_2"}},
	}
	if !reflect.DeepEqual(projects, expected) {
		t.Errorf("Projects.List returned %+v, expected %+v", projects, expected)
	}
}

/*
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
*/

func TestProjects_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/project.query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		fmt.Fprint(w, `{"result":{"data":{"PHID-PROJ-1":{"id":"181","phid":"PHID-PROJ-1","name":"Project 1","profileImagePHID":"PHID-FILE-1","icon":"flag-checkered","color":"disabled","members":["PHID-USER-1","PHID-USER-2"],"slugs":["project_1"],"dateCreated":"1445305386","dateModified":"1446586132"}},"slugMap":[],"cursor":{"limit":1,"after":"1","before":null}},"error_code":null,"error_info":null}`)
	})

	project, _, err := client.Projects.Get("Project 1")
	if err != nil {
		t.Errorf("Projects.Get returned error: %v", err)
	}

	expected := &Project{PHID: "PHID-PROJ-1", Name: "Project 1", Icon: "flag-checkered", Color: "disabled", Members: []string{"PHID-USER-1", "PHID-USER-2"}, Tags: []string{"project_1"}}

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
	/*
		project, _, err := client.Projects.Create(createRequest)
		if err != nil {
			t.Errorf("Project.Create returned error: %v", err)
		}

		expected := &Project{PHID: "PHID-1234567", Name: "Test Project 1"}
		if !reflect.DeepEqual(project, expected) {
			t.Errorf("Projects.Create returned %+v, expected %+v", project, expected)
		}
	*/
}
