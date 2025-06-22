# Autobrr Go Client Library

A Go client library for interacting with the [Autobrr](https://github.com/autobrr/autobrr) API.

## Features

- **Filter Management**: Create, read, update, and delete filters
- **Filter Control**: Enable/disable filters
- **Connection Testing**: Verify API connectivity
- **Full Filter Support**: All filter options including actions, external filters, and advanced criteria

## Installation

To install the package, run:

```bash
go get github.com/cehbz/autobrr
```

## Usage

### Importing the Package

```go
import (
    "github.com/cehbz/autobrr"
)
```

### Initializing the Client

```go
client, err := autobrr.NewClient("your-api-key", "localhost", "10798")
if err != nil {
    log.Fatalf("Failed to create client: %v", err)
}
```

- `apiKey`: Your Autobrr API key
- `addr`: The address where Autobrr is running (e.g., `"127.0.0.1"`)
- `port`: The port number of the Autobrr API (e.g., `"10798"`)

### Getting All Filters

```go
filters, err := client.GetFilters()
if err != nil {
    log.Fatalf("Failed to get filters: %v", err)
}

for _, filter := range filters {
    fmt.Printf("Filter: %s (ID: %d, Enabled: %v)\n", filter.Name, filter.ID, filter.Enabled)
}
```

### Getting a Specific Filter

```go
filter, err := client.GetFilter(123)
if err != nil {
    log.Fatalf("Failed to get filter: %v", err)
}

fmt.Printf("Filter: %s\n", filter.Name)
fmt.Printf("Shows: %v\n", filter.Shows)
fmt.Printf("Resolutions: %v\n", filter.Resolutions)
```

### Creating a Filter

```go
newFilter := &autobrr.Filter{
    Name:                  "New TV Show Filter",
    Enabled:               true,
    Priority:              100,
    SmartEpisode:          true,
    Shows:                 []string{"The Matrix"},
    Resolutions:           []string{"1080p", "2160p"},
    Sources:               []string{"WEB-DL", "WEBRip"},
    MatchReleaseGroups:    "NTb|FLUX",
    Years:                 "2023-2024",
    Actions: []autobrr.Action{
        {
            Name:      "qBittorrent",
            Type:      "QBITTORRENT",
            Enabled:   true,
            ClientID:  1,
            Category:  "tv-shows",
            SavePath:  "/downloads/tv",
        },
    },
}

createdFilter, err := client.CreateFilter(newFilter)
if err != nil {
    log.Fatalf("Failed to create filter: %v", err)
}

fmt.Printf("Created filter with ID: %d\n", createdFilter.ID)
```

### Updating a Filter

```go
filter.Resolutions = append(filter.Resolutions, "720p")
filter.Priority = 150

updatedFilter, err := client.UpdateFilter(filter.ID, filter)
if err != nil {
    log.Fatalf("Failed to update filter: %v", err)
}
```

### Deleting a Filter

```go
err := client.DeleteFilter(123)
if err != nil {
    log.Fatalf("Failed to delete filter: %v", err)
}
```

### Enabling/Disabling a Filter

```go
// Enable a filter
err := client.ToggleFilterEnabled(123, true)
if err != nil {
    log.Fatalf("Failed to enable filter: %v", err)
}

// Disable a filter
err = client.ToggleFilterEnabled(123, false)
if err != nil {
    log.Fatalf("Failed to disable filter: %v", err)
}
```

### Testing Connection

```go
err := client.TestConnection()
if err != nil {
    log.Fatalf("Connection test failed: %v", err)
}
fmt.Println("Successfully connected to Autobrr")
```

## Filter Options

The `Filter` struct supports all Autobrr filter options:

- **Basic**: Name, enabled status, priority
- **Content Matching**: Shows, movies, resolutions, sources, codecs, containers
- **Advanced Matching**: Release groups, uploaders, categories, languages
- **Size Limits**: Min/max size constraints
- **Special Options**: Freeleech, scene releases, smart episode detection
- **Actions**: Download client actions, webhooks, custom commands

## Action Types

Supported action types include:
- `QBITTORRENT`: qBittorrent download client
- `DELUGE`: Deluge download client
- `TRANSMISSION`: Transmission download client
- `RADARR`: Radarr integration
- `SONARR`: Sonarr integration
- `LIDARR`: Lidarr integration
- `WHISPARR`: Whisparr integration
- `WEBHOOK`: Custom webhooks
- `EXEC`: Execute custom commands
- `WATCH_FOLDER`: Watch folder for .torrent files

## Error Handling

The client returns detailed errors for various failure scenarios:

```go
filter, err := client.GetFilter(999)
if err != nil {
    // Handle specific error cases
    if strings.Contains(err.Error(), "404") {
        log.Fatal("Filter not found")
    } else if strings.Contains(err.Error(), "401") {
        log.Fatal("Invalid API key")
    } else {
        log.Fatalf("Error: %v", err)
    }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contribution

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgments

- [Autobrr API Documentation](https://autobrr.com/api)
