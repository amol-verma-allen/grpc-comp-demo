Delete /bin folder before replicating

Steps to Replicate - 

1. Run these commands to build binary for services -

   chmod +x run.sh
   ./run/sh

2. Open Two terminals -

   First run mock-server service in one terminal -

   ./bin/mock-server

   Then run api-server service in another with keploy

   keploy record -c "./bin/api-server"

3. Open third terminal and perform this curl command -

   curl http://localhost:8082/api/taxonomy

    
    
