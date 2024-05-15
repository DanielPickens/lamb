---
meta:
  - name: description
    content: "daniel pickens lamb | Documentation FAQ"
---
## Frequently Asked Questions

### I updated my deployment method to use the new API version and lamb doesn't report anything but kubectl still shows the old API. What gives?

See above in the [Purpose](/#purpose) section of this doc. Kubectl is likely lying to you because it only tells you what the default is for the given kubernetes version even if an object was deployed with a newer API version.

### I don't use helm, how can I do in cluster checks?

Currently, the only in-cluster check we are confident in supporting is helm. If your deployment method can generate yaml manifests for kubernetes, you should be able to use the `detect` or `detect-files` functionality described below after the manifest files have been generated.

### I updated the API version of an object, but lamb still reports that the apiVersion needs to be updated.

lamb looks at the API Versions of objects in releases that are in a `Deployed` state, and Helm has an issue where it might list old revisions of a release as still being in a `Deployed` state. To fix this, look at the release revision history with `helm history <release name>`, and determine if older releases still show a `Deployed` state. If so, delete the Helm release secret(s) associated with the revision number(s). For example, `kubectl delete secret sh.helm.release.v1.my-release.v10` where `10` corresponds to the release number. Then run lamb again to see if the object has been removed from the report.

### I updated the API Type calls of an object, but lamb still reports that the apiType needs to be updated.

lamb looks at the ApiTypes of objects in it's releases that are in a `deployed` state, and Helm has a similar issue where it will list old revisions of a release that is still currently in `deployed` state. To fix this issue, delete the Helm release secret(s) associated with the new revision numbers(s). For example `kubectl delete secret sh.helm.release.v2.mynew-release.v20` where `20` corresponds to the release number. You can then test the api call again in swagger to see if lamb will detect if the object has been removed for the revision.

### Why API is version check on a live cluster using the "last-applied-configuration" annotation not reliable?

When using `--detect-api-resources` or `--detect-all-in-cluster`, there are some potential issues to be aware of:

  * The annotation `kubectl.kubernetes.io/last-applied-configuration` on an object in your cluster holds the API version by which it was created. In fact, others have pointed out that updating the same object with `kubectl patch` will **remove** the annotation. Hence this is not a reliable method to detect deprecated API's from a live cluster.
  * You may get false positives in the first change after fixing the apiVersion. Please see [this issue](https://github.com/danielpickens/lamb/issues/495) for more details.
