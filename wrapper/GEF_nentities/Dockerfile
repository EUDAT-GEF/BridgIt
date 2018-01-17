FROM weblicht.sfs.uni-tuebingen.de:5000/clarind/wl-webservices/service-illinois-ner

MAINTAINER Claus Zinn <claus.zinn@uni-tuebingen.de>

LABEL "eudat.gef.service.name"="Illinois NER (GEF)"
LABEL "eudat.gef.service.description"="Wrapper around the Illinois Named Entity Recognizer"
LABEL "eudat.gef.service.version"="1.0"
LABEL "eudat.gef.service.input.1.name"="Input Directory"
LABEL "eudat.gef.service.input.1.path"="/root/input"
LABEL "eudat.gef.service.output.1.name"="Output Directory"
LABEL "eudat.gef.service.output.1.path"="/root/output"


# install python
RUN apk add --update \
    python \
    python-dev \
    py-pip \
    build-base \
  && pip install requests \
  && rm -rf /var/cache/apk/*

ADD wrapper_gef.py wrapper/wrapper_gef.py

WORKDIR /wrapper

ENV PYTHONUNBUFFERED=0
ENV _JAVA_OPTIONS="-Xmx9g"
ENV JAVA_OPTS="-Xmx9g"
ENTRYPOINT ["/usr/bin/python", "wrapper_gef.py"]



