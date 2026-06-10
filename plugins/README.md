# plugins/ (deploy only)

Native binaries are **not built in this folder**.

Build and copy from the server repo:

```powershell
cd D:\vcmp-go-server
.\build.ps1
```

Expected file: `goserver04rel64.dll` (Windows) or `goserver04rel64.so` (Linux).

Do not use `goplugin04rel64.dll` from the blank example — that is a different plugin with no Safari commands.
