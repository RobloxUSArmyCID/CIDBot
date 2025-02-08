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
  }: let
    name = "cidbot";
    src = ./.;

    systemIndependent = {
      homeManagerModules.default = {
        config,
        lib,
        pkgs,
        ...
      }: let
        cfg = config.services.${name};
      in {
        options.services.${name} = {
          enable = lib.mkEnableOption "CIDBot service";
          package = lib.mkOption {
            type = lib.types.package;
            default = self.packages.${pkgs.system}.default;
            description = "The package to use for CIDBot.";
          };
        };

        config = lib.mkIf cfg.enable {
          home.packages = [cfg.package];
          systemd.user.services.cidbot = {
            Unit = {
              Description = "The CID Bot";
              After = ["network.target"];
            };

            Service = {
              ExecStart = "${cfg.package}/bin/cidbot";
              Restart = "always";
            };

            Install = {
              WantedBy = ["default.target"];
            };
          };
        };
      };
    };

    systemDependent = flake-utils.lib.eachDefaultSystem (system: let
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
  in
    systemDependent // systemIndependent;
}
