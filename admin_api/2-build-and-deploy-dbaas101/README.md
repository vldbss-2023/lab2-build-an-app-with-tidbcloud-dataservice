# Step 2: Build and deploy dbaas101 to EKS cluster.

The following steps guides you through the process of building and deploying dbaas101 service to EKS cluster.

The source code is under `dbaas101` folder.

1. Install docker, skip if docker is already installed.

    https://docs.docker.com/engine/install/#installation

2. Create a ECR repository in your AWS account.
    ```bash
    aws ecr create-repository \
        --repository-name lab2/dbaas101 \
        --region <region>
    ```

3. Authenticate to your default registry.

    ```bash
    aws ecr get-login-password --region <region> | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.region.amazonaws.com
    ```

4. Build dbaas101 image an push to ECR. The following instructions run under `dbaas101/` folder.
    
    1. Get go dependencies.
        ```bash
        make vendor && go mod tidy
        ```
    
    2. Build image and push to ECR.
        ```bash
        GOOS=linux GOARCH=amd64 LDFLAGS="" make build
        docker build --platform=linux/amd64 -q -f Dockerfile -t lab2/dbaas101:alpha .
        docker tag lab2/dbaas101:alpha <aws_account_id>.dkr.ecr.<region>.amazonaws.com/lab2/dbaas101:alpha
        docker push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/lab2/dbaas101:alpha
        ```

    3. Deploy dbaas101 to EKS.
        ```bash
        sed -i "s/image: (v.*)/image: <your image url>/g" manifests/dbaas101-resources.yaml
        kubectl apply -f manifests/dbaas101-resources.yaml
        ```

    4. Check if dbass101 is running.
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
