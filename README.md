Reference a pubblicazioni simili (chiamati in introduzione numerazione)  ●
Citare hadoop
Citare github codice cluster
Citare slides professoressa 

AGGIUNGERE TRACCIA SU PROGETTO ● 

- IMMAGINE MAPREDUCE ● 

- CONSOLE E IMG CONSOLE  ●
Mostr. Clusters, num. Iterations (master), time (client)

- SETUP EC2 INSTANCE
init.sh
Insert ip, keyPath 

- CONSOLE DISPLAY. 
         Client
Display intro luccicante
Get K 
Get PathPoints
Show result 

        Test
Display intro luccicante
Select default or personalized
Attention to select nmap or k to large. It may take too much time.
Get nMap(from 1 to inserted value/ must be minor than the initialized)
Get k (from 1 to inserted value)
Get path
Get threshold
Get MaxIter 

To generate new file points (inside points)
Go run genPoints.go number name
Es. Go run genPoints.go 1000 p1000


How to:
IMMAGINE MAP REDUCE 

Ansible service is used to automate the Docker installation and to copy the application code 

# To check the ec2 instance (optional) ansible -i hosts.ini -m ping all # To execute 

ansible in SDCC/sdcc/ansible ansible-playbook -v -i hosts.ini deploy.yaml 
# Connect to the EC2 instance 

ssh -i "key.pem" ubuntu@ip_instance 

# To execute the whole application. Inside the points folder are present some example dataset
./start.sh [nMaps] [Threshold] [MaxIters] 

Client execution:
Go run client [k] [pathPoints] 

Test execution (inside test path):
Go run test.go [nMap] [threshold] [m
MaxIter]
---- TEST CHART      ●
---- COMMENTI TEST      ●
