package helm

import (
	"encoding/json"
	"fmt"

	helmclient "github.com/kube-tarian/kad/integrator/common-pkg/plugins/helm/go-helm-client"
	"github.com/kube-tarian/kad/integrator/model"
)

func (h *HelmCLient) Delete(payload model.RequestPayload) (json.RawMessage, error) {
	h.logger.Infof("Helm client Install invoke started")

	req := &model.Request{}
	err := json.Unmarshal(payload.Data, req)
	if err != nil {
		h.logger.Errorf("payload unmarshal failed, %v", err)
		return nil, err
	}

	helmClient, err := h.getHelmClient(req)
	if err != nil {
		h.logger.Errorf("helm client initialization failed, %v", err)
		return nil, err
	}

	err = h.addOrUpdate(helmClient, req)
	if err != nil {
		h.logger.Errorf("helm repo add failed, %v", err)
		return nil, err
	}

	// Define the released chart to be uninstalled.
	chartSpec := helmclient.ChartSpec{
		ReleaseName: req.ReleaseName,
		ChartName:   fmt.Sprintf("%s/%s", req.RepoName, req.ChartName),
		Namespace:   req.Namespace,
		Wait:        true,
	}

	// Uninstall the chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	err = helmClient.UninstallRelease(&chartSpec)
	if err != nil {
		h.logger.Errorf("helm uninitialization for request %+v failed, %v", req, err)
		return nil, err
	}

	h.logger.Infof("helm uninstall of app %s successful in namespace: %v", req.ReleaseName, req.Namespace)
	h.logger.Infof("Helm client Install invoke finished")
	return json.RawMessage(fmt.Sprintf("{\"status\": \"Application %s successful with helm client\"}", req.ReleaseName)), nil
}
