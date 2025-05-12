package argocd

import "fmt"

type (
	ArgocdApp interface {
		Sync(req *ArgocdAppReq) error
	}

	ArgocdAppReq struct {
		Name string
	}

	argocdAppHandler struct{}
)

func NewArgocdAppHandler() ArgocdApp {
	return &argocdAppHandler{}
}

func (a *argocdAppHandler) Sync(req *ArgocdAppReq) error {
	if req.Name == "" {
		return fmt.Errorf("no application name specified")
	}

	// Call argocd...

	return nil
}
