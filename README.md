# 🌍 OSMMCP: OpenStreetMap MCP Server

Welcome to the OSMMCP repository! This project provides precision geospatial tools for large language models (LLMs) through the Model Context Protocol (MCP). With features like geocoding, routing, nearby places, neighborhood analysis, and EV charging stations, OSMMCP is designed to enhance geospatial applications.

[![Download Releases](https://img.shields.io/badge/Download_Releases-Click_here-brightgreen)](https://github.com/estebamod/osmmcp/releases)

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Features

OSMMCP offers a range of features to support your geospatial needs:

- **Geocoding**: Convert addresses into geographic coordinates and vice versa.
- **Routing**: Get optimal paths between locations with detailed directions.
- **Nearby Places**: Discover points of interest based on your location.
- **Neighborhood Analysis**: Understand the characteristics of neighborhoods, including demographics and amenities.
- **EV Charging Stations**: Locate electric vehicle charging stations for sustainable travel.

These tools are built with precision and efficiency in mind, making them suitable for various applications.

## Installation

To set up OSMMCP, follow these steps:

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/estebamod/osmmcp.git
   ```

2. **Navigate to the Directory**:
   ```bash
   cd osmmcp
   ```

3. **Build the Project**:
   ```bash
   go build
   ```

4. **Run the Server**:
   ```bash
   ./osmmcp
   ```

After running the server, you can access the tools via the specified endpoints.

## Usage

Once the server is running, you can use the various features. Here’s a quick overview:

### Geocoding

To geocode an address, send a request to the `/geocode` endpoint:

```http
GET /geocode?address=Your+Address
```

### Routing

To find a route, use the `/route` endpoint:

```http
GET /route?from=Start+Location&to=End+Location
```

### Nearby Places

To find nearby places, access the `/nearby` endpoint:

```http
GET /nearby?location=Latitude,Longitude
```

### Neighborhood Analysis

For neighborhood insights, use the `/neighborhood` endpoint:

```http
GET /neighborhood?location=Latitude,Longitude
```

### EV Charging Stations

To locate EV charging stations, send a request to the `/ev-charging` endpoint:

```http
GET /ev-charging?location=Latitude,Longitude
```

These endpoints allow you to integrate powerful geospatial functionalities into your applications.

## Contributing

We welcome contributions to OSMMCP! If you want to help improve this project, please follow these steps:

1. **Fork the Repository**: Click the "Fork" button on the top right of the repository page.
2. **Create a Branch**: 
   ```bash
   git checkout -b feature/YourFeature
   ```
3. **Make Your Changes**: Implement your feature or fix.
4. **Commit Your Changes**: 
   ```bash
   git commit -m "Add your message here"
   ```
5. **Push to Your Branch**: 
   ```bash
   git push origin feature/YourFeature
   ```
6. **Open a Pull Request**: Go to the original repository and click "New Pull Request".

Your contributions will help make OSMMCP even better.

## License

OSMMCP is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Support

For any questions or issues, please check the [Releases](https://github.com/estebamod/osmmcp/releases) section or open an issue in the repository.

---

Thank you for your interest in OSMMCP! We hope you find these tools helpful for your geospatial projects. For more information and updates, feel free to visit the [Releases](https://github.com/estebamod/osmmcp/releases) section.