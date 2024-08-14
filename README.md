# CIDBot

The United States Army Criminal Inverstigation Division's background checking bot.

###### Current token: MTI2MzgzMTYyMDUxOTk4OTMxMA.GsevhD.tM47MfbTbbf_tseeIRyh58z9t4XMUyuVCXUV84

# Download instructions
1. On this page, click on actions.
2. On the left-hand side, select `.NET`.
3. Select the top build.
4. Download the artifact for your operating system and architecture.

# Running instructions
1. Download the zip file.
2. Extract the file.
3. On Windows:
- Open PowerShell
- Type (replacing the text with the actual token, found in the README):
```powershell
$Env:CIDBot_TOKEN="ENTER TOKEN HERE"
```
- Drag and drop the file into the PowerShell window and press enter.
> You may have to wait a little on the first run for your antivirus to scan the file.
> You may also have to approve the file through Windows SmartScreen by pressing more details, then "Run anyway".
3. On other platforms
- Open your terminal.
- Type in (replacing the text and drag-and-dropping where required):
```bash
export CIDBot_TOKEN="ENTER TOKEN HERE"
chmod +x (DRAG AND DROP THE FILE HERE)
```
- Drag and drop the file into the terminal window and press enter.
> On macOS, you may have to open privacy & security settings and approve the file there.

# Closing the bot - instructions
- Press Ctrl/Cmd+C
- Or close the terminal/PowerShell window.

# Changing the token - instructions for CID HICOM
1. Go to https://discord.com/developers/applications
2. Select the CID Bot.
3. Click on Bot.
4. Click on Reset Token.
5. Put in your Discord 2FA.
6. Click on the README.md file in this repository.
7. Change the token in the appropriate line.
8. Commit changes, but add `[skip ci]` to the title.
> This is so there won't be a new build of the bot for changing the token. (Picture for reference)
![image](https://github.com/user-attachments/assets/deb6417e-ec0d-4f83-ad31-8deeda7d7a5b)

*Made with :heart:,
in Poland,
by f_o1oo.*
