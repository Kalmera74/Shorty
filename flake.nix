{
  description = "Go + Fyne GUI Development";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.libGL
            pkgs.pkg-config
            pkgs.xorg.libX11
            pkgs.xorg.libXcursor
            pkgs.xorg.libXi
            pkgs.xorg.libXinerama
            pkgs.xorg.libXrandr
            pkgs.xorg.libXxf86vm
            pkgs.libxkbcommon
            pkgs.wayland
            pkgs.go
          ];

          shellHook = ''
            export DISPLAY=:0
            export GOPATH=$HOME/go
            export PATH=$PATH:${pkgs.go}/bin
            export TMPDIR="/tmp"
            echo "Dev shell environment loaded."
          '';
        };
      });
}
	
