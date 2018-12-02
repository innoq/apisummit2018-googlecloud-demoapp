# Google Cloud Functions ImageMagick sample

This sample shows how to resize images using ImageMagick in a Storage-triggered Cloud Function.

``` 
# create bucket

export TEST_BUCKET=eimer-$RANDOM

gsutil mb "gs://$TEST_BUCKET"

gcloud functions deploy resize_image --trigger-bucket=gs://$TEST_BUCKET/ --runtime python37

gsutil cp random-image.jpg gs://$TEST_BUCKET/

gsutil ls gs://$TEST_BUCKET/

```

https://console.cloud.google.com/storage/browser/