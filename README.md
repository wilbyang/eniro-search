# Eniro Search
This repo is for a rest services written in Golang, using Eniro search. 

## Installing dependencies
check [Golang installation](https://golang.org/doc/install)



### How to run the program
1. go to the repo root and run
```
go build -o eniro-search main.go eniro.go
```

2. run locally
```
./eniro-search
```

p.s. you can specify a http parameter to the commandline, like 
```
./eniro-search -http=:8081
```

### Help
if you can not run the project locally for some resone, I have deployed to Google Coud Appengine
please check [https://qieruzhengde.appspot.com/search?q=suzhi&include=companyInfo,location,address](https://qieruzhengde.appspot.com/search?q=suzhi&include=companyInfo,location,address) for playground

### Features
1. search
 
2. timeout


3. fields filtering


### Frondend
Frontend is using React and React Hooks for state management

check the repo

## Author
yang.wilby@gmail.com