### GitHub Actions Exporter

A utiliity to collect useful data for github actions and plot the into meaningful dashboard

## Build

```bash
    make docker-build
```
## Deploy
Build and Push docker image of the exporter and push to dockerhub
Deploy Helm chart to deploy Prometheous, grafana, github exporter and dashboards.
```bash
    make docker-push
    make deploy
```