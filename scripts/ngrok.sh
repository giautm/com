#!/bin/sh

echo "Start ngrok in background on port [ $PORT ]"
mkdir -p logs
ngrok http ${PORT} -bind-tls=true --log=stdout > logs/ngrok.log &

echo -n "Extracting ngrok public url ."
NGROK_PUBLIC_URL=""
while [ -z "$NGROK_PUBLIC_URL" ]; do
  # Run 'curl' against ngrok API and extract public (using 'sed' command)
  NGROK_PUBLIC_URL=$(curl --silent --max-time 10 --connect-timeout 5 \
                    --show-error http://127.0.0.1:4040/api/tunnels | \
                    sed -nE 's/.*public_url":"https:..([^"]*).*/\1/p')
  sleep 1
  echo -n "."
done

export NGROK_PUBLIC_URL="https://$NGROK_PUBLIC_URL"
export CLOUD_RUN_SERVICE_URL=$NGROK_PUBLIC_URL

echo
echo "NGROK_PUBLIC_URL => [ $NGROK_PUBLIC_URL ]"

