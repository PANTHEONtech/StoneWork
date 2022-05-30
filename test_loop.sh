RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

STEP=1
while true; do
  echo "***********************  STEP ************************ $STEP"
#  export TEST_LOOP_STATE=0
  make test
  RES=$?
#  export TEST_LOOP_STATE=1
  echo "$RES"
  if [ $RES -ne 0 ]; then
    echo "Fail"
    echo "------------------------------------------------"
    echo -e "${RED}[FAIL] ${3}${NC}"
    echo "------------------------------------------------"
    #    kill -s TERM $TOP_PID
    break
  else
    echo "OK"
  fi
  if [ $STEP -ge $1 ]; then
    break
  fi
  STEP=$((STEP + 1))
done
