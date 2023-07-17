# Step 1: Create an EKS cluster and deploy tidb-operator


After going thought [lab1 step1](https://github.com/vldbss-2023/lab1-deploy-tidb-cluster-on-aws-eks/tree/main/1-create-an-eks-cluster) 
and [lab1 step2](https://github.com/vldbss-2023/lab1-deploy-tidb-cluster-on-aws-eks/tree/main/2-deploy-tidb-with-tidb-operator), 
you already know how to create a eks cluster and deploy tidb-operator and tidb cluster into your EKS cluster.

If you already have an EKS cluster with tidb-operator and tidb cluster deployed, you can skip this step and start step2. 
You can refer the instructions and the pulumi codes in `lab1-step1` and `lab2-step2` to recreate an EKS cluster with tidb-operaor and tidb cluster deployed.

The following instructions are just summary from `lab1-step1` and `lab1-step2`, you may refer the original guideline for more details.

If you have already installed and configured the AWS CLI and Pulumi CLI, you can skip the setup steps.

## Create an EKS cluster

1. Set up AWS CLI, skip if you already installed and configured the AWS CLI.

    1. Install AWS CLI

        - https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
        - `[aws-iam-authenticator](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html)`:
          Amazon EKS uses IAM to provide secure authentication to your Kubernetes cluster.

    2. Config AWS credentials

        Run this command to quickly set and view your credentials, region, and output format. The following example shows
        sample values.

        ```bash
        $ aws configure
        AWS Access Key ID [None]: XXXXXXXXX
        AWS Secret Access Key [None]: XXXXXXXXX
        Default region name [None]: us-west-2
        Default output format [None]: yaml
        ```

2. Set up Pulumi, skip if you already installed and configured the Pulumi CLI.

    1. Install Pulumi

        https://www.pulumi.com/docs/get-started/install/

    2. Initialize Pulumi Stack

        ```bash
        $ pulumi login --local
        $ export PULUMI_CONFIG_PASSPHRASE="" # Set passphrase env to `""`. This passphrase is required by Pulumi and was created by Lab maintainer.
        ```

3. Create the EKS cluster via Pulumi (may take more than **_10_** minutes)
    
    following instructions should be run under `lab1-deploy-tidb-cluster-on-aws-eks/1-create-an-eks-cluster` folder

    ```bash
    $ pulumi stack select default -c # Select the `default` stack.
    $ pulumi up
    Updating (default):

         Type                      Name                            Status
    +   pulumi:pulumi:Stack        1-create-an-eks-cluster-default created
    +   └─ eks:index:Cluster       my-eks                          created
        ... dozens of resources omitted ...
    ```

4. Interact with the newly created EKS cluster

    ```bash
    $ pulumi stack output kubeconfig > kubeconfig.yaml
    $ export KUBECONFIG=$PWD/kubeconfig.yaml
    $ kubectl get nodes
    NAME                                            STATUS   ROLES    AGE   VERSION
    ip-xxx-xxx-xxx-xxx.us-west-2.compute.internal   Ready    <none>   27m   v1.27.1-eks-2f008fe
    ```

## Delpoy tidb operator and tidb cluster

1. Install Helm

    1. Download
        ```bash
        $ curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
        $ chmod 700 get_helm.sh
        $ ./get_helm.sh
        ```
    2. Init

        ```bash
        $ helm init --client-only
        ```
2. Deploy TiDB Operator and TiDB Cluster via Pulumi

    following instructions should be run under `lab1-deploy-tidb-cluster-on-aws-eks/2-deploy-tidb-with-tidb-operator` folder

    ```bash
    $ pulumi stack select default -c # Select the `default` stack.
    $ pulumi up
    Updating (default):
         Type                                                                  Name                                      Status
         pulumi:pulumi:Stack                                                   2-deploy-tidb-with-tidb-operator-default
     +-  ├─ kubernetes:helm.sh/v3:Release                                      tidb-operator                             craeted (22s)
         ├─ kubernetes:yaml:ConfigGroup                                        tidb-operator-crds
         │  └─ kubernetes:yaml:ConfigFile                                      crds/tidb-operator-v1.4.4.yaml
     +-  │     ├─ kubernetes:apiextensions.k8s.io/v1:CustomResourceDefinition  tidbinitializers.pingcap.com              craeted (2s)
         ... dozens of resources omitted ...
    ```
