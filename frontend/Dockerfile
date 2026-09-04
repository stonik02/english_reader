FROM node:22-alpine AS build

WORKDIR /app
COPY package.json package-lock.json ./
COPY src/api/gen/package.json ./src/api/gen/package.json
RUN npm ci

COPY . ./
ARG VITE_GRPC_WEB_URL
ARG VITE_HTTP_API_URL
ENV VITE_GRPC_WEB_URL=$VITE_GRPC_WEB_URL
ENV VITE_HTTP_API_URL=$VITE_HTTP_API_URL
RUN npm run build

FROM nginx:1.27-alpine
COPY nginx.production.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
