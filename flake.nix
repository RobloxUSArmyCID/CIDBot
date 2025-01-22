{
  description = "A flake for the CID bot";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      name = "cidbot";
      src = ./.;
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.default = pkgs.buildGoModule {
        inherit src name;

        buildPhase = ''
          runHook preBuild
          mkdir -p $TMPDIR/bin
          go build -o $TMPDIR/bin/${name} ./main/main.go
          runHook postBuild
        '';

        vendorHash = "sha256-mW9uHq9CrWYVg1STmxnJ3uRNiOK5Wxv9hJvGWNkplbU=";

        installPhase = ''
          runHook preInstall
          install -m755 -D $TMPDIR/bin/${name} $out/bin/${name}
          runHook postInstall
        '';
      };

      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          go-tools
          gopls
        ];
      };
    });
}
