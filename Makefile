install_charts:
	helm install postgres-release oci://registry-1.docker.io/bitnamicharts/postgresql -f postgres/values.yml && \
	helm repo add elastic https://helm.elastic.co && \
	helm install elasticsearch-release elastic/elasticsearch -f elk/elastic_values.yml && \
	mkdir /tmp/hostpath_pv && \
	cp elk/postgres-jdbc-42.7.4.jar /tmp/hostpath_pv/postgres-jdbc-42.7.4.jar && \
	helm install logstash-release elastic/logstash -f elk/logstash_values.yml