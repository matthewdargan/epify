{
  inputs = {
    nix-go = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:matthewdargan/nix-go";
    };
    nixpkgs.url = "nixpkgs/nixos-unstable";
    parts.url = "github:hercules-ci/flake-parts";
    pre-commit-hooks = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:cachix/pre-commit-hooks.nix";
    };
  };
  outputs = inputs:
    inputs.parts.lib.mkFlake {inherit inputs;} {
      imports = [inputs.pre-commit-hooks.flakeModule];
      perSystem = {
        config,
        inputs',
        lib,
        pkgs,
        ...
      }: {
        devShells.default = pkgs.mkShell {
          packages = [inputs'.nix-go.packages.go inputs'.nix-go.packages.golangci-lint];
          shellHook = "${config.pre-commit.installationScript}";
        };
        packages = {
          epify = inputs'.nix-go.legacyPackages.buildGoModule {
            meta = with lib; {
              description = "Categorize shows and movies using the Jellyfin naming scheme";
              homepage = "https://github.com/matthewdargan/epify";
              license = licenses.bsd3;
              maintainers = with maintainers; [matthewdargan];
            };
            pname = "epify";
            src = ./.;
            vendorHash = null;
            version = "0.1.0";
          };
          trdone = inputs'.nix-go.legacyPackages.buildGoModule {
            meta = with lib; {
              description = "Organize completed torrent downloads";
              homepage = "https://github.com/matthewdargan/epify";
              license = licenses.bsd3;
              maintainers = with maintainers; [matthewdargan];
            };
            pname = "trdone";
            src = ./.;
            vendorHash = null;
            version = "0.1.0";
          };
        };
        pre-commit = {
          settings = {
            hooks = {
              alejandra.enable = true;
              deadnix.enable = true;
              golangci-lint = {
                enable = true;
                package = inputs'.nix-go.packages.golangci-lint;
              };
              gotest = {
                enable = true;
                package = inputs'.nix-go.packages.go;
              };
              statix.enable = true;
            };
            src = ./.;
          };
        };
      };
      systems = ["aarch64-darwin" "aarch64-linux" "x86_64-darwin" "x86_64-linux"];
    };
}
