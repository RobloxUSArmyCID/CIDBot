# CIDBot

The United States Army Criminal Investigation Division's background checking bot.

#### This bot is now hosted 24/7, and will be as long as I remain in CID.

# Usage instructions

The bot exposes a package and a Home Manager Module.
You can use the Home Manager Module in your flake:

```nix
{
    inputs = {
        nixpkgs.url = "nixpkgs/nixos-unstable";
        home-manager = {
            url = "github:nix-community/home-manager";
            inputs.nixpkgs.follows = "nixpkgs";
        };
        cidbot = {
            url = "git+ssh://git@github.com/RobloxUSArmyCID/CIDBot"
            inputs.nixpkgs.follows = "nixpkgs";
        };
    };

    outputs = { nixpkgs,  cidbot, ...}: let
        username = "user";
        system = "x86_64-linux" # confirmed - it works on aarch64-linux, not sure about darwin, it _should_ tho
        pkgs = import nixpkgs { inherit system; };
    in
    {
        homeConfigurations.${user} = home-manager.lib.homeManagerConfiguration {
            inherit pkgs;
            modules = [
                cidbot.homeManagerModules.default
                # your modules
            ];
        };
    };
}
```

An example configuration of a module would be:

```nix
{...}: {
    services.cidbot = {
        enable = true;
        extraConfig = {
            token = "INSERT TOKEN HERE";
            is_development = true;
            admin_server_id = "ADMIN SERVER ID";
            whitelist_path = "PATH/TO/WHITELIST";
        };
    };
}
```

# FAQ

- Is it any good?
  - [Yes.](https://news.ycombinator.com/item?id=3067434)

# Changing the token - instructions for CID HICOM

1. Go to https://discord.com/developers/applications
2. Select the CID Bot.
3. Click on Bot.
4. Click on Reset Token.
5. Put in your Discord 2FA.
6. Change the token in the configuration.

# Building instructions

The **only** supported building environment is Linux. It _should_ work on macOS, but I don't guarantee anything. If you're on Windows, use WSL.

```bash
cd <PATH_TO_SOURCE>
nix build
```

_Made with :heart:,
in Poland,
by f_o1oo._
