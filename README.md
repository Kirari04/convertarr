# Convertarr

![Alt text](img/image.png)

## Install using Docker

### Software encoding

```bash
docker run -d \
    --name convertarr \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Etc/UTC \
    -v /path/to/data:/app/database \
    -v /path/to/videofiles:/videofiles `#optional` \
    -p 8080:8080 \
    kirari04/convertarr:latest
```

or using docker-compose

```yaml
---
services:
  convertarr:
    image: kirari04/convertarr:latest
    container_name: convertarr
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
    volumes:
      - /path/to/data:/app/database
      - /path/to/videofiles:/videofiles #optional
    ports:
      - 8080:8080
    restart: unless-stopped
```

### Hardware encoding

#### Nvidia

```bash
docker run -d \
    --name convertarr \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Etc/UTC \
    -v /path/to/data:/app/database \
    -v /path/to/videofiles:/videofiles `#optional` \
    -p 8080:8080 \
    --gpus all \
    kirari04/convertarr:latest-nvidia
```

or using docker-compose

```yaml
---
services:
  convertarr:
    image: kirari04/convertarr:latest-nvidia
    container_name: convertarr
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
    volumes:
      - /path/to/data:/app/database
      - /path/to/videofiles:/videofiles #optional
    ports:
      - 8080:8080
    deploy:
      resources:
        reservations:
          devices:
            - driver: nvidia
              count: all
              capabilities: [gpu,compute,video]
    restart: unless-stopped
```

#### VAAPI

```bash
docker run -d \
    --name convertarr \
    -e PUID=1000 \
    -e PGID=1000 \
    -e TZ=Etc/UTC \
    -v /path/to/data:/app/database \
    -v /path/to/videofiles:/videofiles `#optional` \
    -p 8080:8080 \
    --device /dev/dri:/dev/dri \
    kirari04/convertarr:latest-vaapi
```

or using docker-compose

```yaml
---
services:
  convertarr:
    image: kirari04/convertarr:latest-vaapi
    container_name: convertarr
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Etc/UTC
    volumes:
      - /path/to/data:/app/database
      - /path/to/videofiles:/videofiles #optional
    ports:
      - 8080:8080
    devices:
      - /dev/dri:/dev/dri
    restart: unless-stopped
```

## Developement

### Server Application

```bash
go run main.go serve
```

### Watch Templ

```bash
templ generate --watch
```
### Build

```bash
docker build --platform linux/amd64 -t kirari04/convertarr:latest --push .
```

### Build On System

```bash
go build -o main.bin main.go
docker build --platform linux/amd64 -t kirari04/convertarr:latest --push -f Dockerfile.main
```
