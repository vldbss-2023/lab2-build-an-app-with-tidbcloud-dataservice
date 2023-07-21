# Step 3: Implement the create handler and deploy the new dbaas101.


The following steps guides you through the process of building and deploying dbaas101 service to EKS cluster.

The source code is under `dbaas101` folder.

1. Implement the `CreateTidbCluster` function
Read the code in `dbaas101/pkg/service/api.go` and implement the `CreateTidbCluster` function.

You can read the `GetTidbCluster`, `ListTidbCluster` and `DeleteTidbCluster` to learn how to use 
the `client.Client` to interact with k8s apiserver.

Commit your implementation.


Other references:

- https://github.com/kubernetes-sigs/controller-runtime/tree/main/examples
- https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client


2. Deploy dbaas101 again
    0. Export some env variable. `STUDENT_NAME` should be your own name
        ```bash
         export STUDENT_NAME=XXXX
         export AWS_ACCOUNT_ID=335771843383
         export REGION=ap-southeast-1
        ```

    1. build a new image and push to ECR.
        ```bash
        GOOS=linux GOARCH=amd64 LDFLAGS="" make build
        docker build --platform=linux/amd64 -q -f Dockerfile -t lab4/dbaas101:${STUDENT_NAME}_modify .
        docker tag lab4/dbaas101:${STUDENT_NAME} ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/lab4/dbaas101:${STUDENT_NAME}_modify
        docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/lab4/dbaas101:${STUDENT_NAME}_modify
        ```

    2. Deploy your dbaas101 to EKS
        ```bash
        sed -i "s#image: .*#image: ${AWS_ACCOUNT_ID}.dkr.ecr.${REGION}.amazonaws.com/lab4/dbaas101:${STUDENT_NAME}_modify#g" manifests/dbaas101-resources.yaml
        kubectl apply -f manifests/dbaas101-resources.yaml
        ```

3. Create a new tidb cluster via api
    
    1. Setup port forward in a different shell session.
        ```bash
        kubectl port-forward deploy/dbaas101 8082:8082
        ```

    2. use curl to access create api.

        ```bash
        curl -XPOST http://localhost:8082/api/v1/tidbclusters -H "Content-Type: application/json" -d @tidb-cluster.json
        ```
    
