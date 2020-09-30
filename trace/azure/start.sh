echo ""
echo "`date` ==== Starting Servers ===="
echo ""

../sbin/node.sh s -s settings.sh

t=5
echo ""
echo "`date` ==== Waiting $t seconds to start clients ===="
sleep $t

echo ""
echo "`date` ==== Starting Clients ===="
echo ""

../sbin/node.sh c -s settings.sh
