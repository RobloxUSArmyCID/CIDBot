# CIDBot

The United States Army Criminal Inverstigation Division's background checking bot.

###### Current token: MTI2MzgzMTYyMDUxOTk4OTMxMA.GsevhD.tM47MfbTbbf_tseeIRyh58z9t4XMUyuVCXUV84

# Download Instructions
1. On this page, click on actions.
2. Select the top build.
3. Download the artifact for your operating system and architecture.

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

# Changing the token - instructions for CID HICOM
1. Click on the README.md file in this repository.
2. Change the token in the appropriate line.
3. Commit changes, but add `[skip ci]` to the title.
> This is so there won't be a new build of the bot for changing the token.


*Made with :heart:,
in Poland,
by f_o1oo.*
