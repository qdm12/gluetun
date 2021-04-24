# Development container

Development container that can be used with VSCode.

It works on Linux, Windows and OSX.

## Requirements

- [VS code](https://code.visualstudio.com/download) installed
- [VS code remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed
- [Docker](https://www.docker.com/products/docker-desktop) installed and running
    - If you don't use Linux or WSL 2, share your home directory `~/` and the directory of your project with Docker Desktop
- [Docker Compose](https://docs.docker.com/compose/install/) installed
- Ensure your host has the following and that they are accessible by Docker:
    - `~/.ssh` directory
    - `~/.gitconfig` file (can be empty)

## Setup

1. Open the command palette in Visual Studio Code (CTRL+SHIFT+P).
1. Select `Remote-Containers: Open Folder in Container...` and choose the project directory.
1. For Docker running on Windows HyperV, if you want to use SSH keys, bind mount them at `/tmp/.ssh` by changing the `volumes` section in the [docker-compose.yml](docker-compose.yml).

## Customization

### Customize the image

You can make changes to the [Dockerfile](Dockerfile) and then rebuild the image. For example, your Dockerfile could be:

```Dockerfile
FROM qmcgaw/godevcontainer
USER root
RUN apk add curl
USER vscode
```

Note that you may need to use `USER root` to build as root, and then change back to `USER vscode`.

To rebuild the image, either:

- With VSCode through the command palette, select `Remote-Containers: Rebuild and reopen in container`
- With a terminal, go to this directory and `docker-compose build`

### Customize VS code settings

You can customize **settings** and **extensions** in the [devcontainer.json](devcontainer.json) definition file.

### Entrypoint script

You can bind mount a shell script to `/home/vscode/.welcome.sh` to replace the [current welcome script](shell/.welcome.sh).

### Publish a port

To access a port from your host to your development container, publish a port in [docker-compose.yml](docker-compose.yml).

### Run other services

1. Modify [docker-compose.yml](docker-compose.yml) to launch other services at the same time as this development container, such as a test database:

    ```yml
      database:
        image: postgres
        restart: always
        environment:
          POSTGRES_PASSWORD: password
    ```

1. In [devcontainer.json](devcontainer.json), change the line `"runServices": ["vscode"],` to `"runServices": ["vscode", "database"],`.
1. In the VS code command palette, rebuild the container.
