MySQL Operator for Kubernetes
=============================

The MYSQL Operator for Kubernetes is an Operator for Kubernetes managing
MySQL Instance setups inside a Kubernetes Cluster.

Release Status
--------------
The MySQL Operator for Kubernetes currently is in a preview state.
DO NOT USE IN PRODUCTION.

Installation of the MySQL Operator
----------------------------------

The MYSQL Operator can be installed using `kubectl`:

```sh
kubectl apply -f https://raw.githubusercontent.com/shubhindia/mysql-operator/develop/deploy/mysql-operator.yaml
```

Note: The propagation of the CRDs can take a few seconds depending on the size
of your Kubernetes cluster. Best is to wait a second or two between those
commands. If the second command fails due to missing CRD apply it a second
time.

To verify the operator is running check the deployment managing the 
operator, inside the `mysql-operator` namespace.

```sh
kubectl get po -n mysql-operator-system
```

Once the Operator is ready the output should be like

``` 
NAME                                                 READY   STATUS    RESTARTS   AGE
mysql-operator-controller-manager-5c4b67bbc8-q6g6h   2/2     Running   0          66s
```

Using the MySQL Operator to setup a MySQL Instance
-------------------------------------------------------

For creating a Mysql instance you need to use/refer below yaml

```
apiVersion: apps.shubhindia.me/v1beta1
kind: Mysql
metadata:
  name: mysql-sample
spec:
  usepvc: true
  pvcspec:
    name: test
    size: 1Gi
    storageclass: standard
```

With that the sample cluster can be created:

```sh
kubectl apply -f mysql-instance.yaml
```

This sample will create a Mysql instance

```sh
kubectl get mysql --watch
```

Once the instance is ready the output should be like
```
NAME           AGE   STATUS
mysql-sample   5s    Ready
```
Connecting to the MYSQL Instance
-------------------------------------

For connecting to the Mysql instance a `Service` is created inside the 
Kubernetes cluster.

```sh
kubectl get service 
```

``` 
NAME           TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
mysql-sample   NodePort    10.111.90.125   <none>        3306:30942/TCP   80s
```

Get instance password from secret
```
kubectl get secret mysql-sample-user-password --template={{.data.password}} | base64 -d
```
And then in a second terminal:

```sh
mysqlsh -h <node-ip> -P <node-port> -u root -p
```

When promted enter the password, which received from secret.

