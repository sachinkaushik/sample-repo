FROM intel/edge_video_analytics_microservice:0.7.2
USER root
COPY . /app
RUN pip3 install -r /app/requirements.txt -t /usr/local/lib/python3.8/dist-packages
USER eiiuser
ENTRYPOINT [ "./run.sh"]