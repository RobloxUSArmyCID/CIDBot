# CIDBot

The United States Army Criminal Investigation Division's background checking bot.

#### This bot is now hosted 24/7, and will be as long as I remain in CID.

# Download instructions
1. Go to https://github.com/RobloxUSArmyCID/CIDBot/releases/latest
2. Download the file appropriate for your Operating System and architecture

# Running instructions
1. On Windows (Powershell):
```powershell
cd <PATH_TO_CIDBOT_FOLDER>
.\CIDBot-vX.X.X-windows-arch.exe --token <TOKEN>
```
OR
```powershell
cd <PATH_TO_CIDBOT_FOLDER>
.\CIDBot-vX.X.X-windows-arch.exe --token-path <PATH_TO_FILE_CONTAINING_TOKEN>
```
> You may have to wait a little on the first run for your antivirus to scan the file.
> You may also have to approve the file through Windows SmartScreen by pressing more details, then "Run anyway".
2. On other platforms:
```bash
cd <PATH_TO_CIDBOT_FOLDER>
chmod +x CIDBot-vX.X.X-os-arch
./CIDBot-vX.X.X-os-arch --token <TOKEN>
```
OR
```bash
cd <PATH_TO_CIDBOT_FOLDER>
chmod +x CIDBot-vX.X.X-os-arch
./CIDBot-vX.X.X-os-arch --token-path <PATH_TO_FILE_CONTAINING_TOKEN>
```
> On macOS, you may have to open privacy & security settings and approve the file there.

# Closing the bot - instructions
- Press Ctrl/Cmd+C
- Or close the terminal/PowerShell window.

# FAQ
* Is it any good?
  * [Yes.](https://news.ycombinator.com/item?id=3067434)

# Changing the token - instructions for CID HICOM
1. Go to https://discord.com/developers/applications
2. Select the CID Bot.
3. Click on Bot.
4. Click on Reset Token.
5. Put in your Discord 2FA.

# Building instructions
The **only** supported building environment is Linux. It *should* work on macOS, but I don't guarantee anything. If you're on Windows, use WSL.
- Testing:
```bash
cd <PATH_TO_SOURCE>
make build # make clean build to clean out the bin/ folder beforehand
```
- Release:
```bash
cd <PATH_TO_SOURCE>
make CIDBOT_VERSION=<VERSION> clean release # do not use v for the version (ex. v2.0.0)
```

*Made with :heart:,
in Poland,
by f_o1oo.*
