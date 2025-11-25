# WatchContainer

<p align="center">
  <img src="https://img.shields.io/badge/status-active-brightgreen" />
  <img src="https://img.shields.io/badge/docker-ready-blue" />
  <img src="https://img.shields.io/badge/license-MIT-purple" />
  <img src="https://img.shields.io/badge/slack-notifications-orange" />
</p>

**WatchContainer** is an open‑source, lightweight container auto‑updater that monitors your running Docker containers and automatically updates them whenever a new version of their image is published on Docker Hub. It also supports optional Slack notifications for update events.

This project is similar to Watchtower but more focused, simpler, and supports per‑container configuration through Docker labels.

---

## 🚀 Features

* 🔍 Watches running containers for new image versions
* ⬆️ Automatically pulls the latest image
* 🔄 Recreates the container with the new image (keeping all existing config)
* 🧹 Optional cleanup of old images
* 🔔 Slack notifications (global or per‑container)
* 🏷️ Per‑container configuration using labels
* ⏱️ Configurable polling interval

---

## 📦 Installation

Clone this repo or pull the image directly:

```bash
docker pull aasourav/watchcontainer:latest
```

---

## 🛠️ Usage

### Docker Compose Example

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
    image: aasourav/watchcontainer:latest
    container_name: watcher
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    # - /var/lib/docker:/var/lib/docker:ro # optional but useful for metadata
    evironment:
      GLOBAL_SLACK_WEBHOOK: "slack/xxxxxxx"
      GLOBAL_SLACK_CHANNEL: "#amirath-lube"
      WATCH_INTERVAL: "30"
      CLEAN_OLD_IMAGE: "true"
    restart: unless-stopped
```

Run the watcher using Docker Compose:

```yaml
services:
  watcher:
    image: aasourav/watchcontainer:latest
    container_name: watcher

    volumes:
      - /var/run/docker.sock:/var/run/docker.sock

    environment:
      GLOBAL_SLACK_WEBHOOK: ""
      GLOBAL_SLACK_CHANNEL: ""
      WATCH_INTERVAL: "30"      # default: 10 sec
      CLEAN_OLD_IMAGE: "true"   # optional

    restart: unless-stopped
```

The watcher will now automatically scan for containers that include `io.watcher.enable=true`.

---

## ⚠️ Important Notes

### 1. Use only the `latest` tag

WatchContainer currently checks for updates based on the `latest` tag. Make sure the containers you want to monitor are running images tagged as `latest`.

### 2. **Do NOT build images on the same server where WatchContainer is running**

If you build a new image locally on the same server while WatchContainer is active, the following can happen:

* When you rebuild an image with the **same tag (`latest`)**, Docker temporarily un-tags the existing image.
* During this untagged window, WatchContainer may detect:

  * The **untagged old image** as if it's the *new* one.
  * And the **new freshly built image** as *untagged*.
* This causes WatchContainer to pull incorrect layers and trigger an unnecessary update.

This scenario **only happens when building images locally on the same host** where WatchContainer is running.

➡️ In real production deployments, images should be built in CI/CD or a separate build server — **not on the live Docker host**.

---

## 🧩 Container Labels

Attach these labels to **any container you want WatchContainer to manage**.

| Label                      | Type    | Description                        | Default |
| -------------------------- | ------- | ---------------------------------- | ------- |
| `io.watcher.enable`        | boolean | Enables watcher for this container | `false` |
| `io.watcher.slack.enable`  | boolean | Enables Slack notifications        | `false` |
| `io.watcher.slack.webhook` | string  | Slack webhook (container-specific) | `""`    |
| `io.watcher.slack.channel` | string  | Slack channel (container-specific) | `""`    |

---

## ▶️ Example: Running a Watched Container

```bash
docker run -d \
  --name newConn1 \
  -p 8080:80 \
  --label io.watcher.enable=true \
  --label io.watcher.slack.enable=true \
  --label io.watcher.slack.webhook="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" \
  --label io.watcher.slack.channel="#amirath-lube" \
  aasourav/conn:latest
```

### Label Explanation

* `io.watcher.enable=true` → The watcher monitors this container
* `io.watcher.slack.enable=true` → Enable Slack notifications
* `io.watcher.slack.webhook="..."` → Slack webhook specific to this container
* `io.watcher.slack.channel="#channel"` → Slack channel for this container

If container-specific Slack settings are missing, the watcher falls back to global Slack settings.

---

## 🌍 Global Environment Variables (Watcher)

| Variable               | Description                        | Default |
| ---------------------- | ---------------------------------- | ------- |
| `GLOBAL_SLACK_WEBHOOK` | Fallback Slack webhook             | `""`    |
| `GLOBAL_SLACK_CHANNEL` | Fallback Slack channel             | `""`    |
| `WATCH_INTERVAL`       | Watch interval (seconds)           | `10`    |
| `CLEAN_OLD_IMAGE`      | Remove old images (`true`/`false`) | `false` |

---

## 🔔 Slack Notification Logic

Watcher sends Slack notifications when:

* A new image is detected
* The container is being upgraded
* The update succeeds or fails

Notification priority:

1. Container labels (`io.watcher.slack.*`)
2. Global variables (`GLOBAL_SLACK_*`)
3. No Slack notification if both are missing

---

## 🧹 Image Cleanup

If `CLEAN_OLD_IMAGE=true`, WatchContainer will automatically remove old image versions after redeploying the updated container.

---

## 🛣️ Roadmap

* Support for private registries (Harbor, ECR, GCR, GitHub Container Registry)
* Support for Telegram / Discord notifications
* Web dashboard to view updates
* Kubernetes version (DaemonSet)

---

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit PRs.

---

## 📄 License

MIT License.
