# Automate Http/Https Load testing
The wg/wrk extension tool provides continuous load testing in range of connection. Easy to find out the bottleneck and increase your web application performance. Using together with Jenkins are available via http call(curl to API). Understand your web application load characteristic in each number of connection.

* Latency
* Data-Transfer/Second
* Socket Error
* Non-2xx Response

### Mode
* Test by case
* Realtime Test
* Soaking Test

### require
* [wg/wrk](https://github.com/wg/wrk)
* [golang](https://golang.org/)

### Available
* Linux
* OSX

### Install Instruction
* Install Golang
* Install wrk
```
git clone https://github.com/wg/wrk.git
cd wrk
make
sudo install 0755 wrk /bin
```
* Install Ahlt
```
git clone https://github.com/tspn/ahlt.git
cd ahlt
make
```
