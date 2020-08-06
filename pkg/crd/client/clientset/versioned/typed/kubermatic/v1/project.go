// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"time"

	scheme "k8c.io/kubermatic/v2/pkg/crd/client/clientset/versioned/scheme"
	v1 "k8c.io/kubermatic/v2/pkg/crd/kubermatic/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ProjectsGetter has a method to return a ProjectInterface.
// A group's client should implement this interface.
type ProjectsGetter interface {
	Projects() ProjectInterface
}

// ProjectInterface has methods to work with Project resources.
type ProjectInterface interface {
	Create(*v1.Project) (*v1.Project, error)
	Update(*v1.Project) (*v1.Project, error)
	UpdateStatus(*v1.Project) (*v1.Project, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Project, error)
	List(opts metav1.ListOptions) (*v1.ProjectList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Project, err error)
	ProjectExpansion
}

// projects implements ProjectInterface
type projects struct {
	client rest.Interface
}

// newProjects returns a Projects
func newProjects(c *KubermaticV1Client) *projects {
	return &projects{
		client: c.RESTClient(),
	}
}

// Get takes name of the project, and returns the corresponding project object, and an error if there is any.
func (c *projects) Get(name string, options metav1.GetOptions) (result *v1.Project, err error) {
	result = &v1.Project{}
	err = c.client.Get().
		Resource("projects").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Projects that match those selectors.
func (c *projects) List(opts metav1.ListOptions) (result *v1.ProjectList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ProjectList{}
	err = c.client.Get().
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested projects.
func (c *projects) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("projects").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a project and creates it.  Returns the server's representation of the project, and an error, if there is any.
func (c *projects) Create(project *v1.Project) (result *v1.Project, err error) {
	result = &v1.Project{}
	err = c.client.Post().
		Resource("projects").
		Body(project).
		Do().
		Into(result)
	return
}

// Update takes the representation of a project and updates it. Returns the server's representation of the project, and an error, if there is any.
func (c *projects) Update(project *v1.Project) (result *v1.Project, err error) {
	result = &v1.Project{}
	err = c.client.Put().
		Resource("projects").
		Name(project.Name).
		Body(project).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *projects) UpdateStatus(project *v1.Project) (result *v1.Project, err error) {
	result = &v1.Project{}
	err = c.client.Put().
		Resource("projects").
		Name(project.Name).
		SubResource("status").
		Body(project).
		Do().
		Into(result)
	return
}

// Delete takes name of the project and deletes it. Returns an error if one occurs.
func (c *projects) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("projects").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *projects) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("projects").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched project.
func (c *projects) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Project, err error) {
	result = &v1.Project{}
	err = c.client.Patch(pt).
		Resource("projects").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
