
t=$1
if [ -z $t ]; then
  echo "Usage: <exp_time in seconds>"; exit 1
fi

./build.sh

./deploy.sh client server exp.config location.config

./start.sh

echo "`date` ==== Wait $t seconds for the experiment to be done ===="
#t=$((15*1))
sleep $t
echo "`date` ==== The experiment should be done after $t seconds ===="

./process.sh

echo "`date` ==== Important: logs NOT collected ===="

exit 1

#./log-check.sh
./log-collect.sh
#./stop.sh
#./log-rm.sh
