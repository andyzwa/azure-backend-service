#!/bin/bash

set -eo pipefail

#Local appname
APP="azure-backend-service"
IMAGE_TAG="latest"

#Image name in docker hub
IMAGE_NAME="u115725/$APP"

#Call and log
callAndLog() {
  echo "${1}"
  ${1}
}

#Build Image
build() {

  # build service
  env GOOS=linux GOARCH=amd64 go build -o ./app/myservice.bin

  #Go settings
  echo "# Building go application and docker image ${APP}:${IMAGE_TAG}"
  callAndLog "docker build -t $APP:$IMAGE_TAG ."

}

test(){

  #run docker
  callAndLog "docker run --rm -d --name $APP -p 8000:8000 $APP"

  #check return
  callUrl=http://localhost:8000/
  set +e
  countPong=$(curl -s $callUrl | grep -c "Welcome to the HomePage!")
  set -e

  #stop docker
  callAndLog "docker stop $APP"

  if (( countPong > 0 )) ; then
      echo "OK test call url: $callUrl"
      return 0
  else
      echo "ERROR test call url: $callUrl"
      return 1
  fi
}

tag() {
  echo "# Tagging image for docker hub"
  callAndLog "docker tag $APP:$IMAGE_TAG $IMAGE_NAME:$IMAGE_TAG"
}

push() {
  echo "# Pushing image -> Maybe you have to login to docker hub first: docker login"
  callAndLog "docker push $IMAGE_NAME:$IMAGE_TAG"
}


#Tag the local image and push it
tagpush() {
  tag
  push
}

#Run the local image
run() {
  echo "run local docker.."
  echo "GET Articles: http://localhost:8000/articles"
  echo "SwaggerUI: http://localhost:8000/swaggerui/"
  docker run --rm -it --name $APP -p 8000:8000 $APP
}

#Run curl commands
get()
{
   echo "Get articles"
   curl --header "Content-Type: application/json" --request GET http://localhost:8000/articles
}

add() {
   echo "Add article"
   curl \
    -H "Accept: application/json" \
    -X POST --data '{"id":"","title":"Article 3","content":"Article Description 3","color":"#fff2"}' \
    http://localhost:8000/article

}

del() {
   echo "Delete article"
       curl -i \
    -X DELETE http://localhost:8000/article/3
}

#Validate command
if ! [[ "$1" =~ ^(build|test|tag|push|tagpush|run|get|add|del)$ ]]; then
    echo "No valid target given, possible values: build|test|tag|push|tagpush|run"
    exit 1
fi

# call the command
$1

