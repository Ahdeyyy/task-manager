# Define the directory to watch and the command to run
$watchDir = "C:\Users\ahdey\task-manager"
$command = "sh build.sh; go run ."

# Create a FileSystemWatcher to monitor the directory
$watcher = New-Object System.IO.FileSystemWatcher
$watcher.Path = $watchDir
$watcher.Filter = "*.*"
$watcher.IncludeSubdirectories = $true
$watcher.NotifyFilter = [System.IO.NotifyFilters]'FileName, LastWrite'

# Define the event handler function
$action = {
    $changedFile = $Event.SourceEventArgs.FullPath
    if ($changedFile -match "\.templ$" -or $changedFile -match "\.go$") {
        Write-Host "Detected change in $changedFile, running command..."
        Invoke-Expression $command
    }
}

# Register event handlers
Register-ObjectEvent $watcher 'Changed' -Action $action
Register-ObjectEvent $watcher 'Created' -Action $action
Register-ObjectEvent $watcher 'Deleted' -Action $action
Register-ObjectEvent $watcher 'Renamed' -Action $action

# Start watching
$watcher.EnableRaisingEvents = $true

# Keep the script running
Write-Host "Watching directory $watchDir for changes. Press [Enter] to exit."
[void][System.Console]::ReadLine()

# Unregister events when done
Unregister-Event -SourceIdentifier $watcher.Changed
Unregister-Event -SourceIdentifier $watcher.Created
Unregister-Event -SourceIdentifier $watcher.Deleted
Unregister-Event -SourceIdentifier $watcher.Renamed
$watcher.Dispose()
