FROM debian:bullseye-slim

ARG BASEDIR=/go/src/github.com/stackpulse

ARG STEP_BASEPATH
COPY ${STEP_BASEPATH}/ .

RUN apt-get update
RUN apt-get install -y python3 pip git
RUN pip install git+https://github.com/bluejeans/api-rest-meetings.git@pip-repo
RUN chmod +x doit.py

ENTRYPOINT ["./doit.py"]
