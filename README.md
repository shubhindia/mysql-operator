#Mysql Operator
Just a simple mysql operator which gives us a Mysql kind in kubernetes. Once the kind is created, operator starts the mysql instance in the given namespace. This is not even in alpha stage. There are a lot of things which I need to add before it can even start a single usable mysql instance.
Currently you can pass below specs to the kind. 
1. image: <Mysql docker image>
2. size: <replicas>
3. password: <mysql password>

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
```

Note: This is something which I am building to get familier with golang and operator-sdk. Its not even in alpha stage but I will be adding more and more features to it.