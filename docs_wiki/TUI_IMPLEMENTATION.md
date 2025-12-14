# TUI Implementation Documentation

This document outlines the implementation of the Terminal User Interface (TUI) for the application, constructed using the Bubble Tea framework.

## Overview

The application features a sophisticated TUI that provides visual feedback during the boot sequence and a live dashboard for monitoring system resources and service status. The implementation utilizes the Model-View-Update (MVU) architecture provided by Bubble Tea.

## Key Technologies

- **Bubble Tea**: The primary framework for the TUI loop.
- **Lipgloss**: Used for styling, layouts, and color management.
- **Bubbles**: Provides pre-built components like spinners and progress bars.

## Components

The TUI is divided into two main components: the Boot Sequence and the Live Dashboard.

### 1. Boot Sequence (`pkg/tui/boot.go`)

The boot sequence handles the visualization of service initialization. It guides the user through the startup process with real-time feedback.

#### Model Structure

The `BootModel` struct manages the state of the boot process:
- **State Tracking**: Monitors the current phase (starting, initializing, complete, countdown, error).
- **Service Queue**: Maintains a list of services to be initialized and their current status (pending, loading, success, error, skipped).
- **Animation**: manages frame counters for visual effects like the wave animation.

#### Features

- **Phased Execution**: The logic transitions through distinct phases:
    1.  **Starting**: Brief intro animation.
    2.  **Initializing**: Iterates through the service queue, executing initialization functions.
    3.  **Complete/Countdown**: Displays a success message and an optional countdown before proceeding.
- **Visual Feedback**:
    -   Uses a spinner for active tasks.
    -   Displays a progress bar showing overall completion.
    -   Lists individual services with color-coded status indicators.
- **Interactivity**: The user can skip the countdown or quit the application using keyboard commands (q, esc, ctrl+c).

### 2. Live Dashboard (`pkg/tui/dashboard.go`)

Once the application is running, the dashboard provides a real-time view of system metrics and component health.

#### Model Structure

The `DashboardModel` struct holds the monitoring data:
- **System Stats**: CPU usage, Memory usage (used vs total), and Goroutine count.
- **Infrastructure Status**: Connection states for external dependencies (e.g., databases, caches).
- **Service Status**: Health status of internal services.

#### Features

- **Real-time Updates**: A background tick command triggers updates every 500ms.
- **System Monitoring**:
    -   Retrieves CPU and Memory statistics using `gopsutil`.
    -   Visualizes usage with color-coded progress bars (Green < 50%, Yellow < 80%, Red > 80%).
- **Layout**:
    -   **Header**: Displays app name, version, and running status.
    -   **System Box**: Shows resource usage details.
    -   **Infrastructure Box**: Lists external connections.
    -   **Services Box**: Lists internal service statuses.
- **Responsive Design**: Adapts to terminal window size changes.

## Styling System (`pkg/tui/styles.go`)

The application uses `lipgloss` to define a consistent design language. The color palette appears to be inspired by the Dracula theme.

- **Colors**:
    -   Primary/Accents: Pink (#FF79C6), Purple (#BD93F9), Cyan (#8BE9FD).
    -   Status: Green (#50FA7B) for success, Yellow (#F1FA8C) for warning, Red (#FF5555) for error.
    -   UI Elements: Dark Grey (#6272A4) for borders and muted text.
- **Typography**: Uses bold text for headers and distinct colors for labels vs values.
- **Animations**:
    -   **Wave**: A string-array based animation frame loop in the boot screen.
    -   **Pulse**: Color cycling on headers to indicate activity.

## Architecture & Data Flow

1.  **Initialization (`Init`)**:
    -   Starts the spinner tick.
    -   Starts the custom tick loop (boot tick or dashboard tick) to drive animations and updates.

2.  **Update Loop (`Update`)**:
    -   **KeyMsg**: Handles user input for quitting.
    -   **WindowSizeMsg**: Recalculates layout dimensions when the terminal is resized.
    -   **TickMsg**:
        -   Updates animation frames.
        -   Advanced boot logic (transitions phases, starts services).
        -   Refreshes system statistics (Dashboard only).

3.  **Rendering (`View`)**:
    -   Constructs the string representation of the UI using `lipgloss` styles.
    -   Uses `strings.Builder` for efficient string concatenation.
    -   Renders sub-components (boxes, progress bars) and joins them spatially.

## Usage

To use these components, the application entry point calls the respective Run functions:

```go
// Run Boot Sequence
results, err := tui.RunBootSequence(config, initQueue)

// Run Dashboard
err := tui.RunDashboardTUI(config, infraStatus, serviceStatus)
```

Both functions encapsulate the `tea.NewProgram` creation and execution, handling the alternative screen buffer automatically.
