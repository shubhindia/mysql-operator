# **Mysql Operator**
Just a simple mysql operator which gives us a Mysql kind in kubernetes. Once the kind is created, operator starts the mysql instance in the given namespace. This is not even in alpha stage. There are a lot of things which I need to add before it can even start a single usable mysql instance.
Currently you can pass below specs to the kind. 
1. image: "Mysql docker image"
2. size: "replicas"
3. password: "mysql password"
4. usepvc: "True/False"
5. pvcsize: 5Gi

**Operator Installation**
1. clone this repo using ``` git clone https://github.com/shubhindia/mysql-operator.git ```
2. ```kubectl applf -f deploy/```

**How-To**
You can refer the yaml given below:
```
apiVersion: apps.shcn.me/v1
kind: Mysql
metadata:
  name: mysql-sample
spec:
  image: "mysql:5.6"
  size: 1
  password: "password"
  usepvc: true
  pvcsize: 5Gi
```

**Bugs**
1. Seting replica size more than 1 causes issues. 
2. You have to use PVC for now. Setting usepvc to false will break casuse issues.


*Note*: This is something which I am building to get familier with golang and operator-sdk. Its not even in alpha stage but I will be adding more and more features to it.