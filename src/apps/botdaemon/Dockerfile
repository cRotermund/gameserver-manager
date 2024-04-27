FROM python:3-alpine AS base

FROM base AS prepare

RUN apk update && apk add bash
SHELL ["/bin/bash", "-c"]

#Install key build utilities to base image if not there
RUN pip install pipx
ENV PATH=~/root/.local/share/pipx:~/.local/bin:$PATH
RUN pipx install poetry

FROM prepare as build

#Assemble lib package, install, clean
WORKDIR /tmp/appbuild
COPY ./libs/gsmclient ./gsmclient
WORKDIR ./gsmclient
RUN poetry build
RUN find ./dist -name *.whl | xargs pip install

#Assemble daemon
WORKDIR /tmp/appbuild
COPY ./apps/botdaemon ./botdaemon
WORKDIR ./botdaemon
RUN poetry build
RUN find ./dist -name *.whl | xargs pip install

#Clean up after the build, runtime is ready to execute
FROM build AS final
WORKDIR /
RUN rm -rf /tmp/appbuild

CMD ["botdaemon"]