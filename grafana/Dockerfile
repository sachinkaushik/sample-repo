
FROM grafana/grafana:9.1.2

USER root

RUN grafana-cli --pluginsDir "/var/lib/grafana/plugins" plugins install yesoreyeram-infinity-datasource
RUN grafana-cli --pluginsDir "/var/lib/grafana/plugins" plugins install grafana-simple-json-datasource

COPY provisioning/datasources.yml /etc/grafana/provisioning/datasources/datasources.yml
COPY provisioning/dashboards.yml /etc/grafana/provisioning/dashboards/dashboards.yml
COPY provisioning/dashboard.json /etc/dashboards/dashboard.json

USER grafana