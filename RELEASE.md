# How to release a new version

1. Update `CHANGELOG.md`
  - Open PR and get it merged
1. Create a new tag that follows semantic versioning:
    ```bash
    $ tag=v0.1.0
    $ git tag -s "${tag}" -m "${tag}"
    $ git push origin "${tag}"
    ```
1. Publish the updated Docker image
    ```bash
    $ IMAGE_TAG="${tag}" make publish-images
    ```
1. Update Helm Chart
  - Repository https://github.com/grafana/helm-charts/tree/main/charts/rollout-operator
  - [Example PR](https://github.com/grafana/helm-charts/pull/3177/files)