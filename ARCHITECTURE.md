# Overview

## Three Tier Architecture

- Backing control plane REST API
- Web application for user control and observability
- Discord bot for in-chat control and convenience

## PRIMARY MOTIVATIONS

- extremely cheap, ideally totally free to run (ignoring backing kubernetes infra, already established).

# Subsystems

# REST API
- golang dockerized 
- AuthZ/AuthN managed with oAuth
- command and control endpoints for server management
- observability/metrics retrieval for the game server (uptime, connected clients, issued control commands, etc)

# Web Application
- consumes the REST API
- offers simpled start/stop capabilities
- shows observability and server status
- shows status of the command and control API (simple availability monitor)

# Discord bot
- super basic chat interface for discord server
- known start/stop/status commands
- IDEA: basic in-mem semantic analysis for natural language interpretation?

# Observability
- oTel?
- grafana?
- prometheus, probably?
