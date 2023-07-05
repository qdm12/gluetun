# Development container

Development container that can be used with VSCode.

It works on Linux, Windows and OSX.

## Requirements

- [VS code](https://code.visualstudio.com/download) installed
- [VS code remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed
- [Docker](https://www.docker.com/products/docker-desktop) installed and running
- [Docker Compose](https://docs.docker.com/compose/install/) installed

## Setup

1. Create the following files on your host if you don't have them:

    ```sh
    touch ~/.gitconfig ~/.zsh_history
    ```

    Note that the development container will create the empty directories `~/.docker`, `~/.ssh` and `~/.kube` if you don't have them.

1. **For Docker on OSX or Windows without WSL**: ensure your home directory `~` is accessible by Docker.
1. Open the command palette in Visual Studio Code (CTRL+SHIFT+P).
1. Select `Remote-Containers: Open Folder in Container...` and choose the project directory.

## Customization

### Customize the image

You can make changes to the [Dockerfile](Dockerfile) and then rebuild the image. For example, your Dockerfile could be:

```Dockerfile
FROM qmcgaw/godevcontainer
RUN apk add curl
```

To rebuild the image, either:

- With VSCode through the command palette, select `Remote-Containers: Rebuild and reopen in container`
- With a terminal, go to this directory and `docker-compose build`

### Customize VS code settings

You can customize **settings** and **extensions** in the [devcontainer.json](devcontainer.json) definition file.

### Entrypoint script

You can bind mount a shell script to `/root/.welcome.sh` to replace the [current welcome script](https://github.com/qdm12/godevcontainer/blob/master/shell/.welcome.sh).

### Publish a port

To access a port from your host to your development container, publish a port in [docker-compose.yml](docker-compose.yml). You can also now do it directly with VSCode without restarting the container.

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
