FROM registry-gitlab.i.wish.com/contextlogic/protobuf-go-base-image/release:v0.0.1 as build

ARG ITA_JOB_TOKEN
ARG GITLAB_ACCESS_TOKEN=${ITA_JOB_TOKEN}
ARG ITA_PROJECT_NAME

# Validating the Arguments
RUN if [[ -z "$GITLAB_ACCESS_TOKEN" ]]; then \
    echo -e "\n[ Error ] GITLAB_ACCESS_TOKEN is missing\n"; \
    exit 1; \
    fi; \
    if [[ -z "$ITA_PROJECT_NAME" ]]; then \
    echo -e "\n[ Error ] ITA_PROJECT_NAME is missing. This should be the same as the value of your project\n"; \
    exit 1; \
    fi

COPY . /go/src/github.com/ContextLogic/${ITA_PROJECT_NAME}
WORKDIR /go/src/github.com/ContextLogic/${ITA_PROJECT_NAME}

RUN git config --global url.https://gitlab-ci-token:${GITLAB_ACCESS_TOKEN}@gitlab.i.wish.com/ContextLogic.insteadof 'https://github.com/ContextLogic'
RUN echo ${ITA_PROJECT_NAME}
ENV GOSUMDB off
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o bin/${ITA_PROJECT_NAME} github.com/ContextLogic/${ITA_PROJECT_NAME}

# Temporarily switcing to git to fix docker rate-limiting
FROM registry-gitlab.i.wish.com/contextlogic/tooling-image/alpine:3.7

ARG ITA_PROJECT_NAME
ENV BIN_NAME=/bin/${ITA_PROJECT_NAME}

# install timezone data in alpine which enables timezone conversion.
RUN apk add tzdata
ENV ZONEINFO=/usr/share/timezone

COPY --from=build /go/src/github.com/ContextLogic/${ITA_PROJECT_NAME}/bin/${ITA_PROJECT_NAME} /bin/${ITA_PROJECT_NAME}
COPY /config/service.json /config/service.json

# Exposing the ports, so the tests in Gitlab can work
EXPOSE 8080
EXPOSE 8081

CMD [ "${BIN_NAME}" ]
