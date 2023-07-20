# Step 3: Implement the create handler and deploy the new dbaas101.


The following steps guides you through the process of building and deploying dbaas101 service to EKS cluster.

The source code is under `dbaas101` folder.

1. Implement the `UpdateTidbCluster` function
Read the code in `dbaas101/pkg/service/api.go` and implement the `UpdateTidbCluster` function.

You can read the `CreateTidbCluster`, `GetTidbCluster`, `ListTidbCluster` and `DeleteTidbCluster` to learn how to use 
the `client.Client` to interact with k8s apiserver.

Commit your implementation.


Other references:

- https://github.com/kubernetes-sigs/controller-runtime/tree/main/examples
- https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/client


2. Deploy dbaas101 again

    1. build a new image and push to ECR.
        ```bash
        GOOS=linux GOARCH=amd64 LDFLAGS="" make build
        docker build --platform=linux/amd64 -q -f Dockerfile -t lab2/dbaas101:latest .
        docker tag lab2/dbaas101:latest <aws_account_id>.dkr.ecr.<region>.amazonaws.com/lab2/dbaas101:$(git rev-parse --short HEAD)
        docker push <aws_account_id>.dkr.ecr.<region>.amazonaws.com/lab2/dbaas101:$(git rev-parse --short HEAD)
        ```

    2. Deploy your dbaas101 to EKS
        ```bash
        sed -i "s/image: (v.*)/image: <your image url>/g" manifests/dbaas101-resources.yaml
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
    
