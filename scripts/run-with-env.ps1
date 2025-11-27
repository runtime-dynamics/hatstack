# Wrapper script to load .env file and run the application
# Used by Air for live-reload with environment variables

param(
    [Parameter(Mandatory=$true)]
    [string]$ExecutablePath
)

# Function to load .env file
function Load-EnvFile {
    param (
        [string]$EnvFilePath
    )
    
    if (Test-Path $EnvFilePath) {
        Get-Content $EnvFilePath | ForEach-Object {
            $line = $_.Trim()
            # Skip empty lines and comments
            if ($line -and -not $line.StartsWith("#")) {
                $parts = $line -split "=", 2
                if ($parts.Count -eq 2) {
                    $key = $parts[0].Trim()
                    $value = $parts[1].Trim()
                    # Remove quotes if present
                    $value = $value -replace '^["'']|["'']$', ''
                    [Environment]::SetEnvironmentVariable($key, $value, "Process")
                }
            }
        }
    }
}

# Load .env file from project root
$envFilePath = Join-Path (Join-Path $PSScriptRoot "..") ".env"
Load-EnvFile -EnvFilePath $envFilePath

# Run the executable with all environment variables
& $ExecutablePath
