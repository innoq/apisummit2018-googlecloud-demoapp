# Google Cloud Functions pandoc sample

This sample shows how to convert markdown to html using pandoc and HTTP-triggered Cloud Function.

``` 
# download pandoc binary
curl -Lo pandoc.tar.gz https://github.com/jgm/pandoc/releases/download/2.5/pandoc-2.5-linux.tar.gz
tar -zxvf pandoc.tar.gz  pandoc-2.5/bin/pandoc  --strip-components=2
rm ./pandoc.tar.gz

#deploy
gcloud functions deploy markdown_to_html --runtime python37 --trigger-http

#test
curl -vF 'doc=@test.md' https://[...].cloudfunctions.net/markdown_to_html


```


