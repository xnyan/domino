echo ""
echo "`date` ==== Stoping Servers ===="
echo ""

../sbin/node.sh s -t settings.sh

sleep 2

echo ""
echo "`date` ==== Stoping Clients ===="
echo ""

../sbin/node.sh c -t settings.sh
