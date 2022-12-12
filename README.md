# Specification :notebook:
The aim of the project is to realise a distributed application in the Go programming language that implements the k-means clustering algorithm in a distributed version according to the MapReduce computing paradigm. 

The distributed application must fulfil the requirements listed below:
1. **Master-Worker structure**: Master distributes the workload among the nodeworkers, who implement the mappers and reducers.
2. **MapReduce paradigm**:
3. Deploy in **Docker containers** inside **EC2** instance
4. **Fail tolerant**

Infrastructure:
![f2-mapreduce-diagram](https://user-images.githubusercontent.com/50273854/206888337-7c0116c8-f17f-432d-8454-de26cde7667b.png)





# How to :computer:

IMMAGINE MAP REDUCE 

## AWS EC2 instance with Docker Containers
Ansible service is used to automate the Docker installation and to copy the application code
Aws Cli needs to be already configured.
```
# To check the ec2 instance (optional): Inside ansible folder (set personal ip and key)
ansible -i hosts.ini -m ping all

# To execute ansible in ansible: Inside ansible folder
ansible-playbook -v -i hosts.ini deploy.yaml 

# Connect to the EC2 instance 
ssh -i "key.pem" ubuntu@ip_instance 

# enter code folder
cd code

# To execute the whole application: inside code folder. Threshold = 1 is 0.0.001
./start.sh [nMaps] [Threshold] [MaxIters] 
#or
sudo NUMMAP=[nMaps] MAXITER=[Threshold] THRESHOLD=[MaxIters]  docker-compose up
```

## Client execution
```
# Inside the points folder are present some example dataset
Go run client [k] [pathPoints] [ip master]
```
example: 
```go run client 10 ./points/rand10000.txt 0.0.0.0```

## Test execution: Inside test folder
```Go run test.go [ip master]```




# Results :chart_with_upwards_trend:
The figure shown presents the algorithm running on
10 000 random points. The 3 curves have different number
of clusters: k = 20 (yellow), k = 10 (blue), k = 5 (red) the
execution times are shown in ms according to the number of
mappers (from 1 to 8) As expected, there is a substantial
increase in performance as the number of mappings
increases.
1. From 1 to 3 mappers the performance increases by
about 200%.
2. From 1 to 8 mappers the performance increases by
about 390%.

Moreover, as it was easy to expect, there is a degradation of
performance as the number of clusters increases. Doubling
the number of clusters would also seem to double the
execution times.
![f9-test-chart](https://user-images.githubusercontent.com/50273854/206888279-11031146-7dd7-45b0-bc9c-962824746334.png)

