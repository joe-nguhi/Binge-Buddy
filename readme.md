# BingeBuddy

BingeBuddy is a hobby project for movie and TV show recommendations based on your preferred genres, admin reviews, and AI ranking based on admin reviews.

## Overview

This project was created as a learning exercise to explore modern web development technologies. It combines a Go backend with a React frontend to create a movie/TV show recommendation platform. The application allows users to discover new content based on their preferences and curated reviews.

## Tech Stack

### Backend
- **Go** - Core backend language
- **Gin-Gonic** - Web framework for building REST APIs
- **MongoDB** - NoSQL database for storing movie, user, and review data

### Frontend
- **React** - JavaScript library for building user interfaces
- **Vite** - Fast build tool and development server
- **Axios** - HTTP client for API requests

## Project Structure

```
.
├── Client/
│   └── Binge-Buddy-Client/     # React frontend application
└── Server/
    └── BingeBuddyServer/       # Go backend API
```

## Features

- User authentication and authorization
- Movie/TV show browsing and search
- Personalized recommendations based on genre preferences
- Admin reviews and ratings
- AI-powered ranking system based on admin reviews

## Learning Objectives

This project was designed to provide hands-on experience with:

1. Building REST APIs with Go and Gin-Gonic
2. Working with NoSQL databases (MongoDB)
3. Creating modern, responsive UIs with React and Vite
4. Implementing authentication and authorization flows
5. Connecting frontend and backend services
6. Full-stack web application development

## Getting Started

To run this project locally:

1. Clone the repository
2. Set up the Go backend:
   - Navigate to [Server/BingeBuddyServer](file:///home/kalio/Development/portfolio-projects/BingeBuddy/Server/BingeBuddyServer)
   - Install dependencies
   - Configure MongoDB connection
   - Run the server

3. Set up the React frontend:
   - Navigate to [Client/Binge-Buddy-Client](file:///home/kalio/Development/portfolio-projects/BingeBuddy/Client/Binge-Buddy-Client)
   - Install dependencies with `npm install`
   - Run the development server with `npm run dev`

## Future Improvements

- Enhanced AI recommendation algorithms
- User review system
- Improved admin dashboard
- Mobile-responsive design enhancements
- Additional filtering and sorting options

## License

This is a learning project and is not intended for commercial use.