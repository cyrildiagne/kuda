The authentication service serves a simple static page that lets user
authenticate using [Cloud Identity Platform / Firebase Auth](https://console.cloud.google.com/marketplace/details/google-cloud-platform/customer-identity).

## Setup

First enable the [Cloud Identity Platform](https://console.cloud.google.com/marketplace/details/google-cloud-platform/customer-identity) on your project and configure at least one provider.

Add your domain to the list of authorized domain with the prefix `auth.kuda`.
For example:`auth.kuda.12.34.56.78.xip.io`.

If you're using the Google Auth Provider, you also have to configure the Oauth content screen.

Then click on "Application setup details" to find out values of the following variables.
ToS (Terms of Service) and PP (Privacy policy) urls can be left blank.

```bash
export KUDA_AUTH_API_KEY="your Firebase Auth API key"
export KUDA_AUTH_DOMAIN="your auth domain"
export KUDA_AUTH_TOS_URL="your terms and service url"
export KUDA_AUTH_PP_URL="your privacy policy url"
```

## Build and run locally (optional)

```bash
docker build \
  -t gcr.io/kuda-project/auth \
  -f install/auth/Dockerfile \
  .
```

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

## Deply

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

### Dev with Skaffold

To dev:

```
skaffold dev -f install/auth/skaffold.yaml
```

To run:

```bash
skaffold run -f install/auth/skaffold.yaml
```

### Deploy with kubectl

```bash
kubectl apply -f install/auth/service.yaml
```

## Try it

To check if your authentication module is working, first retrieve the URL of the
service:

```
kubectl get ksvc auth --namespace kuda
```

If the service was deployed successfully, you should see a ligne similar to:

```
NAME   URL                                     LATESTCREATED   LATESTREADY   READY   REASON
auth   http://auth.kuda.12.34.56.78.xip.io     auth-w4lrk      auth-w4lrk    True
```

You can then click on the URL and try to login.
