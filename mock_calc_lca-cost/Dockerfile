# Use Node.js runtime
FROM node:20-alpine

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source files
COPY . .

# Build TypeScript
RUN npm run build

# Create the directory and copy WASM file
RUN mkdir -p /app/dist/src/services && \
  cp node_modules/web-ifc/web-ifc-node.wasm /app/dist/src/services/ && \
  echo "Contents of /app/dist/src/services/:" && \
  ls -la /app/dist/src/services/ && \
  echo "Contents of /app/dist:" && \
  ls -la /app/dist && \
  echo "Contents of node_modules/web-ifc:" && \
  ls -la node_modules/web-ifc/

EXPOSE 4242

# Start the application with more verbose Node.js output
CMD ["node", "--trace-warnings", "dist/index.js"]