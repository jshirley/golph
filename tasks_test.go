package golph

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	listTaskJSON   = `{"result":{"PHID-TASK-1":{"id":"1000","phid":"PHID-TASK-1","authorPHID":"PHID-USER-1","ownerPHID":null,"ccPHIDs":["PHID-USER-2","PHID-USER-3"],"status":"resolved","statusName":"Resolved","isClosed":true,"priority":"Needs Triage","priorityColor":"violet","title":"Test List","description":"We only test one item in the list, because Phab returns a map and the order isn't fixed :(","projectPHIDs":["PHID-PROJ-1","PHID-PROJ-2"],"uri":"https://phabricator.example.com/T1000","auxiliary":{},"objectName":"T1000","dateCreated":"1415646583","dateModified":"1451336014","dependsOnTaskPHIDs":[]}},"error_code":null,"error_info":null}`
	updateTaskJSON = `{"result":{"id":"1000","phid":"PHID-TASK-1","authorPHID":"PHID-USER-1","ownerPHID":null,"ccPHIDs":["PHID-USER-2","PHID-USER-3"],"status":"resolved","statusName":"Resolved","isClosed":true,"priority":"Needs Triage","priorityColor":"violet","title":"Awesome Task","description":"A description goes here","projectPHIDs":["PHID-PROJ-1","PHID-PROJ-2"],"uri":"https://phabricator.example.com/T1000","auxiliary":{},"objectName":"T1000","dateCreated":"1415646583","dateModified":"1451336014","dependsOnTaskPHIDs":[]},"error_code":null,"error_info":null}`
)

func TestTasks_ListTasks(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/maniphest.query", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, listTaskJSON)
	})

	tasks, _, err := client.Tasks.List(nil)
	if err != nil {
		t.Errorf("Tasks.List returned error: %v", err)
	}

	expected := []Task{
		{
			PHID:          "PHID-PROJ-1",
			Author:        "PHID-USER-1",
			Owner:         "",
			Title:         "Test List",
			Description:   "We only test one item in the list, because Phab returns a map and the order isn't fixed :(",
			Projects:      []string{"PHID-PROJ-1"},
			CCs:           []string{"PHID-USER-2", "PHID-USER-3"},
			Status:        "resolved",
			StatusName:    "Resolved",
			IsClosed:      true,
			Priority:      "Needs Triage",
			PriorityColor: "violet",
		},
	}

	if !reflect.DeepEqual(tasks, expected) {
		t.Errorf("Tasks.List returned %+v, expected %+v", tasks, expected)
	}
}

func TestTasks_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/maniphest.info", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		fmt.Fprint(w, listTaskJSON)
	})

	task, _, err := client.Tasks.Get("2000")
	if err != nil {
		t.Errorf("Tasks.Get returned error: %v", err)
	}

	expected := Task{
		PHID:          "PHID-PROJ-1",
		Author:        "PHID-USER-1",
		Owner:         "",
		Title:         "Test List",
		Description:   "We only test one item in the list, because Phab returns a map and the order isn't fixed :(",
		Projects:      []string{"PHID-PROJ-1"},
		CCs:           []string{"PHID-USER-1", "PHID-USER-2"},
		Status:        "resolved",
		StatusName:    "Resolved",
		IsClosed:      true,
		Priority:      "Needs Triage",
		PriorityColor: "violet",
		URI:           "https://phabricator.example.com/T2000",
		ObjectName:    "T2000",
	}

	if !reflect.DeepEqual(task, expected) {
		t.Errorf("Tasks.Get returned %+v, expected %+v", task, expected)
	}
}

func TestTasks_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &TaskCreateRequest{
		Title:       "Test List",
		Description: "We only test one item in the list, because Phab returns a map and the order isn't fixed :(",
		Projects:    `["PHID-PROJ-1"]`,
		OwnerPHID:   "PHID-USER-1",
		CCs:         `["PHID-USER-2"]`,
		Priority:    "Needs Triage",
	}

	mux.HandleFunc("/api/maniphest.createtask", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		if r.PostFormValue("title") != createRequest.Title {
			t.Errorf("Form name = %+v, expected %+v", r.PostFormValue("title"), createRequest.Title)
		}
		fmt.Fprint(w, `{"result":{"id":"2000","phid":"PHID-TASK-2","authorPHID":"PHID-USER-1","ownerPHID":null,"ccPHIDs":["PHID-USER-2"],"status":"open","statusName":"Open","isClosed":false,"priority":"Normal","priorityColor":"green","title":"Testing Golph","description":"Golph created this task","projectPHIDs":["PHID-PROJ-1"],"uri":"https://phabricator.example.com/T2000","auxiliary":{},"objectName":"T2000","dateCreated":"1451337180","dateModified":"1451337180","dependsOnTaskPHIDs":[]},"error_code":null,"error_info":null}`)
	})

	task, _, err := client.Tasks.Create(createRequest)
	if err != nil {
		t.Errorf("Task.Create returned error: %v", err)
	}

	expected := &Task{
		PHID:          "PHID-TASK-2",
		Author:        "PHID-USER-1",
		Owner:         "",
		Title:         "Testing Golph",
		Description:   "Golph created this task",
		Projects:      []string{"PHID-PROJ-1"},
		CCs:           []string{"PHID-USER-2"},
		Status:        "open",
		StatusName:    "Open",
		IsClosed:      false,
		Priority:      "Normal",
		PriorityColor: "green",
		URI:           "https://phabricator.example.com/T2000",
		ObjectName:    "T2000",
	}

	if !reflect.DeepEqual(task, expected) {
		t.Errorf("Task.Create returned:\n%+v\nExpected:\n%+v", task, expected)
	}
}

func TestTasks_Update(t *testing.T) {
	setup()
	defer teardown()

	t.Errorf("Task.Update")
}
