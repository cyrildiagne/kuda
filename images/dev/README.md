To build & deploy this image:

```bash
cd $KUDA_DIR

docker build -f images/dev/Dockerfile -t gcr.io/kuda-project/dev:1.0.0 .
docker push $KUDA_VERSION
```