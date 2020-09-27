NAME="tasarruf"
IMAGENAME="tasarruf-server"

echo "=== building docker image"

docker build -t $IMAGENAME .

echo "=== docker image built succesfully"

echo "=== Logging in to AWS ECR"

aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 177403612005.dkr.ecr.us-east-1.amazonaws.com

echo "=== Logged in to AWS ECR"

docker tag tasarruf-server:latest 177403612005.dkr.ecr.us-east-1.amazonaws.com/tasarruf-server:latest

echo "=== pushing docker image to AWS ECR"

docker push 177403612005.dkr.ecr.us-east-1.amazonaws.com/tasarruf-server:latest

echo "--- Pushed the docker image to ECR"

# echo "=== Cleaning untagged images"
#
# for image_id in $(docker images --filter "dangling=true" -q);do docker image rm $image_id;done
#
# echo "--- Cleaned old images tagged"
