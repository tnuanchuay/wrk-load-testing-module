#!/bin/sh

git clone https://github.com/wg/wrk.git
cd wrk
make
ls -s wrk /bin/wrk
cd ..
./install-dependency.sh
./run.sh


