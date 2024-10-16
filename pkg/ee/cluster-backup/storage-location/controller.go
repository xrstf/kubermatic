//go:build ee

/*
                  Kubermatic Enterprise Read-Only License
                         Version 1.0 ("KERO-1.0”)
                     Copyright © 2023 Kubermatic GmbH

   1.	You may only view, read and display for studying purposes the source
      code of the software licensed under this license, and, to the extent
      explicitly provided under this license, the binary code.
   2.	Any use of the software which exceeds the foregoing right, including,
      without limitation, its execution, compilation, copying, modification
      and distribution, is expressly prohibited.
   3.	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND,
      EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
      MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
      IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
      CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
      TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
      SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

   END OF TERMS AND CONDITIONS
*/

package storagelocation

import (
	"context"
	"fmt"

	velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	"go.uber.org/zap"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	"k8c.io/kubermatic/v2/pkg/ee/cluster-backup/storage-location/backupstore"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	ControllerName = "cluster-backup-storage-location-controller"
)

type reconciler struct {
	ctrlruntimeclient.Client

	recorder record.EventRecorder
	log      *zap.SugaredLogger
}

func Add(mgr manager.Manager, numWorkers int, log *zap.SugaredLogger) error {
	reconciler := &reconciler{
		Client:   mgr.GetClient(),
		recorder: mgr.GetEventRecorderFor(ControllerName),
		log:      log,
	}

	c, err := controller.New(ControllerName, mgr, controller.Options{
		Reconciler:              reconciler,
		MaxConcurrentReconciles: numWorkers,
	})
	if err != nil {
		return fmt.Errorf("failed to create controller: %w", err)
	}

	if err := c.Watch(source.Kind(mgr.GetCache(), &kubermaticv1.ClusterBackupStorageLocation{}), &handler.EnqueueRequestForObject{}); err != nil {
		return fmt.Errorf("failed to create watch for ClusterBackupStorageLocation: %w", err)
	}
	return nil
}

func (r *reconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := r.log.With("request", request)
	log.Debug("Reconciling")

	cbsl := &kubermaticv1.ClusterBackupStorageLocation{}
	if err := r.Client.Get(ctx, request.NamespacedName, cbsl); err != nil {
		log.Debug("ClusterBackupStorageLocation not found")
		return reconcile.Result{}, nil
	}

	if cbsl.Spec.Provider != "aws" {
		log.Infow("unsupported provider, skipping..", "provider", cbsl.Spec.Provider)
		return reconcile.Result{}, nil
	}
	if cbsl.Spec.Credential == nil {
		log.Info("no credentials secret reference, skipping..")
		return reconcile.Result{}, nil
	}
	err := r.reconcile(ctx, cbsl)
	if err != nil {
		r.recorder.Event(cbsl, corev1.EventTypeWarning, "ReconcilingError", err.Error())
	}
	return reconcile.Result{}, nil
}

func (r *reconciler) reconcile(ctx context.Context, cbsl *kubermaticv1.ClusterBackupStorageLocation) error {
	creds := &corev1.Secret{}
	err := r.Client.Get(ctx,
		types.NamespacedName{
			Name:      cbsl.Spec.Credential.Name,
			Namespace: cbsl.Namespace,
		},
		creds,
	)
	if err != nil {
		return fmt.Errorf("failed to get the credentials secret: %w", err)
	}
	store, err := backupstore.NewBackupStore(ctx, cbsl, creds)
	if err != nil {
		return fmt.Errorf("failed to create backup store: %w", err)
	}

	if err = store.IsValid(ctx); err != nil {
		return r.updateCBSLStatus(ctx,
			cbsl,
			velerov1.BackupStorageLocationPhaseUnavailable,
			err.Error(),
		)
	}
	return r.updateCBSLStatus(ctx,
		cbsl,
		velerov1.BackupStorageLocationPhaseAvailable,
		"ClusterBackupStoreLocation is available",
	)
}

func (r *reconciler) updateCBSLStatus(ctx context.Context, cbsl *kubermaticv1.ClusterBackupStorageLocation, phase velerov1.BackupStorageLocationPhase, message string) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := r.Client.Get(ctx,
			types.NamespacedName{
				Namespace: cbsl.Namespace,
				Name:      cbsl.Name},
			cbsl,
		); err != nil {
			return err
		}
		updatedCBSL := cbsl.DeepCopy()
		updatedCBSL.Status.Message = message
		updatedCBSL.Status.Phase = phase
		now := metav1.Now()
		updatedCBSL.Status.LastValidationTime = &now
		// we patch anyway even if there is no changes because we want to update the LastValidationTime.
		return r.Client.Status().Patch(ctx, updatedCBSL, ctrlruntimeclient.MergeFrom(cbsl))
	})
}
