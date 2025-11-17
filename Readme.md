```yaml
docker run -d \
  --name newConn1 \
  -p 8080:80 \
  --label io.watcher.enable=true \
  --label io.watcher.slack.enable=true \
  --label io.watcher.slack.webhook="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  --label io.watcher.slack.channel="#amirath-lube" \
  aasourav/conn:latest

```

```yaml
services:
  newConn1:
    image: aasourav/conn:latest
    container_name: newConn1

    ports:
      - "8080:80"

    labels:
      io.watcher.enable: "true"
      io.watcher.slack.enable: "true"
      io.watcher.slack.webhook: "xadxxxxxxxxxxxxxxxxxxxx"
      io.watcher.slack.channel: "#amirath-lube"

    restart: unless-stopped

  watcher:
    image: aasourav/containerwatcher:latest
    container_name: watcher
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      GLOBAL_SLACK_WEBHOOK: "slack/xxxxxxx"
      GLOBAL_SLACK_CHANNEL: "#amirath-lube"
      WATCH_INTERVAL: "30"
      CLEAN_OLD_IMAGE: "true"
    restart: unless-stopped
```