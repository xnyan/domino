echo ""
echo "`date` ==== Collecting logs from remote servers and clients ===="
echo ""

## Parallel
#../sbin/log.sh -d settings.sh
## Sequential
../sbin/log.sh -ds settings.sh
