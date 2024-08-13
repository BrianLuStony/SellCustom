# Use the official Node.js image
FROM node:18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Build the React application
RUN npm run build

# Install a simple web server to serve the static files
RUN npm install -g serve

# Expose port 5000
EXPOSE 5000

# Command to run the web server
CMD ["serve", "-s", "build", "-l", "5000"]

