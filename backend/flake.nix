{
  description = "Server Manager Backend";

  inputs = {
    nixpkgs.url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/25.05.tar.gz";
  };

  outputs = { self, nixpkgs }:
  let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
    modname = "servermanagerbackend";
  in
  {
    packages.x86_64-linux.default = pkgs.buildGoModule {
      pname = modname;
      version = "0.1.0";
      src = ./.;
      vendorHash = "sha256-g+yaVIx4jxpAQ/+WrGKxhVeliYx7nLQe/zsGpxV4Fn4=";
      doCheck = false;
    };

    devShells.x86_64-linux.default = pkgs.mkShell {
      name = modname;
      shellHook = ''export PS1="\[\e[1;32m\][nix:${modname}]\]$ \e[0m"'';
      packages = with pkgs; [
        docker
        go
        gcc
      ];
    };
  };
}