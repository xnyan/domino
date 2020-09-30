
echo ""
echo "`date` ==== Checking Server Processes ===="
echo ""

#Checks server process
../sbin/node.sh s -c settings.sh

echo ""
echo "`date` ==== Checking Client Processes ===="
echo ""

#Checks client process
../sbin/node.sh c -c settings.sh
