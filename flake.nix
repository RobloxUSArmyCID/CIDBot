{
  description = "A flake for the CID bot";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    yaml = {
      url = "github:folospior/yaml.arm64.nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
    yaml,
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
        extraCfg =
          if cfg.extraConfig == {}
          then cfg.configFile
          else yaml.lib.${pkgs.system}.toYaml cfg.extraConfig;
      in {
        options.services.${name} = {
          enable = lib.mkEnableOption "CIDBot service";
          package = lib.mkOption {
            type = lib.types.package;
            default = self.packages.${pkgs.system}.default;
            description = "The package to use for CIDBot.";
          };
          extraConfig = lib.mkOption {
            type = lib.types.attrs;
            default = {};
            description = "Configuration to put in the config.yml file. Takes prescedence over configFile.";
          };
          configFile = lib.mkOption {
            type = lib.types.path;
            default = "";
            description = "Configuration to put in the config.yml file in the form of YAML file.";
          };
          configPath = lib.mkOption {
            default = "${config.xdg.configHome}/CIDBot/config.yml";
            type = lib.types.str;
            description = "Where to put the config.yml file";
          };
        };

        config = lib.mkIf cfg.enable {
          home.packages = [cfg.package];
          home.file.${cfg.configPath}.source = extraCfg;
          systemd.user.services.cidbot = {
            Unit = {
              Description = "The CID Bot";
              After = ["network.target"];
            };

            Service = {
              ExecStart = "${cfg.package}/bin/cidbot --config-path ${cfg.configPath}";
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

        vendorHash = "sha256-jQ/OoLciPd4KvRf4AErxENPQvjNNl/4tq9ZxLohVtas=";

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
