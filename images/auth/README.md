```bash
export KUDA_AUTH_API_KEY="your auth API key"
export KUDA_AUTH_DOMAIN="your auth domain"
export KUDA_AUTH_TOS_URL="your terms and service url"
export KUDA_AUTH_PP_URL="your privacy policy url"
```

## Build

```bash
docker build \
  -t gcr.io/kuda-project/auth \
  -f ./Dockerfile \
  .
```

## Run

```bash
docker run --rm \
  -e KUDA_AUTH_API_KEY=$KUDA_AUTH_API_KEY \
  -e KUDA_AUTH_DOMAIN=$KUDA_AUTH_DOMAIN \
  -e KUDA_AUTH_TOS_URL=$KUDA_AUTH_TOS_URL \
  -e KUDA_AUTH_PP_URL=$KUDA_AUTH_PP_URL \
  -e PORT=80 \
  -p 8080:80 \
  gcr.io/kuda-project/auth
```

## Deploy

```bash
KUDA_AUTH_TOS_URL=$(echo $KUDA_AUTH_TOS_URL | sed 's/\//\\\//g')
KUDA_AUTH_PP_URL=$(echo $KUDA_AUTH_PP_URL | sed 's/\//\\\//g')
cp service.tpl.yaml service.yaml
sed -i'.bak' "s/value: <your-auth-api-key>/value: $KUDA_AUTH_API_KEY/g" service.yaml
sed -i'.bak' "s/value: <your-auth-domain>/value: $KUDA_AUTH_DOMAIN/g" service.yaml
sed -i'.bak' "s/value: <your-tos-url>/value: $KUDA_AUTH_TOS_URL/g" service.yaml
sed -i'.bak' "s/value: <your-pp-url>/value: $KUDA_AUTH_PP_URL/g" service.yaml
rm service.yaml.bak
```
