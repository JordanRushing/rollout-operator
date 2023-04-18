package admission

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	LastDownscaleAnnotationKey = "grafana.com/last-prepared-for-downscale"
)

func addPreparedForDownscaleAnnotationToPod(ctx context.Context, api kubernetes.Interface, namespace, stsName string, podNr int) error {
	client := api.CoreV1().Pods(namespace)
	labelSelector := v1.LabelSelector{MatchLabels: map[string]string{
		"name":                               stsName,
		"statefulset.kubernetes.io/pod-name": fmt.Sprintf("%v-%v", stsName, podNr),
	}}
	pods, err := client.List(ctx,
		v1.ListOptions{LabelSelector: labels.Set(labelSelector.MatchLabels).String()})
	if err != nil {
		return err
	}
	if len(pods.Items) != 1 {
		return fmt.Errorf("multiple or no pods found for statefulset %v and index %v", stsName, podNr)
	}

	pod := pods.Items[0]

	annotations := pod.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	// The value of the annotation is not important. It is set to the current date and time.
	// This is to the benefit of the operator only.
	annotations[LastDownscaleAnnotationKey] = time.Now().UTC().String()
	pod.SetAnnotations(annotations)

	_, err = client.Update(ctx, &pod, v1.UpdateOptions{})
	return err
}

func addDownscaledAnnotationToStatefulSet(ctx context.Context, api kubernetes.Interface, namespace, stsName string) error {
	client := api.AppsV1().StatefulSets(namespace)
	sts, err := client.Get(ctx, stsName, v1.GetOptions{})
	if err != nil {
		return err
	}
	annotations := sts.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations[LastDownscaleAnnotationKey] = time.Now().UTC().String()
	sts.SetAnnotations(annotations)

	_, err = client.Update(ctx, sts, v1.UpdateOptions{})
	return err
}
