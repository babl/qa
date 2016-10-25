FROM registry.babl.sh/logstash:stable
RUN /opt/logstash/bin/logstash-plugin install logstash-output-logmatic
RUN /opt/logstash/bin/logstash-plugin install logstash-output-kafka
CMD ["/bin/sh -c start"]
