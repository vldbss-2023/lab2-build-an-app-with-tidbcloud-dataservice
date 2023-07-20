![image](https://github.com/vldbss-2023/lab2-build-an-app-with-tidbcloud-dataservice/assets/4160015/8ff165cc-5ac0-4e4e-ab3f-eb2cb1883069)# Step 2: Build and deploy dbaas101 to EKS cluster.

The following steps guides you through the process of building and deploying dbaas101 service to EKS cluster.

The source code is under `dbaas101` folder.
   
0. Export some env variable. `STUDENT_NAME` should be your own name
    ```bash
    export STUDENT_NAME=XXXX
    export AWS_ACCOUNT_ID=335771843383
    export REGION=ap-southeast-1
    ```

1. Install docker, skip if docker is already installed.

    https://docs.docker.com/engine/install/#installation

2. Create a ECR repository in your AWS account.
    ```bash
    aws ecr create-repository \
        --repository-name lab2/dbaas101 \
        --region ${REGION}
    ```

3. Authenticate to your default registry.

    ```bash
    aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com
    ```

4. Build dbaas101 image an push to ECR. The following instructions run under `dbaas101/` folder.
    
    1. Get go dependencies.
        ```bash
        make vendor && go mod tidy
        ```
    
    2. Build image and push to ECR.
        ```bash
        GOOS=linux GOARCH=amd64 LDFLAGS="" make build
        docker build --platform=linux/amd64 -q -f Dockerfile -t lab2/dbaas101:${STUDENT_NAME} .
        docker tag lab2/dbaas101:${STUDENT_NAME} ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/lab2/dbaas101:${STUDENT_NAME}
        docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/lab2/dbaas101:${STUDENT_NAME}
        ```

    3. Deploy dbaas101 to EKS. (Note: you need to init kubenetes configuration file at first.[click here](https://github.com/vldbss-2023/lab1-deploy-tidb-cluster-on-aws-eks/tree/main/1-create-an-eks-cluster#25-scoring-point-interact-with-the-newly-created-eks-cluster) )
        ```bash
        sed -i "s#<your image url>#${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/lab2/dbaas101:${STUDENT_NAME}#g" manifests/dbaas101-resources.yaml
        kubectl apply -f manifests/dbaas101-resources.yaml
        ```

    5. Check if dbass101 is running.
        ```bash
        kubectl get -n default deploy
        ```

5. Access dbaas101 list api.

    1. Open a new shell session and setup port forward.
        ```bash
        kubectl port-forward deploy/dbaas101 8082:8082
        ```

    2. Back to the original shell session and curl the api.
        ```bash
        curl http://localhost:8082/api/v1/tidbclusters
        ```
