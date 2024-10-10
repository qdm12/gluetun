# Development container

Development container that can be used with VSCode.

It works on Linux, Windows (WSL2) and OSX.

## Requirements

- [VS code](https://code.visualstudio.com/download) installed
- [VS code dev containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) installed
- [Docker](https://www.docker.com/products/docker-desktop) installed and running

## Setup

1. Create the following files on your host if you don't have them:

    ```sh
    touch ~/.gitconfig ~/.zsh_history
    ```

1. **For OSX hosts**: ensure your home directory `~` is accessible by Docker.
1. Open the command palette in Visual Studio Code (CTRL+SHIFT+P).
1. Select `Dev-Containers: Open Folder in Container...` and choose the project directory.

## Customization

For customizations to take effect, you should "rebuild and reopen":

1. Open the command palette in Visual Studio Code (CTRL+SHIFT+P)
2. Select `Dev-Containers: Rebuild Container`

Customizations available are notably:

- Changes to the Docker image in [Dockerfile](Dockerfile)
- Changes to VSCode **settings** and **extensions** in [devcontainer.json](devcontainer.json).
- Change the entrypoint script by adding in [devcontainer.json](devcontainer.json) a bind mount to a shell script to `/root/.welcome.sh` to replace the [current welcome script](https://github.com/qdm12/godevcontainer/blob/master/shell/.welcome.sh). For example:

    ```json
    {
        "source": "/yourpath/.welcome.sh",
        "target": "/root/.welcome.sh",
        "type": "bind"
    },
    ```

- More customizations available are documented in the [devcontainer.json reference](https://containers.dev/implementors/json_reference/).
