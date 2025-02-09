package metrics

import (
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"kubevirt.io/ssp-operator/internal/common"
	"kubevirt.io/ssp-operator/internal/operands"
)

// Define RBAC rules needed by this operand:
// +kubebuilder:rbac:groups=monitoring.coreos.com,resources=prometheusrules,verbs=get;list;watch;create;update;patch;delete

type metrics struct{}

func (m *metrics) Name() string {
	return operandName
}

func (m *metrics) AddWatchTypesToScheme(scheme *runtime.Scheme) error {
	return promv1.AddToScheme(scheme)
}

func (m *metrics) WatchTypes() []client.Object {
	return []client.Object{&promv1.PrometheusRule{}}
}

func (m *metrics) WatchClusterTypes() []client.Object {
	return nil
}

func (m *metrics) Reconcile(request *common.Request) ([]common.ResourceStatus, error) {
	return common.CollectResourceStatus(request,
		reconcilePrometheusRule,
	)
}

func (m *metrics) Cleanup(*common.Request) error {
	return nil
}

var _ operands.Operand = &metrics{}

func GetOperand() operands.Operand {
	return &metrics{}
}

const (
	operandName      = "metrics"
	operandComponent = common.AppComponentMonitoring
)

func reconcilePrometheusRule(request *common.Request) (common.ResourceStatus, error) {
	return common.CreateOrUpdate(request).
		NamespacedResource(newPrometheusRule(request.Namespace)).
		WithAppLabels(operandName, operandComponent).
		UpdateFunc(func(newRes, foundRes client.Object) {
			foundRes.(*promv1.PrometheusRule).Spec = newRes.(*promv1.PrometheusRule).Spec
		}).
		Reconcile()
}
