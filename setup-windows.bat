mkdir "C:\Program Files\docsel"
copy build\bin\docsel-windows_amd64.exe "C:\Program Files\docsel\docsel.exe"
setx PATH "C:\Program Files\docsel;%PATH%"
