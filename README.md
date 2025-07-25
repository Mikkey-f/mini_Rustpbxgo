# MiniRustpbxgo

A real-time voice communication system integrating WebRTC, Go backend, and Rust services, enabling seamless audio transmission, speech processing, and AI interaction.

## Overview

MiniRustpbxgo combines WebRTC technology with a Go-based signaling server and Rust backend services to create a robust voice communication platform. It supports real-time audio streaming, speech recognition (ASR), text-to-speech (TTS), and AI-powered responses through LLM integration.

## Core Features

- **WebRTC Integration**: Real-time peer-to-peer audio communication with ICE candidate management
- **WebSocket Signaling**: Secure WebSocket connections for session setup and message routing
- **Speech Recognition (ASR)**: Real-time and final speech-to-text conversion
- **Text-to-Speech (TTS)**: High-quality audio generation from text with configurable parameters
- **AI Interaction**: Integration with LLM models for intelligent responses
- **Multi-provider Support**: Compatible with various ASR/TTS providers (Tencent, Aliyun, etc.)
- **Robot Management**: Configurable robot profiles with custom voices and behaviors

## Technical Stack

- **Frontend**: HTML5, JavaScript (WebRTC API)
- **Backend**: Go 1.24+ with Gin framework
- **Services**: Rust integration for performance-critical operations
- **Database**: MySQL for persistent storage
- **Caching**: Redis for session management
- **Network**: WebSocket for signaling, WebRTC for media streaming
- **AI**: OpenAI API integration for LLM capabilities

## Prerequisites

- Go 1.24 or higher
- Rust toolchain
- MySQL database
- Redis server
- Modern web browser with WebRTC support
- Valid API keys for ASR/TTS providers (if using cloud services)

## Installation

- Clone the repository：

```bash
git clone https://github.com/yourusername/miniRustpbxgo.git
cd miniRustpbxgo
```

- Install Go dependencies:

```bash
go mod download
```

- Configure environment variables:

    - Database connection strings

    - Redis server address

    - API keys for ASR/TTS providers

    - LLM API credentials

- Build and run the application:

```bash
go build -o miniRustpbxgo cmd/main.go
./miniRustpbxgo
```

- Start the Rust service (see Rust service documentation for details)

- Open index.html in a web browser to access the client interface

## Usage

1. Click "Establish Voice Connection" to initialize WebRTC session
2. Speak into your microphone for speech recognition
3. Receive real-time transcription and AI-generated responses
4. Adjust TTS parameters (speed, volume, speaker) through robot configuration
5. Click "Hang Up" to terminate the connection

## API Endpoints

- `GET /health`: Service health check
- `POST /robot/create`: Create new robot configuration
- `WS /out/webrtc/setup`: WebSocket endpoint for WebRTC signaling

## Configuration

Robot configurations and API keys are stored in the database. Use the `RobotCreate` API to configure new robot profiles with custom:

- Speech speed (0.5-2.0)
- Volume levels (0-10)
- Speaker voices
- Emotional tones
- System prompts for AI behavior

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

- Ensure WebSocket connection is not blocked by firewalls
- Verify SSL configuration for production deployments
- Check API key validity for ASR/TTS services
- Confirm proper CORS settings for WebRTC functionality