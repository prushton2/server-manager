FROM nixos/nix

EXPOSE 3000

WORKDIR /app
COPY . /app

RUN mkdir -p ~/.config/nix && \
    echo "experimental-features = nix-command flakes" >> ~/.config/nix/nix.conf

RUN nix-env -iA nixpkgs.curl
RUN nix-env -iA nixpkgs.docker
RUN nix build .
CMD ["/app/result/bin/server-manager"]