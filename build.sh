#!/bin/bash

APP_NAME=go-meshblu-http-server
TMP_DIR=/tmp/$APP_NAME/
IMAGE_NAME=local/$APP_NAME

build() {
  docker build --tag $IMAGE_NAME .
}

copy() {
  cp $TMP_DIR/$APP_NAME .
  cp $TMP_DIR/$APP_NAME entrypoint/
}

init() {
  rm -rf $TMP_DIR/ \
   && mkdir -p $TMP_DIR/
}

package() {
  docker build --tag $IMAGE_NAME:entrypoint entrypoint
}

run() {
  docker run --rm \
    --volume $TMP_DIR:/export/ \
    $IMAGE_NAME \
      cp $APP_NAME /export
}

panic() {
  local message=$1
  echo $message
  exit 1
}

main() {
  init    || panic "init failed"
  build   || panic "build failed"
  run     || panic "run failed"
  copy    || panic "copy failed"
  package || panic "package failed"
}
main
