# NEPTUNE

A **Work In Progress** YouTube based song player, made with libmpv, GTK and GoLang.

_If you want to listen to song inside your terminal, take a look at 
[Tuner](https://github.com/Pauloo27/tuner)._

## Installing
### Arch Linux
You can install Neptune from the AUR, the package is called `neptune-git`.

### Other distros
First, install the required packages to build Neptune. They are:
- the GoLang compiler;
- mpv (with libmpv);
- GCC, libgtk 3 and libappindicator for gtk3 (required by the systray package);
- YouTube-DL (it's not required to compile, but it's required to run);

_if you want a MPRIS integration, install mpv-mpris_.

Then, clone the repository and run `make install`.

## License

<img src="https://i.imgur.com/AuQQfiB.png" alt="GPL Logo" height="100px" />

This project is licensed under [GNU General Public License v2.0](./LICENSE).

This program is free software; you can redistribute it and/or modify 
it under the terms of the GNU General Public License as published by 
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.
