Delete the bin folder then follow these steps - 

docker build -t api-server .

keploy test -c "docker run --name api-server --network keploy-network \
  -p 8082:8082 -e ENV=stage \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/credentials:/app/credentials \
  --rm api-server"

Try for multiple times, out.txt has the POC
    
