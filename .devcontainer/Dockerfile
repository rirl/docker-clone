ARG REGISTRY=registry.hub.docker.com
ARG IMAGE=library
ARG VARIANT=golang:alpine3.19
FROM ${REGISTRY}/${IMAGE}/${VARIANT}

ARG USERNAME=rirl
ARG USER_UID=1000
ARG USER_GID=${USER_UID}
ARG DOCKER_GROUP=docker
ARG DOCKER_GID=1001

# Create the user, install sudo, and enable password-less login
RUN apk update && apk add --no-cache shadow \ 
    && groupadd --force --non-unique --gid $DOCKER_GID $DOCKER_GROUP \
    && groupadd --force --non-unique --gid $USER_GID $USERNAME \
    && useradd --non-unique --uid $USER_UID --gid $USER_GID -m $USERNAME \
    && usermod -G $DOCKER_GROUP $USERNAME \
    && apk add --no-cache sudo bash git \
    && mkdir -p /etc/sudoers.d && echo "$USERNAME  ALL=(ALL) NOPASSWD:ALL">> /etc/sudoers.d/$USERNAME \
    && apk add --no-cache docker-cli

# [Optional] Set the default user. Omit if you want to keep the default as root.
USER $USERNAME

ENV PATH="/usr/bin:$PATH"