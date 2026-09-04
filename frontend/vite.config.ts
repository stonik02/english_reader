import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  resolve: {
    preserveSymlinks: true,
  },
  optimizeDeps: {
    include: [
      '@reader/proto/reader/v1/auth_pb.js',
      '@reader/proto/reader/v1/library_pb.js',
      '@reader/proto/reader/v1/LibraryServiceClientPb.ts',
      '@reader/proto/reader/v1/reader_pb.js',
      '@reader/proto/reader/v1/ReaderServiceClientPb.ts',
      '@reader/proto/reader/v1/vocabulary_pb.js',
      '@reader/proto/reader/v1/VocabularyServiceClientPb.ts',
    ],
  },
  build: {
    commonjsOptions: {
      include: [/node_modules/, /src\/api\/gen/],
    },
  },
})
