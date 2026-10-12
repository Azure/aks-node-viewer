{
  description = "aks-node-viewer: a tool for visualizing dynamic node usage within a Kubernetes cluster";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f (import nixpkgs { inherit system; }));
      version = "0.1.0";
    in
    {
      packages = forAllSystems (
        pkgs:
        let
          aks-node-viewer = pkgs.buildGoModule {
            pname = "aks-node-viewer";
            inherit version;
            src = self;
            # go.mod requests a toolchain >= go1.24.1; the nixpkgs go satisfies that
            # without downloading a toolchain inside the sandbox.
            go = pkgs.go;
            vendorHash = "sha256-x+xTXQ3cSTOOgKzXXFEanjVu8qaZbJ5kR+sm4KAJ+VM=";
            subPackages = [ "cmd/aks-node-viewer" ];
            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
              "-X main.builtBy=nix"
            ];
            meta = {
              description = "Visualize dynamic node usage within a Kubernetes cluster";
              homepage = "https://github.com/Azure/aks-node-viewer";
              license = pkgs.lib.licenses.asl20;
              mainProgram = "aks-node-viewer";
            };
          };
        in
        {
          default = aks-node-viewer;
          aks-node-viewer = aks-node-viewer;
        }
      );

      apps = forAllSystems (pkgs: {
        default = {
          type = "app";
          program = "${self.packages.${pkgs.stdenv.hostPlatform.system}.default}/bin/aks-node-viewer";
        };
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.git
            pkgs.goreleaser
            pkgs.golangci-lint
            pkgs.go-licenses
          ];
        };
      });
    };
}
