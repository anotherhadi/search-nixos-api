{
  description = "Search NixOS Api";

  inputs = { nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable"; };

  outputs = { self, nixpkgs }:
    let pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
      packages.x86_64-linux.search-nixos-api = pkgs.buildGoModule {
        pname = "search-nixos-api";
        version = "0.1.0";
        src = ./.;
        vendorHash =
          "sha256-EcaOsOWqr4j5bTZBMxwsOI9l3Gb+aZOPDul9Whrn1mg="; # `nix build .#search-nixos-api`
      };

      nixosModules.search-nixos-api = { config, lib, pkgs, ... }: {
        options.services.search-nixos-api = {
          enable = lib.mkEnableOption "Enable the search-nixos-api service";
          port = lib.mkOption {
            type = lib.types.int;
            default = 8090;
            description = "Port for the search-nixos-api service";
          };
          indexPath = lib.mkOption {
            type = lib.types.path;
            default = "/var/lib/search-nixos-api/index.json";
            description =
              "Path to the index file for the search-nixos-api service. Make sure to set the correct permissions.";
          };
          user = lib.mkOption {
            type = lib.types.str;
            default = "search-nixos-api";
            description = "User for the search-nixos-api service";
          };
          group = lib.mkOption {
            type = lib.types.str;
            default = "search-nixos-api";
            description = "Group for the search-nixos-api service";
          };
          interval = lib.mkOption {
            type = lib.types.str;
            default = "12h";
            description = "Interval for the search-nixos-api service";
          };
        };

        config = lib.mkIf config.services.search-nixos-api.enable {
          systemd.services.search-nixos-api = {
            description = "Search NixOS API";
            after = [ "network.target" ];
            wantedBy = [ "multi-user.target" ];
            serviceConfig = {
              ExecStart =
                "${self.packages.x86_64-linux.search-nixos-api}/bin/cmd";
              Restart = "always";
              User = config.services.search-nixos-api.user;
              Group = config.services.search-nixos-api.group;
              DynamicUser = true;
              StateDirectory = "search-nixos-api";
              ReadWritePaths = [ "/var/lib/search-nixos-api" ];
              Environment = [
                "PRODUCTION=true"
                "PORT=${toString config.services.search-nixos-api.port}"
                "INTERVAL=${config.services.search-nixos-api.interval}"
                "INDEX_PATH=${config.services.search-nixos-api.indexPath}"
              ];
            };
          };
        };
      };
    };
}
