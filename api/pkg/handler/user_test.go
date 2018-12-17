package handler

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	apiv1 "github.com/kubermatic/kubermatic/api/pkg/api/v1"
	kubermaticapiv1 "github.com/kubermatic/kubermatic/api/pkg/crd/kubermatic/v1"

	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/diff"
	clienttesting "k8s.io/client-go/testing"
)

func TestGetUsersForProject(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		Name                        string
		ExpectedResponse            []apiv1.User
		ExpectedResponseString      string
		ExpectedActions             int
		ExpectedUserAfterInvitation *kubermaticapiv1.User
		ProjectToGet                string
		HTTPStatus                  int
		ExistingAPIUser             apiv1.User
		ExistingKubermaticObjs      []runtime.Object
	}{
		{
			Name:         "scenario 1: get a list of user for a project 'foo'",
			HTTPStatus:   http.StatusOK,
			ProjectToGet: "foo-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("foo", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("bar", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("zorg", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("foo-ID", "john@acme.com", "owners"),
				genBinding("bar-ID", "john@acme.com", "editors"),
				genBinding("foo-ID", "alice@acme.com", "viewers"),
				genBinding("foo-ID", "bob@acme.com", "editors"),
				genBinding("bar-ID", "bob@acme.com", "editors"),
				genBinding("zorg-ID", "bob@acme.com", "editors"),
				/*add users*/
				func() *kubermaticapiv1.User {
					user := genUser("", "john", "john@acme.com")
					user.CreationTimestamp = metav1.NewTime(time.Date(2013, 02, 03, 19, 54, 0, 0, time.UTC))
					return user
				}(),
				func() *kubermaticapiv1.User {
					user := genUser("", "alice", "alice@acme.com")
					user.CreationTimestamp = metav1.NewTime(time.Date(2013, 02, 03, 19, 55, 0, 0, time.UTC))
					return user
				}(),
				func() *kubermaticapiv1.User {
					user := genUser("", "bob", "bob@acme.com")
					user.CreationTimestamp = metav1.NewTime(time.Date(2013, 02, 03, 19, 56, 0, 0, time.UTC))
					return user
				}(),
			},
			ExistingAPIUser: *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: []apiv1.User{
				apiv1.User{
					ObjectMeta: apiv1.ObjectMeta{
						ID:                "4b2d8785b49bad23638b17d8db76857a79bf79441241a78a97d88cc64bbf766e",
						Name:              "john",
						CreationTimestamp: time.Date(2013, 02, 03, 19, 54, 0, 0, time.UTC),
					},
					Email: "john@acme.com",
					Projects: []apiv1.ProjectGroup{
						apiv1.ProjectGroup{
							GroupPrefix: "owners",
							ID:          "foo-ID",
						},
					},
				},

				apiv1.User{
					ObjectMeta: apiv1.ObjectMeta{
						ID:                "0a0a58273565a8f3dcf779375d9debd0f685d94dc56651a16bff3bf901c0b127",
						Name:              "alice",
						CreationTimestamp: time.Date(2013, 02, 03, 19, 55, 0, 0, time.UTC),
					},
					Email: "alice@acme.com",
					Projects: []apiv1.ProjectGroup{
						apiv1.ProjectGroup{
							GroupPrefix: "viewers",
							ID:          "foo-ID",
						},
					},
				},

				apiv1.User{
					ObjectMeta: apiv1.ObjectMeta{
						ID:                "405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6",
						Name:              "bob",
						CreationTimestamp: time.Date(2013, 02, 03, 19, 56, 0, 0, time.UTC),
					},
					Email: "bob@acme.com",
					Projects: []apiv1.ProjectGroup{
						apiv1.ProjectGroup{
							GroupPrefix: "editors",
							ID:          "foo-ID",
						},
					},
				},
			},
		},
		{
			Name:         "scenario 2: get a list of user for a project 'foo' for external user",
			HTTPStatus:   http.StatusForbidden,
			ProjectToGet: "foo2InternalName",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("foo2", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("bar2", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("bar-ID", "alice2@acme.com", "viewers"),
				genBinding("foo2-ID", "bob@acme.com", "editors"),
				genBinding("bar2-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "alice2", "alice2@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:        *genAPIUser("alice2", "alice2@acme.com"),
			ExpectedResponseString: `{"error":{"code":403,"message":"forbidden: The user \"alice2@acme.com\" doesn't belong to the given project = foo2InternalName"}}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/projects/%s/users", tc.ProjectToGet), nil)
			res := httptest.NewRecorder()
			kubermaticObj := []runtime.Object{}
			kubermaticObj = append(kubermaticObj, tc.ExistingKubermaticObjs...)
			ep, _, err := createTestEndpointAndGetClients(tc.ExistingAPIUser, nil, []runtime.Object{}, []runtime.Object{}, kubermaticObj, nil, nil)
			if err != nil {
				t.Fatalf("failed to create test endpoint due to %v", err)
			}

			ep.ServeHTTP(res, req)

			if res.Code != tc.HTTPStatus {
				t.Fatalf("Expected HTTP status code %d, got %d: %s", tc.HTTPStatus, res.Code, res.Body.String())
			}

			if len(tc.ExpectedResponse) > 0 {
				actualUsers := newUserV1SliceWrapper{}
				actualUsers.DecodeOrDie(res.Body, t).Sort()

				wrappedExpectedUsers := newUserV1SliceWrapper(tc.ExpectedResponse)
				wrappedExpectedUsers.Sort()

				actualUsers.EqualOrDie(wrappedExpectedUsers, t)
			}

			if len(tc.ExpectedResponseString) > 0 {
				compareWithResult(t, res, tc.ExpectedResponseString)
			}
		})
	}

}

func TestDeleteUserFromProject(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		Name                         string
		Body                         string
		ExpectedResponse             string
		ExpectedBindingIDAfterDelete string
		ProjectToSync                string
		UserIDToDelete               string
		HTTPStatus                   int
		ExistingAPIUser              apiv1.User
		ExistingKubermaticObjs       []runtime.Object
	}{
		// scenario 1
		{
			Name:          "scenario 1: john the owner of the plan9 project removes bob from the project",
			Body:          `{"id":"bobID", "email":"bob@acme.com", "projects":[{"id":"plan9", "group":"editors"}]}`,
			HTTPStatus:    http.StatusOK,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("plan9-ID", "bob@acme.com", "viewers"),
				genBinding("planX-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToDelete:               genDefaultUser().Name,
			ExistingAPIUser:              *genAPIUser("john", "john@acme.com"),
			ExpectedResponse:             `{}`,
			ExpectedBindingIDAfterDelete: genBinding("plan9-ID", "bob@acme.com", "viewers").Name,
		},

		// scenario 2
		{
			Name:          "scenario 2: john the owner of the plan9 project removes bob, but bob is not a member of the project",
			Body:          `{"id":"bobID", "email":"bob@acme.com", "projects":[{"id":"plan9", "group":"editors"}]}`,
			HTTPStatus:    http.StatusBadRequest,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("planX-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToDelete:   genDefaultUser().Name,
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":400,"message":"cannot delete the user = bob@acme.com from the project plan9-ID because the user is not a member of the project"}}`,
		},

		// scenario 3
		{
			Name:          "scenario 3: john the owner of the plan9 project removes himself from the projec",
			Body:          fmt.Sprintf(`{"id":"%s", "email":"%s", "projects":[{"id":"plan9", "group":"owners"}]}`, testUserID, testUserEmail),
			HTTPStatus:    http.StatusForbidden,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
			},
			UserIDToDelete:   genUser("", "john", "john@acme.com").Name,
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":403,"message":"you cannot delete yourself from the project"}}`,
		},

		// scenario 4
		{
			Name:          "scenario 4: email case insensitive. Remove bob from the project where email is Bob@acme.com instead bob@acme.com",
			Body:          `{"id":"bobID", "email":"Bob@acme.com", "projects":[{"id":"plan9", "group":"editors"}]}`,
			HTTPStatus:    http.StatusOK,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("plan9-ID", "bob@acme.com", "viewers"),
				genBinding("planX-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToDelete:               genDefaultUser().Name,
			ExistingAPIUser:              *genAPIUser("john", "john@acme.com"),
			ExpectedResponse:             `{}`,
			ExpectedBindingIDAfterDelete: genBinding("plan9-ID", "bob@acme.com", "viewers").Name,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/projects/%s/users/%s", tc.ProjectToSync, tc.UserIDToDelete), strings.NewReader(tc.Body))
			res := httptest.NewRecorder()
			kubermaticObj := []runtime.Object{}
			kubermaticObj = append(kubermaticObj, tc.ExistingKubermaticObjs...)
			ep, clients, err := createTestEndpointAndGetClients(tc.ExistingAPIUser, nil, []runtime.Object{}, []runtime.Object{}, kubermaticObj, nil, nil)
			if err != nil {
				t.Fatalf("failed to create test endpoint due to %v", err)
			}

			ep.ServeHTTP(res, req)

			if res.Code != tc.HTTPStatus {
				t.Fatalf("Expected HTTP status code %d, got %d: %s", tc.HTTPStatus, res.Code, res.Body.String())
			}
			compareWithResult(t, res, tc.ExpectedResponse)

			kubermaticFakeClient := clients.fakeKubermaticClient
			{
				if len(tc.ExpectedBindingIDAfterDelete) > 0 {
					actionWasValidated := false
					for _, action := range kubermaticFakeClient.Actions() {
						if action.Matches("delete", "userprojectbindings") {
							deleteAction, ok := action.(clienttesting.DeleteAction)
							if !ok {
								t.Fatalf("unexpected action %#v", action)
							}
							if deleteAction.GetName() != tc.ExpectedBindingIDAfterDelete {
								t.Fatalf("wrong binding removed, wanted = %s, actual = %s", tc.ExpectedBindingIDAfterDelete, deleteAction.GetName())
							}
							actionWasValidated = true
							break
						}
					}
					if !actionWasValidated {
						t.Fatal("create action was not validated, a binding for a user was not updated ?")
					}
				}
			}
		})
	}
}

func TestEditUserInProject(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		Name                       string
		Body                       string
		ExpectedResponse           string
		ExpectedBindingAfterUpdate *kubermaticapiv1.UserProjectBinding
		ProjectToSync              string
		UserIDToUpdate             string
		HTTPStatus                 int
		ExistingAPIUser            apiv1.User
		ExistingKubermaticObjs     []runtime.Object
	}{
		// scenario 1
		{
			Name:          "scenario 1: john the owner of the plan9 project changes the group for bob from viewers to editors",
			Body:          `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6", "email":"bob@acme.com", "projects":[{"id":"plan9-ID", "group":"editors"}]}`,
			HTTPStatus:    http.StatusOK,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("plan9-ID", "bob@acme.com", "viewers"),
				genBinding("my-third-project-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToUpdate:   genDefaultUser().Name,
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6","name":"Bob","creationTimestamp":"0001-01-01T00:00:00Z","email":"bob@acme.com","projects":[{"id":"plan9-ID","group":"editors"}]}`,
			ExpectedBindingAfterUpdate: func() *kubermaticapiv1.UserProjectBinding {
				binding := genBinding("plan9-ID", "bob@acme.com", "editors")
				// the name of the original binding was derived from projectID, email and group
				binding.Name = genBinding("plan9-ID", "bob@acme.com", "viewers").Name
				return binding
			}(),
		},

		// scenario 2
		{
			Name:          "scenario 2: john the owner of the plan9 project changes the group for bob, but bob is not a member of the project",
			Body:          `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6", "email":"bob@acme.com", "projects":[{"id":"plan9-ID", "group":"editors"}]}`,
			HTTPStatus:    http.StatusBadRequest,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToUpdate:   genDefaultUser().Name,
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":400,"message":"cannot change the membership of the user = bob@acme.com for the project plan9-ID because the user is not a member of the project"}}`,
		},

		// scenario 3
		{
			Name:          "scenario 3: john the owner of the plan9 project changes the group for bob from viewers to owners",
			Body:          `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6", "email":"bob@acme.com", "projects":[{"id":"plan9-ID", "group":"owners"}]}`,
			HTTPStatus:    http.StatusForbidden,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("plan9-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToUpdate:   genDefaultUser().Name,
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":403,"message":"the given user cannot be assigned to owners group"}}`,
		},

		// scenario 4
		{
			Name:          "scenario 4: john the owner of the plan9 project changes the group for bob from viewers to admins(wrong name)",
			Body:          `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6", "email":"bob@acme.com", "projects":[{"id":"plan9-ID", "group":"admins"}]}`,
			HTTPStatus:    http.StatusBadRequest,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("plan9-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToUpdate:   genDefaultUser().Name,
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":400,"message":"invalid group name admins"}}`,
		},

		// scenario 5
		{
			Name:          "scenario 5: email case insensitive. Changes the group for bob from viewers to editors where email is BOB@ACME.COM instead bob@acme.com",
			Body:          `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6", "email":"BOB@ACME.COM", "projects":[{"id":"plan9-ID", "group":"editors"}]}`,
			HTTPStatus:    http.StatusOK,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("plan9-ID", "bob@acme.com", "viewers"),
				genBinding("my-third-project-ID", "bob@acme.com", "viewers"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			UserIDToUpdate:   genDefaultUser().Name,
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6","name":"Bob","creationTimestamp":"0001-01-01T00:00:00Z","email":"bob@acme.com","projects":[{"id":"plan9-ID","group":"editors"}]}`,
			ExpectedBindingAfterUpdate: func() *kubermaticapiv1.UserProjectBinding {
				binding := genBinding("plan9-ID", "bob@acme.com", "editors")
				// the name of the original binding was derived from projectID, email and group
				binding.Name = genBinding("plan9-ID", "bob@acme.com", "viewers").Name
				return binding
			}(),
		},
	}
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/projects/%s/users/%s", tc.ProjectToSync, tc.UserIDToUpdate), strings.NewReader(tc.Body))
			res := httptest.NewRecorder()
			kubermaticObj := []runtime.Object{}
			kubermaticObj = append(kubermaticObj, tc.ExistingKubermaticObjs...)
			ep, clients, err := createTestEndpointAndGetClients(tc.ExistingAPIUser, nil, []runtime.Object{}, []runtime.Object{}, kubermaticObj, nil, nil)
			if err != nil {
				t.Fatalf("failed to create test endpoint due to %v", err)
			}

			ep.ServeHTTP(res, req)

			if res.Code != tc.HTTPStatus {
				t.Fatalf("Expected HTTP status code %d, got %d: %s", tc.HTTPStatus, res.Code, res.Body.String())
			}
			compareWithResult(t, res, tc.ExpectedResponse)

			kubermaticFakeClient := clients.fakeKubermaticClient
			{
				if tc.ExpectedBindingAfterUpdate != nil {
					actionWasValidated := false
					for _, action := range kubermaticFakeClient.Actions() {
						if action.Matches("update", "userprojectbindings") {
							updateAction, ok := action.(clienttesting.UpdateAction)
							if !ok {
								t.Fatalf("unexpected action %#v", action)
							}
							updatedBinding := updateAction.GetObject().(*kubermaticapiv1.UserProjectBinding)
							if !equality.Semantic.DeepEqual(updatedBinding, tc.ExpectedBindingAfterUpdate) {
								t.Fatalf("updated action mismatch %v", diff.ObjectDiff(updatedBinding, tc.ExpectedBindingAfterUpdate))
							}
							actionWasValidated = true
							break
						}
					}
					if !actionWasValidated {
						t.Fatal("create action was not validated, a binding for a user was not updated ?")
					}
				}
			}
		})
	}
}

func TestAddUserToProject(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		Name                           string
		Body                           string
		ExpectedResponse               string
		ExpectedBindingAfterInvitation *kubermaticapiv1.UserProjectBinding
		ProjectToSync                  string
		HTTPStatus                     int
		ExistingAPIUser                apiv1.User
		ExistingKubermaticObjs         []runtime.Object
	}{
		{
			Name:          "scenario 1: john the owner of the plan9 project invites bob to the project as an editor",
			Body:          `{"email":"bob@acme.com", "projects":[{"id":"plan9-ID", "group":"editors"}]}`,
			HTTPStatus:    http.StatusCreated,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "john@acme.com", "editors"),
				genBinding("placeX-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6","name":"Bob","creationTimestamp":"0001-01-01T00:00:00Z","email":"bob@acme.com","projects":[{"id":"plan9-ID","group":"editors"}]}`,
			ExpectedBindingAfterInvitation: &kubermaticapiv1.UserProjectBinding{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: kubermaticapiv1.SchemeGroupVersion.String(),
							Kind:       kubermaticapiv1.ProjectKindName,
							Name:       "plan9-ID",
						},
					},
				},
				Spec: kubermaticapiv1.UserProjectBindingSpec{
					UserEmail: "bob@acme.com",
					Group:     "editors-plan9-ID",
					ProjectID: "plan9-ID",
				},
			},
		},

		{
			Name:          "scenario 2: john the owner of the plan9 project tries to invite bob to another project",
			Body:          `{"email":"bob@acme.com", "projects":[{"id":"moby", "group":"editors"}]}`,
			HTTPStatus:    http.StatusForbidden,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/* add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "john@acme.com", "editors"),
				genBinding("placeX-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":403,"message":"you can only assign the user to plan9-ID project"}}`,
		},

		{
			Name:          "scenario 3: john the owner of the plan9 project tries to invite  himself to another group",
			Body:          fmt.Sprintf(`{"email":"%s", "projects":[{"id":"plan9-ID", "group":"editors"}]}`, testUserEmail),
			HTTPStatus:    http.StatusForbidden,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "john@acme.com", "editors"),
				genBinding("placeX-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":403,"message":"you cannot assign yourself to a different group"}}`,
		},

		{
			Name:          "scenario 4: john the owner of the plan9 project tries to invite bob to the project as an owner",
			Body:          `{"email":"bob@acme.com", "projects":[{"id":"plan9-ID", "group":"owners"}]}`,
			HTTPStatus:    http.StatusForbidden,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "john@acme.com", "editors"),
				genBinding("placeX-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":403,"message":"the given user cannot be assigned to owners group"}}`,
		},

		{
			Name:          "scenario 5: john the owner of the plan9 project invites bob to the project one more time",
			Body:          `{"email":"bob@acme.com", "projects":[{"id":"plan9-ID", "group":"editors"}]}`,
			HTTPStatus:    http.StatusBadRequest,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "john@acme.com", "editors"),
				genBinding("plan9-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":400,"message":"cannot add the user = bob@acme.com to the project plan9-ID because user is already in the project"}}`,
		},

		{
			Name:          "scenario 6: email case insensitive. Bob is invited to the project one more time and emil starts with capital letter",
			Body:          `{"email":"Bob@acme.com", "projects":[{"id":"plan9-ID", "group":"editors"}]}`,
			HTTPStatus:    http.StatusBadRequest,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "john@acme.com", "editors"),
				genBinding("plan9-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"error":{"code":400,"message":"cannot add the user = Bob@acme.com to the project plan9-ID because user is already in the project"}}`,
		},

		{
			Name:          "scenario 7: email case insensitive. Bob is invited to the project as an editor with capital letter email",
			Body:          `{"email":"BOB@ACME.COM", "projects":[{"id":"plan9-ID", "group":"editors"}]}`,
			HTTPStatus:    http.StatusCreated,
			ProjectToSync: "plan9-ID",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("my-first-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("my-third-project", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				genBinding("my-third-project-ID", "john@acme.com", "editors"),
				genBinding("placeX-ID", "bob@acme.com", "editors"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
				genDefaultUser(), /*bob*/
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedResponse: `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6","name":"Bob","creationTimestamp":"0001-01-01T00:00:00Z","email":"bob@acme.com","projects":[{"id":"plan9-ID","group":"editors"}]}`,
			ExpectedBindingAfterInvitation: &kubermaticapiv1.UserProjectBinding{
				ObjectMeta: metav1.ObjectMeta{
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: kubermaticapiv1.SchemeGroupVersion.String(),
							Kind:       kubermaticapiv1.ProjectKindName,
							Name:       "plan9-ID",
						},
					},
				},
				Spec: kubermaticapiv1.UserProjectBindingSpec{
					UserEmail: "bob@acme.com",
					Group:     "editors-plan9-ID",
					ProjectID: "plan9-ID",
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			req := httptest.NewRequest("POST", fmt.Sprintf("/api/v1/projects/%s/users", tc.ProjectToSync), strings.NewReader(tc.Body))
			res := httptest.NewRecorder()
			kubermaticObj := []runtime.Object{}
			kubermaticObj = append(kubermaticObj, tc.ExistingKubermaticObjs...)
			ep, clients, err := createTestEndpointAndGetClients(tc.ExistingAPIUser, nil, []runtime.Object{}, []runtime.Object{}, kubermaticObj, nil, nil)
			if err != nil {
				t.Fatalf("failed to create test endpoint due to %v", err)
			}

			ep.ServeHTTP(res, req)

			if res.Code != tc.HTTPStatus {
				t.Fatalf("Expected HTTP status code %d, got %d: %s", tc.HTTPStatus, res.Code, res.Body.String())
			}
			compareWithResult(t, res, tc.ExpectedResponse)

			kubermaticFakeClient := clients.fakeKubermaticClient
			{
				if tc.ExpectedBindingAfterInvitation != nil {
					actionWasValidated := false
					for _, action := range kubermaticFakeClient.Actions() {
						if action.Matches("create", "userprojectbindings") {
							updateAction, ok := action.(clienttesting.CreateAction)
							if !ok {
								t.Fatalf("unexpected action %#v", action)
							}
							createdBinding := updateAction.GetObject().(*kubermaticapiv1.UserProjectBinding)
							// Name was generated by the test framework, just rewrite it
							tc.ExpectedBindingAfterInvitation.Name = createdBinding.Name
							if !equality.Semantic.DeepEqual(createdBinding, tc.ExpectedBindingAfterInvitation) {
								t.Fatalf("%v", diff.ObjectDiff(createdBinding, tc.ExpectedBindingAfterInvitation))
							}
							actionWasValidated = true
							break
						}
					}
					if !actionWasValidated {
						t.Fatal("create action was not validated, a binding for a user was not created ?")
					}
				}
			}
		})
	}
}

func TestGetCurrentUser(t *testing.T) {
	testcases := []struct {
		Name                   string
		ExpectedResponse       string
		ExpectedStatus         int
		ExistingKubermaticObjs []runtime.Object
		ExistingAPIUser        apiv1.User
	}{
		{
			Name: "scenario 1: get john's profile (no projects assigned)",
			ExistingKubermaticObjs: []runtime.Object{
				/*add users*/
				genUser("", "john", "john@acme.com"),
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: `{"id":"4b2d8785b49bad23638b17d8db76857a79bf79441241a78a97d88cc64bbf766e","name":"john","creationTimestamp":"0001-01-01T00:00:00Z","email":"john@acme.com"}`,
		},

		{
			Name: "scenario 2: get john's profile (one project assigned)",
			ExistingKubermaticObjs: []runtime.Object{
				/*add projects*/
				genProject("moby", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				genProject("plan9", kubermaticapiv1.ProjectActive, defaultCreationTimestamp()),
				/*add bindings*/
				genBinding("plan9-ID", "john@acme.com", "owners"),
				/*add users*/
				genUser("", "john", "john@acme.com"),
			},
			ExistingAPIUser:  *genAPIUser("john", "john@acme.com"),
			ExpectedStatus:   http.StatusOK,
			ExpectedResponse: `{"id":"4b2d8785b49bad23638b17d8db76857a79bf79441241a78a97d88cc64bbf766e","name":"john","creationTimestamp":"0001-01-01T00:00:00Z","email":"john@acme.com","projects":[{"id":"plan9-ID","group":"owners"}]}`,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			kubermaticObj := []runtime.Object{}
			kubermaticObj = append(kubermaticObj, tc.ExistingKubermaticObjs...)
			ep, _, err := createTestEndpointAndGetClients(tc.ExistingAPIUser, nil, []runtime.Object{}, []runtime.Object{}, kubermaticObj, nil, nil)
			if err != nil {
				t.Fatalf("failed to create test endpoint due to %v", err)
			}

			req := httptest.NewRequest("GET", "/api/v1/me", nil)
			res := httptest.NewRecorder()
			ep.ServeHTTP(res, req)

			if res.Code != tc.ExpectedStatus {
				t.Fatalf("Expected HTTP status code %d, got %d: %s", tc.ExpectedStatus, res.Code, res.Body.String())
			}
			compareWithResult(t, res, tc.ExpectedResponse)
		})
	}
}

func TestNewUser(t *testing.T) {
	t.Parallel()
	testcases := []struct {
		Name                      string
		HTTPStatus                int
		ExpectedResponse          string
		ExpectedKubermaticUser    *kubermaticapiv1.User
		ExistingKubermaticObjects []runtime.Object
		ExistingAPIUser           *apiv1.User
	}{
		{
			Name:             "scenario 1: successfully creates a new user resource",
			ExpectedResponse: `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6","name":"Bob","creationTimestamp":"0001-01-01T00:00:00Z","email":"bob@acme.com"}`,
			HTTPStatus:       http.StatusOK,
			ExpectedKubermaticUser: func() *kubermaticapiv1.User {
				apiUser := genDefaultAPIUser()
				expectedKubermaticUser := apiUserToKubermaticUser(*apiUser)
				// the name of the object is derived from the email address and encoded as sha256
				expectedKubermaticUser.Name = fmt.Sprintf("%x", sha256.Sum256([]byte(apiUser.Email)))
				return expectedKubermaticUser
			}(),
			ExistingAPIUser: genDefaultAPIUser(),
		},

		{
			Name:             "scenario 2: fails when creating a user without an email address",
			ExpectedResponse: `{"error":{"code":400,"message":"Email, ID and Name cannot be empty when creating a new user resource"}}`,
			HTTPStatus:       http.StatusBadRequest,
			ExistingAPIUser: func() *apiv1.User {
				apiUser := genDefaultAPIUser()
				apiUser.Email = ""
				return apiUser
			}(),
		},

		{
			Name:             "scenario 3: creating a user if already exists doesn't have effect",
			ExpectedResponse: `{"id":"405ac8384fa984f787f9486daf34d84d98f20c4d6a12e2cc4ed89be3bcb06ad6","name":"Bob","creationTimestamp":"0001-01-01T00:00:00Z","email":"bob@acme.com"}`,
			HTTPStatus:       http.StatusOK,
			ExistingKubermaticObjects: []runtime.Object{
				genDefaultUser(),
			},
			ExistingAPIUser: genDefaultAPIUser(),
		},
	}

	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {

			ep, clientSet, err := createTestEndpointAndGetClients(*tc.ExistingAPIUser, nil, []runtime.Object{}, []runtime.Object{}, tc.ExistingKubermaticObjects, nil, nil)
			if err != nil {
				t.Fatalf("failed to create test endpoint due to %v", err)
			}

			// act
			req := httptest.NewRequest("GET", "/api/v1/me", nil)
			res := httptest.NewRecorder()
			ep.ServeHTTP(res, req)

			// validate
			if res.Code != tc.HTTPStatus {
				t.Fatalf("Expected HTTP status code %d, got %d: %s", http.StatusOK, res.Code, res.Body.String())
			}
			compareWithResult(t, res, tc.ExpectedResponse)

			for _, action := range clientSet.fakeKubermaticClient.Actions() {
				if action.Matches("create", "users") {
					createAction, ok := action.(clienttesting.CreateAction)
					if !ok {
						t.Fatalf("unexpected action %#v", action)
					}
					if !equality.Semantic.DeepEqual(createAction.GetObject().(*kubermaticapiv1.User), tc.ExpectedKubermaticUser) {
						t.Fatalf("%v", diff.ObjectDiff(tc.ExpectedKubermaticUser, createAction.GetObject().(*kubermaticapiv1.User)))
					}
					return /*pass*/
				}
			}
			if tc.ExpectedKubermaticUser != nil {
				t.Fatal("expected to find create action (fake client) but haven't received one.")
			}
		})
	}
}

// genUser generates a User resource
// note if the id is empty then it will be auto generated
func genUser(id, name, email string) *kubermaticapiv1.User {
	if len(id) == 0 {
		// the name of the object is derived from the email address and encoded as sha256
		id = fmt.Sprintf("%x", sha256.Sum256([]byte(email)))
	}
	return &kubermaticapiv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		Spec: kubermaticapiv1.UserSpec{
			Name:  name,
			Email: email,
		},
	}
}

func genDefaultUser() *kubermaticapiv1.User {
	userEmail := "bob@acme.com"
	return genUser("", "Bob", userEmail)
}

func genAPIUser(name, email string) *apiv1.User {
	usr := genUser("", name, email)
	return &apiv1.User{
		ObjectMeta: apiv1.ObjectMeta{
			ID:   usr.Name,
			Name: usr.Spec.Name,
		},
		Email: usr.Spec.Email,
	}
}

func genDefaultAPIUser() *apiv1.User {
	return &apiv1.User{
		ObjectMeta: apiv1.ObjectMeta{
			ID:   genDefaultUser().Name,
			Name: genDefaultUser().Spec.Name,
		},
		Email: genDefaultUser().Spec.Email,
	}
}

func genBinding(projectID, email, group string) *kubermaticapiv1.UserProjectBinding {
	return &kubermaticapiv1.UserProjectBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s-%s", projectID, email, group),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: kubermaticapiv1.SchemeGroupVersion.String(),
					Kind:       kubermaticapiv1.ProjectKindName,
					Name:       projectID,
				},
			},
		},
		Spec: kubermaticapiv1.UserProjectBindingSpec{
			UserEmail: email,
			ProjectID: projectID,
			Group:     fmt.Sprintf("%s-%s", group, projectID),
		},
	}
}

func genDefaultOwnerBinding() *kubermaticapiv1.UserProjectBinding {
	return genBinding(genDefaultProject().Name, genDefaultUser().Spec.Email, "owners")
}
