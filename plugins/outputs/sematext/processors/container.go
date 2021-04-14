package processors

import (
	"os"

	"github.com/influxdata/telegraf"
)

const (
	CONTAINER_NAME_ENV_NAME         = "SEMATEXT_CONTAINER_NAME"
	CONTAINER_ID_ENV_NAME           = "SEMATEXT_CONTAINER_ID"
	CONTAINER_IMAGE_NAME_ENV_NAME   = "SEMATEXT_CONTAINER_IMAGE_NAME"
	CONTAINER_IMAGE_TAG_ENV_NAME    = "SEMATEXT_CONTAINER_IMAGE_TAG"
	CONTAINER_IMAGE_DIGEST_ENV_NAME = "SEMATEXT_CONTAINER_IMAGE_DIGEST"
	K8S_POD_ENV_NAME                = "SEMATEXT_K8S_POD_NAME"
	K8S_NAMESPACE_ENV_NAME          = "SEMATEXT_K8S_NAMESPACE"
	K8S_CLUSTER_ENV_NAME            = "SEMATEXT_K8S_CLUSTER"

	CONTAINER_NAME_TAG         = "container.name"
	CONTAINER_ID_TAG           = "container.id"
	CONTAINER_IMAGE_NAME_TAG   = "container.image.name"
	CONTAINER_IMAGE_TAG_TAG    = "container.image.tag"
	CONTAINER_IMAGE_DIGEST_TAG = "container.image.digest"
	K8S_POD_NAME_TAG           = "kubernetes.pod.name"
	K8S_NAMESPACE_ID_TAG       = "kubernetes.namespace"
	K8S_CLUSTER_TAG            = "kubernetes.cluster.name"
)

// ContainerTags is a metric processor that injects container tags read from env variables
type ContainerTags struct {
	tags map[string]string
}

// NewContainerTags creates new instance of container MetricProcessor
func NewContainerTags() MetricProcessor {
	tags := make(map[string]string)
	tags[CONTAINER_NAME_TAG] = os.Getenv(CONTAINER_NAME_ENV_NAME)
	tags[CONTAINER_ID_TAG] = os.Getenv(CONTAINER_ID_ENV_NAME)
	tags[CONTAINER_IMAGE_NAME_TAG] = os.Getenv(CONTAINER_IMAGE_NAME_ENV_NAME)
	tags[CONTAINER_IMAGE_TAG_TAG] = os.Getenv(CONTAINER_IMAGE_TAG_ENV_NAME)
	tags[CONTAINER_IMAGE_DIGEST_TAG] = os.Getenv(CONTAINER_IMAGE_DIGEST_ENV_NAME)

	tags[K8S_POD_NAME_TAG] = os.Getenv(K8S_POD_ENV_NAME)
	tags[K8S_NAMESPACE_ID_TAG] = os.Getenv(K8S_NAMESPACE_ENV_NAME)
	tags[K8S_CLUSTER_TAG] = os.Getenv(K8S_CLUSTER_ENV_NAME)
	return &ContainerTags{
		tags: tags,
	}
}

// Process is a method where ContainerTags processor injects container tags from env variables to metric
func (c *ContainerTags) Process(metric telegraf.Metric) error {
	for tag, value := range c.tags {
		if value != "" {
			metric.AddTag(tag, value)
		}
	}
	return nil
}

// Close clears the resources processor used, no-op in this case
func (c *ContainerTags) Close() {}
