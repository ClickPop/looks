$LOOKS_PATH = "C:\looks\"
$LOOKS_FILE = $LOOKS_PATH + "looks.exe"
$ARCH = $env:PROCESSOR_ARCHITECTURE.ToLower()
$IN_PATH = $env:PATH.Contains($LOOKS_PATH)
$URI = "https://github.com/clickpop/looks/releases/latest/download/looks-windows-" + $ARCH + ".exe"

if ( $ARCH -ne "amd64" -and $ARCH -ne "arm64" ) {
  echo "unsupported architecture"
  Exit 0
}

if ( -not ( Test-Path -Path $LOOKS_PATH ) ) {
  New-Item -Path "C:\" -Name "looks" -ItemType "directory"
}

if ( -not $IN_PATH ) {
  [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";" + $LOOKS_PATH, [EnvironmentVariableTarget]::Machine)
}

Invoke-RestMethod -Uri $URI -OutFile $LOOKS_FILE