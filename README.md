# gcpboxtest

https://github.com/sinmetalcraft/gcpbox を App Engine 上で試しに動かすもの

## [Special URLs and headers](https://cloud.google.com/iap/docs/special-urls-and-headers-howto)

## Deploy

```
gcloud app deploy .
```

## Test

### Cloud Tasks

`/cloudtasks/appengine/add-task` にアクセスすると、TaskがAddされ、 `/cloudtasks/appengine/json-post-task` が呼び出される

### Cloud Storage Pub/Sub Notify