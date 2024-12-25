let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-unstable";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in
pkgs.mkShellNoCC {
  packages = with pkgs; [
    cobra-cli
    cowsay
    go
    goreleaser
    mr
    pre-commit
  ];

  DIRENV_LOG_FORMAT = "";
  ENVIRONMENT = "mr, pre—commit, Go, Cobra-CLI, GoReleaser";

  shellHook = ''
    echo "$ENVIRONMENT" | cowsay -W 80
    echo -n "1. "; pre-commit install
  '';
}
