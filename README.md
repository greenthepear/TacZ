**TacZ** is a turn based strategy zombie game made in Ebitenengie, inspired by Into the Breach.

# Building

If you have [Go](https://go.dev/doc/install) and git:

    git clone https://github.com/greenthepear/TacZ
    cd TacZ
    go build

On Linux you'll need a gcc compiler and some dependencies for Ebitengine. [See more information here](https://ebitengine.org/en/documents/install.html).

- Debian / Ubuntu: `sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config`
- Fedora: `sudo dnf install mesa-libGL-devel mesa-libGLES-devel libXrandr-devel libXcursor-devel libXinerama-devel libXi-devel libXxf86vm-devel alsa-lib-devel pkg-config`
- Solus: `sudo eopkg install libglvnd-devel libx11-devel libxrandr-devel libxinerama-devel libxcursor-devel libxi-devel libxxf86vm-devel alsa-lib-devel pkg-config`
- Arch: `sudo pacman -S mesa libxrandr libxcursor libxinerama libxi pkg-config`
- Alpine: `sudo apk add alsa-lib-dev libx11-dev libxrandr-dev libxcursor-dev libxinerama-dev libxi-dev mesa-dev`

# Credits
Graphics nearly entirely consist of a [tileset made by Ittai Manero](https://ittaimanero.itch.io/zombie-apocalypse-tileset) who generously allows it to be used in any personal and commercial project. Thanks!
